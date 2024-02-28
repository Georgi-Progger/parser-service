package parser

import (
	"fmt"
	"log"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"main.go/internal/model/annoucement"
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

func Run(proxy string, lastIndex int) ([]annoucement.Annoucement, int, bool) {
	annoucements, lastIndex, isEnd := runWebDriver(link, proxy, userAgent, lastIndex)
	return annoucements, lastIndex, isEnd
}
func runWebDriver(link, proxy, userAgent string, lastIndex int) ([]annoucement.Annoucement, int, bool) {
	fmt.Printf("Driver Parsing %s\n, with proxy:%s", link, proxy)

	chromeDriverService, err := getChromeDriver(9515)
	if err != nil {
		log.Panic(err)
	}
	defer chromeDriverService.Stop()

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

	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", 9515)) // запуск драйвера на порту определенном
	if err != nil {
		log.Panic(err)
	}
	defer wd.Quit()

	carLinks := getLinksAnnoucements(wd, link)
	var annoucements []annoucement.Annoucement
	for i := lastIndex; i < len(carLinks); i++ {

		start := time.Now()
		fmt.Println("Parsing link:", carLinks[i])
		annoucementData := getAnnoucementData(wd, carLinks[i])
		if annoucementData.Model == "" {
			return annoucements, i, len(annoucements) == len(carLinks)
		}
		annoucements = append(annoucements, *annoucementData)

		end := time.Now()

		fmt.Printf("runWebDriver took %v to execute\n", end.Sub(start))
	}

	return annoucements, 0, (len(annoucements) == len(carLinks)) && (len(carLinks) > 0)
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
