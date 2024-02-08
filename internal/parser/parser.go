package parser

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	. "main.go/internal/model/annoucement"
)

func createDBConnection() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

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
		Link, Model, Price, Year, Generation, Mileage, History, PTS, Owners,
		Condition, Modification, Engine_Volume, Engine_Type, Transmission,
		Drive, Equipment, Body_Type, Color, Steering, VIN, Exchange, Location, Description
	) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23)
`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		annoucementInfo.Link, annoucementInfo.Model, annoucementInfo.Price, annoucementInfo.Year, annoucementInfo.Generation,
		annoucementInfo.Mileage, annoucementInfo.History, annoucementInfo.PTS, annoucementInfo.Owners,
		annoucementInfo.Condition, annoucementInfo.Modification, annoucementInfo.EngineVolume,
		annoucementInfo.EngineType, annoucementInfo.Transmission, annoucementInfo.Drive,
		annoucementInfo.Equipment, annoucementInfo.BodyType, annoucementInfo.Color,
		annoucementInfo.Steering, annoucementInfo.VIN, annoucementInfo.Exchange, annoucementInfo.Location, annoucementInfo.Description,
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

func ExtractCarInfo(link, htmlContent, targetClass string) (*Annoucement, error) {
	reader := strings.NewReader(htmlContent)
	tokenizer := html.NewTokenizer(reader)

	carInfo := &Annoucement{}
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
