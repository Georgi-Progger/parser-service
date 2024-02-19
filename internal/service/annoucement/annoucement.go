package annoucement

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"main.go/internal/model/annoucement"
	"main.go/internal/model/proxy"
	"main.go/internal/parser"
)

func (s *Service) GetAnnoucements(c echo.Context) error {
	ctx := c.Request().Context()

	annoucementRepo := annoucement.NewRepository(s.db)
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

	annoucementRepo := annoucement.NewRepository(s.db)
	proxyRepo := proxy.NewRepository(s.db)

	activeProxy, err := proxyRepo.GetActiveProxy(ctx)
	if err != nil {
		return c.String(http.StatusNonAuthoritativeInfo, err.Error())
	}

	annoucementInfo := parser.Run(activeProxy.Body)
	for idx := range annoucementInfo {
		isExists := annoucementRepo.LinkExists(ctx, annoucementInfo[idx].Link)
		if !isExists {
			annoucementRepo.SetAnnoucement(ctx, annoucementInfo[idx])
		}
	}
	return c.JSON(http.StatusOK, "All annoucements is save")
}
