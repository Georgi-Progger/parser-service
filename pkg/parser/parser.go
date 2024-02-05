package parser

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"golang.org/x/net/html"

	. "main.go/internal/model/annoucement"
)

const (
	dbHost     = "localhost"
	dbPort     = "5432"
	dbUser     = "geor"
	dbPassword = "georkryt"
	dbName     = "parser_db"
)

func createDBConnection() (*sql.DB, error) {
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func InsertIntoDB(annoucementInfo Annoucement) error {
	db, err := createDBConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(`
			INSERT INTO annoucement (
				Model, Price, Year, Generation, Mileage, History, PTS, Owners,
				Condition, Modification, Engine_Volume, Engine_Type, Transmission,
				Drive, Equipment, Body_Type, Color, Steering, VIN, Exchange
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
		`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		annoucementInfo.Model, annoucementInfo.Price, annoucementInfo.Year, annoucementInfo.Generation,
		annoucementInfo.Mileage, annoucementInfo.History, annoucementInfo.PTS, annoucementInfo.Owners,
		annoucementInfo.Condition, annoucementInfo.Modification, annoucementInfo.EngineVolume,
		annoucementInfo.EngineType, annoucementInfo.Transmission, annoucementInfo.Drive,
		annoucementInfo.Equipment, annoucementInfo.BodyType, annoucementInfo.Color,
		annoucementInfo.Steering, annoucementInfo.VIN, annoucementInfo.Exchange,
	)

	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func ExtractLinks(htmlContent string, targetClass string) ([]string, error) {
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
								links = append(links, attr.Val)
							}
						}
					}
				}
			}
		}
	}

}
func ExtractCarInfo(htmlContent string, targetClass string) (*Annoucement, error) {
	reader := strings.NewReader(htmlContent)
	tokenizer := html.NewTokenizer(reader)

	carInfo := &Annoucement{}
	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			return carInfo, nil
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			if token.Data == "h1" {
				var foundTargetClass bool

				for _, attr := range token.Attr {
					if attr.Key == "data-marker" && strings.Contains(attr.Val, "item-view/title-info") {
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
						} else if tokenType == html.EndTagToken {
							break
						}
					}
					finalString := strings.Split(inputString, ", ")
					carInfo.Model = finalString[0]
					inputString = ""
				}
			} else if token.Data == "span" {
				var foundTargetClass bool

				for _, attr := range token.Attr {
					if attr.Key == "data-marker" && strings.Contains(attr.Val, "item-view/item-price") {
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
						} else if tokenType == html.EndTagToken {
							break
						}
					}
					carInfo.Price = inputString
					inputString = ""
				}
			} else if token.Data == "li" {
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
							// fieldsInfo := make(map[string]string)
							// fieldsInfo[finalString[0]] = finalString[1]

							setCarInfoField(carInfo, finalString)
							break
						}
					}
					inputString = ""
				}
			}
		}
	}
}

func setCarInfoField(info *Annoucement, infoString []string) {
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
