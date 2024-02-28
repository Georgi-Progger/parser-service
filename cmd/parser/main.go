package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"main.go/internal/app"
)

func main() {
	cfg := app.Config{}
	if err := godotenv.Load(); err != nil {
		panic("No .env file found")
	}
	port := os.Getenv("PORT")
	cfg.EchoPort = string(fmt.Sprintf(":%s", port))
	app, err := app.NewApplication(&cfg)
	if err != nil {
		panic(err)
	}

	if err := app.Start(); err != nil {
		panic(err)
	}
}
