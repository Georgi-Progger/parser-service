package app

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"main.go/internal/service/annoucement"
	"main.go/internal/service/proxy"
)

func newRoute(a *annoucement.Service, p *proxy.Service) *echo.Echo {
	router := echo.New()

	router.Use(middleware.Logger())

	annoucementsGroup := router.Group("/annoucements")
	annoucementsGroup.GET("/search/:page", a.GetAnnoucements)
	annoucementsGroup.POST("/save", a.SetAnnoucements)

	// proxiesGroup := router.Group("/proxy")
	// proxiesGroup.GET("", p.GetActiveProxy)
	return router
}
