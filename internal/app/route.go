package app

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"go.uber.org/zap"
	"main.go/internal/service/annoucement"
	"main.go/internal/service/proxy"
)

func newRoute(a *annoucement.Service, p *proxy.Service) *echo.Echo {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	router := echo.New()

	router.Use(middleware.Logger())
	router.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger.Info("Incoming request",
				zap.String("method", c.Request().Method),
				zap.String("uri", c.Request().RequestURI),
			)
			return next(c)
		}
	})

	annoucementsGroup := router.Group("/annoucements")
	annoucementsGroup.GET("/search/:page", a.GetAnnoucements)
	annoucementsGroup.POST("/save", a.SetAnnoucements)

	// proxiesGroup := router.Group("/proxy")
	// proxiesGroup.GET("", p.GetActiveProxy)

	return router
}
