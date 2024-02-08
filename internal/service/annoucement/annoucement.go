package annoucement

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"main.go/internal/model/annoucement"
)

func (s *Service) GetAnnoucements(c echo.Context) error {
	annoucementRepo := annoucement.NewRepository(s.db)
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Некорректное значение страницы.")
	}
	res, err := annoucementRepo.GetAnnoucement(c.Request().Context(), page)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, res)
}
