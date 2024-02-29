package app

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"main.go/internal/service/annoucement"
	"main.go/internal/service/proxy"
)

type (
	App struct {
		AnnoucementService *annoucement.Service
		ProxyService       *proxy.Service
		config             *Config
		Echo               *echo.Echo
	}

	Config struct {
		EchoPort string
		Db       Db
	}

	Db struct {
		DbConnect string
	}
)

func (a *App) Start() error {
	return a.Echo.Start(a.config.EchoPort)
}

func NewApplication(cfg *Config) (*App, error) {
	db, err := ConnectDatabase()
	if err != nil {
		log.Panic("connectionString error..")
	}
	annoucementService := annoucement.NewService(db)
	proxyService := proxy.NewService(db)
	return &App{
		AnnoucementService: annoucementService,
		ProxyService:       proxyService,
		config:             cfg,
		Echo:               newRoute(annoucementService, proxyService),
	}, nil
}

func ConnectDatabase() (*sql.DB, error) {
	cfg := Config{}
	if err := godotenv.Load(); err != nil {
		panic("No .env file found")
	}

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	// dbEncodePassword := os.Getenv("DB_PASSWORD")
	// dbPassword, err := url.QueryUnescape(dbEncodePassword)
	// if err != nil {
	// 	panic(err)
	// }
	cfg.Db.DbConnect = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", cfg.Db.DbConnect)
	if err != nil {
		log.Panic("no connect to database...")
	}

	//defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Panic(err)
	}
	return db, nil
}
