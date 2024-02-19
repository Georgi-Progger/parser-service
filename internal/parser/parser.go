package parser

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"

	_ "github.com/lib/pq"
	"main.go/internal/model/annoucement"
)

func extractLinks(htmlContent string, targetClass string) ([]string, error) {
	var links []string

	reader := strings.NewReader(htmlContent)
	tokenizer := html.NewTokenizer(reader)

	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			return links, nil
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			if token.Data == "a" {
				for _, attr := range token.Attr {
					if attr.Key == "class" && strings.Contains(attr.Val, targetClass) {
						for _, attr := range token.Attr {
							if attr.Key == "href" {
								link := fmt.Sprintf("https://www.avito.ru%s", attr.Val)
								links = append(links, link)
							}
						}
					}
				}
			}
		}
	}
}

func extractCarInfo(link, htmlContent, targetClass string) (*annoucement.Annoucement, error) {
	reader := strings.NewReader(htmlContent)
	tokenizer := html.NewTokenizer(reader)

	carInfo := &annoucement.Annoucement{}
	carInfo.Link = link
	var foundTitleInfo bool

	carInfo.Description = setDescrtiption(htmlContent)
	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			return carInfo, nil
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			switch token.Data {
			case "h1":
				for _, attr := range token.Attr {
					if attr.Key == "data-marker" && strings.Contains(attr.Val, "item-view/title-info") {
						foundTitleInfo = true
						break
					}
				}
				if foundTitleInfo {
					inputString := ""
					for {
						tokenType := tokenizer.Next()
						token := tokenizer.Token()
						if tokenType == html.TextToken {
							inputString += token.Data
						} else if tokenType == html.EndTagToken {
							break
						}
					}
					finalString := strings.Split(inputString, ", ")
					carInfo.Model = finalString[0]
					inputString = ""
				}
			case "span":
				for _, attr := range token.Attr {
					if attr.Key == "data-marker" && strings.Contains(attr.Val, "item-view/item-price") {
						inputString := ""
						for {
							tokenType := tokenizer.Next()
							token := tokenizer.Token()
							if tokenType == html.TextToken {
								inputString += token.Data
							} else if tokenType == html.EndTagToken {
								break
							}
						}
						carInfo.Price = inputString
						break
					}
					if attr.Key == "class" && strings.Contains(attr.Val, "style-item-address__string-wt61A") {
						inputString := ""
						for {
							tokenType := tokenizer.Next()
							token := tokenizer.Token()
							if tokenType == html.TextToken {
								inputString += token.Data
							} else if tokenType == html.EndTagToken {
								break
							}
						}
						carInfo.Location = inputString
						break
					}
				}

			case "li":
				var foundTargetClass bool
				for _, attr := range token.Attr {
					if attr.Key == "class" && strings.Contains(attr.Val, targetClass) {
						foundTargetClass = true
						break
					}
				}
				if foundTargetClass {
					inputString := ""
					for {
						tokenType := tokenizer.Next()
						token := tokenizer.Token()
						if tokenType == html.TextToken {
							inputString += token.Data
						} else if tokenType == html.EndTagToken && token.Data == "li" {
							finalString := strings.Split(inputString, ": ")
							setCarInfoField(carInfo, finalString)
							break
						}
					}
				}
			}
		}
	}
}

func setCarInfoField(info *annoucement.Annoucement, infoString []string) {
	switch infoString[0] {
	case "Год выпуска":
		info.Year = infoString[1]
	case "Поколение":
		info.Generation = infoString[1]
	case "Пробег":
		info.Mileage = infoString[1]
	case "История пробега":
		info.History = infoString[1]
	case "ПТС":
		info.PTS = infoString[1]
	case "Владельцев по ПТС":
		info.Owners = infoString[1]
	case "Состояние":
		info.Condition = infoString[1]
	case "Модификация":
		info.Modification = infoString[1]
	case "Объём двигателя":
		info.EngineVolume = infoString[1]
	case "Тип двигателя":
		info.EngineType = infoString[1]
	case "Коробка передач":
		info.Transmission = infoString[1]
	case "Привод":
		info.Drive = infoString[1]
	case "Комплектация":
		info.Equipment = infoString[1]
	case "Тип кузова":
		info.BodyType = infoString[1]
	case "Цвет":
		info.Color = infoString[1]
	case "Руль":
		info.Steering = infoString[1]
	case "VIN или номер кузова":
		info.VIN = infoString[1]
	case "Обмен":
		info.Exchange = infoString[1]
	}
}

func setDescrtiption(htmlContent string) string {
	reader := strings.NewReader(htmlContent)
	tokenizer := html.NewTokenizer(reader)

	inputString := ""
	inDescription := false
	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			return ""
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			if token.Data == "div" {
				for _, attr := range token.Attr {
					if attr.Key == "data-marker" && attr.Val == "item-view/item-description" {
						inDescription = true
						break
					}
				}
			} else if (token.Data == "li" || token.Data == "p") && inDescription {
				tokenizer.Next()
				textToken := tokenizer.Token()
				inputString += textToken.Data + " "
				inputString = strings.ReplaceAll(inputString, "strong", "")
			}
		case html.EndTagToken:
			token := tokenizer.Token()
			if token.Data == "div" && inDescription {
				return inputString
			}
		}
	}

}
