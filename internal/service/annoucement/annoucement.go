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

	annoucementInfo, lastIndex, isEnd := parser.Run(activeProxy.Body)
	for !isEnd {
		proxyRepo.UpdateProxy(ctx, activeProxy.Body)
		activeProxy, err = proxyRepo.GetActiveProxy(ctx)
		annoucementInfo, lastIndex, isEnd = parser.Run(activeProxy.Body)
	}
	for i := lastIndex; i < len(annoucementInfo); i++ {
		isExists := annoucementRepo.LinkExists(ctx, annoucementInfo[i].Link)
		if !isExists {
			annoucementRepo.SetAnnoucement(ctx, annoucementInfo[i])
		}
	}
	return c.JSON(http.StatusOK, "All annoucements is save")
}
