package handler

import (
	"fmt"
	"go-steam/services"
	"go-steam/util"
	"go-steam/views"

	"github.com/labstack/echo"
)

func GetLibraryFiltered(c echo.Context) error {
	title := c.QueryParam("filter")
	fmt.Printf("Endpoint: ShowLibraryFiltered: %s \n", title)

	data := services.GetLibrary(title)
	fmt.Printf("Returning #%d Games \n", len(data.Cards))
	return util.Render(c, views.LibraryCards(data))
}

func GetLibrary(c echo.Context) error {
	fmt.Println("Endpoint: GetLibrary")

	data := services.GetLibrary("")
	fmt.Printf("Returning #%d Games \n", len(data.Cards))
	return util.Render(c, views.LibraryPage(data))
}

func UpdateLibrary(c echo.Context) error {
	fmt.Println("Endpoint: UpdateLibrary")
	getOwnedGames, err := services.GetSteamUserGames()
	if err != nil {
		fmt.Printf("Fail: %s", err.Error())
	}
	services.UpdateLibrary(getOwnedGames)
	data := services.GetLibrary("")
	fmt.Printf("Returning #%d Games \n", len(data.Cards))
	return util.Render(c, views.LibraryCards(data))
}
