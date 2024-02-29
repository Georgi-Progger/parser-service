package annoucement

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"main.go/internal/parser"
	"main.go/internal/repositories"
)

func (s *Service) GetAnnoucements(c echo.Context) error {
	ctx := c.Request().Context()

	annoucementRepo := repositories.NewRepository(s.db)
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Некорректное значение страницы.")
	}
	res, err := annoucementRepo.GetAnnoucements(ctx, page)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, res)
}

func (s *Service) SetAnnoucements(c echo.Context) error {
	ctx := c.Request().Context()

	proxyRepo := repositories.NewRepository(s.db)

	activeProxy, err := proxyRepo.GetActiveProxy(ctx)
	if err != nil {
		return c.String(http.StatusNonAuthoritativeInfo, err.Error())
	}
	isBlockProxy := parser.Run(activeProxy.Body, s.db, c)
	if isBlockProxy {
		proxyRepo.BlockProxy(ctx, activeProxy.Body)
	}
	return c.JSON(http.StatusOK, "All annoucements is save")
}
