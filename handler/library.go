package handler

import (
	"fmt"
	"go-steam/model"
	"go-steam/services"
	"go-steam/util"
	"go-steam/views"

	"github.com/labstack/echo"
)

func GetLibraryFiltered(c echo.Context) error {
	title := c.QueryParam("filter")
	fmt.Printf("Endpoint: ShowLibraryFiltered: %s \n", title)

	data, err := services.GetLibrary(title)
	if err != nil {
		c.Logger().Error(err)
		return util.Render(c, views.LibraryCards(model.Library{}))
	}

	fmt.Printf("Returning #%d Games \n", len(data.Cards))
	return util.Render(c, views.LibraryCards(data))
}

func GetLibrary(c echo.Context) error {
	data, err := services.GetLibrary("")
	if err != nil {
		c.Logger().Error(err)
		return util.Render(c, views.LibraryPage(model.Library{}))
	}
	fmt.Printf("Returning #%d Games \n", len(data.Cards))
	return util.Render(c, views.LibraryPage(data))
}

func UpdateLibrary(c echo.Context) error {
	getOwnedGames, err := services.GetSteamUserGames()
	if err != nil {
		fmt.Printf("Fail: %s", err.Error())
	}
	services.UpdateLibrary(getOwnedGames)

	data, err := services.GetLibrary("")
	if err != nil {
		c.Logger().Error(err)
		return util.Render(c, views.LibraryCards(model.Library{}))
	}

	fmt.Printf("Returning #%d Games \n", len(data.Cards))
	return util.Render(c, views.LibraryCards(data))
}
