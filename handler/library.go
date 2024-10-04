package handler

import (
	"go-steam/services"
	"go-steam/views"

	"github.com/labstack/echo"
)

type Library struct{}

func (l Library) ShowLibrary(c echo.Context) error {
	data := services.GetLibrary()
	return render(c, views.ShowLibrary(data))
}
