package main

import (
	"go-steam/db"
	"go-steam/handler"

	"github.com/labstack/echo"
)

func main() {
	db.InitDatabase()
	e := echo.New()
	e.GET("/", handler.GetLibrary)
	e.GET("/getLibraryFilter", handler.GetLibraryFiltered)
	e.GET("/updateLibrary", handler.UpdateLibrary)
	e.GET("/getSteamUserLibrary/:AppID", handler.GetDetailsPage)
	e.GET("/updateAchivements/:AppID", handler.UpdateAchievements)
	e.GET("/friendsList", handler.GetFriends)
	e.Logger.Fatal(e.Start(":8000"))
}
