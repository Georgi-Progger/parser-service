package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"golang.org/x/sync/errgroup"
	. "main.go/pkg/parser"
)

func generateRandomUserAgent() string {
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:97.0) Gecko/20100101 Firefox/97.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Edge/18.19041",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
	}

	randomIndex := rand.Intn(len(userAgents))
	return userAgents[randomIndex]
}

var links = []string{
	"https://www.avito.ru/all/avtomobili/asia-ASgBAgICAUTgtg3alyg?cd=1&s=104",
	"https://www.avito.ru/all/avtomobili/arcfox-ASgBAgICAUTgtg2CmuQQ?cd=1&s=104",
}

var mu sync.Mutex

func startChromeDriver(port int) (*selenium.Service, error) {
	opts := []selenium.ServiceOption{}
	service, err := selenium.NewChromeDriverService("./assets/chromedriver", port, opts...)
	if err != nil {
		return nil, err
	}
	return service, nil
}

func main() {
	var g errgroup.Group

	for i, link := range links {
		i, link := i, link

		g.Go(func() error {
			return runWebDriver(i, link)
		})
	}
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}

func runWebDriver(index int, link string) error {
	chromeDriverService, err := startChromeDriver(9515 + index)
	if err != nil {
		return fmt.Errorf("Ошибка при запуске ChromeDriver: %v", err)
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

	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", 9515+index))
	if err != nil {
		return fmt.Errorf("Error creating WebDriver: %v", err)
	}
	defer wd.Quit()

	log.Printf("Driver %d started for %s", index, link)

	err = wd.Get(fmt.Sprintf("%s", link))
	if err != nil {
		return fmt.Errorf("Error opening page: %v", err)
	}
	html, err := wd.PageSource()
	if err != nil {
		return fmt.Errorf("Error getting page source: %v", err)
	}
	fmt.Printf("Driver %d - Ссылки на объявления:\n", index)

	carLinks, err := ExtractLinks(html, "iva-item-sliderLink-uLz1v")
	if err != nil {
		log.Fatal("Error get a links:", err)
	}

	for _, link := range carLinks {
		fmt.Printf("Driver %d - %s\n", index, fmt.Sprintf("https://www.avito.ru%s", link))
		time.Sleep(3 * time.Second)
		err = wd.Get(fmt.Sprintf("https://www.avito.ru%s", link))
		if err != nil {
			return fmt.Errorf("Error opening page: %v", err)
		}

		html, err := wd.PageSource()
		if err != nil {
			return fmt.Errorf("Error getting page source: %v", err)
		}

		carInfo, err := ExtractCarInfo(html, "params-paramsList__item-_2Y2O")
		if err != nil {
			return fmt.Errorf("Error parsing car info: %v", err)
		}

		mu.Lock()
		if err := InsertIntoDB(*carInfo); err != nil {
			mu.Unlock()
			return fmt.Errorf("Error inserting data into the database: %v", err)
		}
		mu.Unlock()

	}
	time.Sleep(time.Second * 5)
	return nil
}
