package handler

import (
	"go-steam/services"
	"go-steam/views"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

type Library struct{}

func (l Library) ShowLibrary(c echo.Context) error {
	data := services.GetLibrary()
	return render(c, views.ShowLibrary(data))
}

func (l Library) ShowGame(c echo.Context) error {
	param := c.Param("AppID")
	id, err := strconv.Atoi(param)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusBadRequest, "Invalid AppID")
	}

	data := services.GetGame(id)
	return render(c, views.ShowGame(data))
}
