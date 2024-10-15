package main

import (
	"go-steam/db"
	"go-steam/handler"
	"go-steam/util"

	"github.com/labstack/echo"
)

func main() {
	db.InitDatabase()
	util.InitConfig()

	e := echo.New()
	e.GET("/", handler.GetLibrary)
	e.GET("/getLibraryFilter/title", handler.GetLibraryFilteredTitle)
	e.POST("/getLibraryFilter/:Genre", handler.GetLibraryFilterGenres)
	e.GET("/updateLibrary", handler.UpdateLibrary)
	e.GET("/getSteamUserLibrary/:AppID", handler.GetDetailsPage)
	e.GET("/filterAchivements/:AppID/:Filter", handler.GetAchievements)
	e.GET("/updateAchivements/:AppID", handler.UpdateAchievements)
	e.GET("/getFriends", handler.GetFriends)
	e.GET("/updateFriends", handler.UpdateFriends)
	e.GET("/settings", handler.SettingsPage)
	// e.GET("/updateSteamKey", handler.UpdateSteamKey)
	e.Logger.Fatal(e.Start(":8000"))
}
