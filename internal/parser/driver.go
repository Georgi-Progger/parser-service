package parser

import (
	"fmt"
	"log"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"main.go/internal/model/annoucement"
)

var (
	userAgent string = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:97.0) Gecko/20100101 Firefox/97.0"

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

func Run(proxy string) []annoucement.Annoucement {
	annoucements := runWebDriver(link, proxy, userAgent)
	return annoucements
}
func runWebDriver(link string, proxy string, userAgent string) []annoucement.Annoucement {
	fmt.Printf("Driver Parsing %s\n", link)

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
	for idx := range carLinks {
		fmt.Println("Parsing link:", carLinks[idx])
		annoucementData := getAnnoucementData(wd, carLinks[idx])
		annoucements = append(annoucements, *annoucementData)
	}

	return annoucements
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
