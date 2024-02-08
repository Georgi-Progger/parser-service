package parser

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:97.0) Gecko/20100101 Firefox/97.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Edge/18.19041",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
}

var links = []string{
	"https://www.avito.ru/all/avtomobili/vaz_lada-ASgBAgICAUTgtg3GmSg?cd=1",
	"https://www.avito.ru/all/avtomobili/toyota-ASgBAgICAUTgtg20mSg?cd=1",
	"https://www.avito.ru/all/avtomobili/kia-ASgBAgICAUTgtg3KmCg?cd=1",
	"https://www.avito.ru/all/avtomobili/mercedes-benz-ASgBAgICAUTgtg3omCg?cd=1",
	"https://www.avito.ru/all/avtomobili/geely-ASgBAgICAUTgtg2gmCg?cd=1",
}
var proxy = []string{
	"50.168.210.232:80",
	"50.204.219.224:80",
	"50.217.226.43:80",
	"172.67.3.7:80",
	"172.67.43.166:80",
}

var mu sync.Mutex

/* func runWebDriverForLinks() []string {
	chromeDriverService, err := startChromeDriver(9514)
	if err != nil {
		log.Fatal("Ошибка при запуске ChromeDriver:", err)
	}
	defer chromeDriverService.Stop()

	caps := selenium.Capabilities{
		"proxy": map[string]interface{}{
			"httpProxy": "172.67.43.72:80",
			"proxyType": "manual",
		},
	}
	capsArgs := []string{"--headless", "--disable-gpu", "--no-sandbox"}
	capsArgs = append(capsArgs, fmt.Sprintf("--user-agent=%s", generateRandomUserAgent()))
	caps.AddChrome(chrome.Capabilities{
		Path:            "",
		Args:            capsArgs,
		ExcludeSwitches: []string{},
		Extensions:      []string{},
		LocalState:      map[string]interface{}{},
		Prefs:           map[string]interface{}{},
		Detach:          new(bool),
		DebuggerAddr:    "",
		MinidumpPath:    "",
		MobileEmulation: &chrome.MobileEmulation{},
		W3C:             false,
	})

	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", 9514))
	if err != nil {
		log.Fatal("Error creating WebDriver:", err)
	}
	defer wd.Quit()

	link := "https://www.avito.ru/all/avtomobili"

	err = wd.Get(fmt.Sprintf("%s", link))
	if err != nil {
		log.Fatal("Error opening page:", err)
	}
	html, err := wd.PageSource()
	if err != nil {
		log.Fatal("Error getting page source:", err)
	}

	carLinks, err := ExtractLinks(html, "popular-rubricator-link-Hrkjd")
	if err != nil {
		log.Fatal("Error get a links:", err)
	}
	return carLinks
} */

func Run() {
	var wg sync.WaitGroup
	for i, link := range links {
		wg.Add(1)
		go func(index int, link string) {
			defer wg.Done()
			if err := runWebDriver(index, fmt.Sprintf("%s&s=104", link), proxy[index], userAgents[index]); err != nil {
				log.Printf("Error in runWebDriver for link %s: %v", link, err)
			}
		}(i, link)
	}
	wg.Wait()
}
func getChromeDriver(port int) (*selenium.Service, error) {
	opts := []selenium.ServiceOption{}
	service, err := selenium.NewChromeDriverService("./assets/chromedriver", port, opts...)
	if err != nil {
		return nil, err
	}
	return service, nil
}

func runWebDriver(index int, link string, proxy string, userAgent string) error {
	chromeDriverService, err := getChromeDriver(9515 + index)
	if err != nil {
		return fmt.Errorf("Ошибка при запуске ChromeDriver: %v", err)
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
		Path:            "",
		Args:            capsArgs,
		ExcludeSwitches: []string{},
		Extensions:      []string{},
		LocalState:      map[string]interface{}{},
		Prefs:           map[string]interface{}{},
		Detach:          new(bool),
		DebuggerAddr:    "",
		MinidumpPath:    "",
		MobileEmulation: &chrome.MobileEmulation{},
		W3C:             false,
	})

	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", 9515+index))
	if err != nil {
		return fmt.Errorf("Error creating WebDriver: %v", err)
	}
	defer wd.Quit()

	fmt.Printf("Driver %d started for %s and proxy:%s\n", index, link, proxy)

	if err := scrapeAndInsertData(index, wd, link); err != nil {
		return fmt.Errorf("Error during scraping and inserting data: %v", err)
	}

	return nil
}

func scrapeAndInsertData(index int, wd selenium.WebDriver, link string) error {
	time.Sleep(time.Duration(rand.Intn(5)) * time.Second)

	err := wd.Get(fmt.Sprintf("%s", link))
	if err != nil {
		return fmt.Errorf("Error opening page: %v", err)
	}
	time.Sleep(time.Duration(rand.Intn(5)) * time.Second)

	html, err := wd.PageSource()
	if err != nil {
		return fmt.Errorf("Error getting page source: %v", err)
	}
	fmt.Printf("Driver %d - Ссылки на объявления:\n", index)
	carLinks, err := ExtractLinks(html, "iva-item-sliderLink-uLz1v")
	if err != nil {
		return fmt.Errorf("Error getting links: %v", err)
	}

	for _, link := range carLinks {
		fmt.Printf("Driver %d - %s\n", index, fmt.Sprintf("https://www.avito.ru%s", link))
		time.Sleep(3 * time.Second)
		fullAddress := fmt.Sprintf("https://www.avito.ru%s", link)
		err = wd.Get(fullAddress)
		if err != nil {
			return fmt.Errorf("Error opening page: %v", err)
		}
		time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
		html, err := wd.PageSource()
		if err != nil {
			return fmt.Errorf("Error getting page source: %v", err)
		}

		carInfo, err := ExtractCarInfo(fullAddress, html, "params-paramsList__item-_2Y2O")
		if err != nil {
			return fmt.Errorf("Error parsing car info: %v", err)
		}

		if carInfo == nil {
			fmt.Printf("Driver %d - No car info found. Restarting...\n", index)
			if err := scrapeAndInsertData(index, wd, fmt.Sprintf("https://www.avito.ru%s", link)); err != nil {
				return fmt.Errorf("Error during scraping and inserting data: %v", err)
			}
			continue
		}

		fmt.Printf("Driver %d", index)
		if err := InsertIntoDB(*carInfo); err != nil {
			mu.Unlock()
			return fmt.Errorf("Error inserting data into the database: %v", err)
		}
	}

	return nil
}
