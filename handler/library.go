package handler

import (
	"fmt"
	"go-steam/model"
	"go-steam/services"
	"go-steam/util"
	"go-steam/views"

	"github.com/labstack/echo"
)

func GetLibraryFilteredTitle(c echo.Context) error {
	title := c.QueryParam("filter")
	fmt.Printf("Endpoint: ShowLibraryFiltered: %s \n", title)

	data, err := services.GetLibrary(title, "")
	if err != nil {
		c.Logger().Error(err)
		return util.Render(c, views.LibraryCards(model.Library{}))
	}

	filterOptions, err := services.GetFilterOptions()
	if err != nil {
		c.Logger().Error(err)
		return util.Render(c, views.LibraryCards(model.Library{}))
	}
	data.FilterOptions = filterOptions

	return util.Render(c, views.LibraryCards(data))
}

func GetLibraryFilterGenres(c echo.Context) error {
	genre := util.GetParam(c, "Genre")

	data, err := services.GetLibrary("", genre)
	if err != nil {
		c.Logger().Error(err)
		return util.Render(c, views.LibraryCards(model.Library{}))
	}

	return util.Render(c, views.LibraryCards(data))
}

func GetLibrary(c echo.Context) error {
	data, err := services.GetLibrary("", "")
	if err != nil {
		c.Logger().Error(err)
		return util.Render(c, views.LibraryPage(model.Library{}))
	}

	filterOptions, err := services.GetFilterOptions()
	if err != nil {
		c.Logger().Error(err)
		return util.Render(c, views.LibraryPage(model.Library{}))
	}
	data.FilterOptions = filterOptions

	return util.Render(c, views.LibraryPage(data))
}

func UpdateLibrary(c echo.Context) error {
	getOwnedGames, err := services.GetSteamUserGames()
	if err != nil {
		fmt.Printf("Fail: %s", err.Error())
	}
	services.UpdateLibrary(getOwnedGames)

	data, err := services.GetLibrary("", "")
	if err != nil {
		c.Logger().Error(err)
		return util.Render(c, views.LibraryCards(model.Library{}))
	}

	return util.Render(c, views.LibraryCards(data))
}
