package parser

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/labstack/echo"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"main.go/internal/model/annoucement"
	"main.go/internal/repositories"
)

var (
	userAgent string = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36"

	link string = "https://www.avito.ru/all/avtomobili?cd=1&s=104"
)

func getChromeDriver(port int) (*selenium.Service, error) {
	opts := []selenium.ServiceOption{}
	service, err := selenium.NewChromeDriverService("./assets/chromedriver", port, opts...)
	if err != nil {
		return nil, err
	}
	return service, nil
}

func Run(proxy string, db *sql.DB, c echo.Context) bool {
	isBlockProxy := runWebDriver(proxy, db, c)
	return isBlockProxy
}

func getCaps(proxy string) selenium.Capabilities {
	caps := selenium.Capabilities{
		"proxy": map[string]interface{}{
			"httpProxy": proxy,
			"proxyType": "manual",
		},
	}
	capsArgs := []string{"--headless", "--disable-gpu", "--no-sandbox"}
	capsArgs = append(capsArgs, fmt.Sprintf("--user-agent=%s", userAgent))
	caps.AddChrome(chrome.Capabilities{
		// MobileEmulation: &chrome.MobileEmulation{},
		// Path:            "",
		// Args:            capsArgs,
		// ExcludeSwitches: []string{},
		// Extensions:      []string{},
		// LocalState:      map[string]interface{}{},
		// Prefs:           map[string]interface{}{},
		// Detach:          new(bool),
		// DebuggerAddr:    "",
		// MinidumpPath:    "",
		// W3C:             false,
	})
	return caps
}
func runWebDriver(proxy string, db *sql.DB, c echo.Context) bool {
	ctx := c.Request().Context()
	caps := getCaps(proxy)
	chromeDriverService, err := getChromeDriver(9515)
	if err != nil {
		log.Panic(err)
	}
	defer chromeDriverService.Stop()
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", 9515)) // запуск драйвера на порту определенном
	if err != nil {
		log.Panic(err)
	}
	defer wd.Quit()

	annoucementRepo := repositories.NewRepository(db)
	carLinks := getLinksAnnoucements(wd, link)

	for idx := range carLinks {
		start := time.Now()
		annoucementData := getAnnoucementData(wd, carLinks[idx])
		if annoucementRepo.LinkExists(ctx, carLinks[idx]) {
			break
		}
		if annoucementData.Model == "" {
			return true
		}
		fmt.Println("Parsing link:", carLinks[idx])
		annoucementRepo.SetAnnoucement(c.Request().Context(), *annoucementData)
		end := time.Now()

		fmt.Printf("runWebDriver took %v to execute\n", end.Sub(start))
	}
	return false
}

func getLinksAnnoucements(wd selenium.WebDriver, link string) []string {
	if err := wd.Get(link); err != nil {
		log.Panic(err)
	}

	html, err := wd.PageSource()
	if err != nil {
		log.Panic(err)
	}
	carLinks, err := extractLinks(html, "iva-item-sliderLink-uLz1v")
	if err != nil {
		log.Panic(err)
	}

	return carLinks
}

func getAnnoucementData(wd selenium.WebDriver, link string) *annoucement.Annoucement {
	if err := wd.Get(link); err != nil {
		log.Panic(err)
	}

	html, err := wd.PageSource()
	if err != nil {
		log.Panic(err)
	}

	carInfo, err := extractCarInfo(link, html, "params-paramsList__item-_2Y2O")
	if err != nil {
		log.Panic(err)
	}

	return carInfo
}
