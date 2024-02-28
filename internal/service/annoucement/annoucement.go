package annoucement

import (
	"log"
	"net/http"
	"strconv"
	"time"

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

	annoucementRepo := repositories.NewRepository(s.db)
	proxyRepo := repositories.NewRepository(s.db)

	activeProxy, err := proxyRepo.GetActiveProxy(ctx)
	if err != nil {
		return c.String(http.StatusNonAuthoritativeInfo, err.Error())
	}
	lastIndex := 0
	annoucementInfo, lastIndex, isEnd := parser.Run(activeProxy.Body, lastIndex)
	for !isEnd {

		time.Sleep(2 * time.Second)
		proxyRepo.UpdateProxy(ctx, activeProxy.Body)
		activeProxy, err = proxyRepo.GetActiveProxy(ctx)
		if err != nil {
			log.Panic(err)
		}
		annoucementInfo, lastIndex, isEnd = parser.Run(activeProxy.Body, lastIndex)
	}
	for i := lastIndex; i < len(annoucementInfo); i++ {
		isExists := annoucementRepo.LinkExists(ctx, annoucementInfo[i].Link)
		if !isExists {
			annoucementRepo.SetAnnoucement(ctx, annoucementInfo[i])
		}
	}
	return c.JSON(http.StatusOK, "All annoucements is save")
}
