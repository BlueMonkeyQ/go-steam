package main

import (
	"go-steam/db"
	"go-steam/handler"

	"github.com/labstack/echo"
)

func main() {
	db.InitDatabase()
	e := echo.New()

	l := handler.Library{}
	e.GET("/", l.ShowLibrary)
	e.GET("/getLibrary", l.GetLibrary)
	e.GET("/getSteamUserLibrary/:AppID", l.ShowGame)
	e.GET("/updateAchivements/:AppID", l.UpdateAchievements)
	e.Logger.Fatal(e.Start(":8000"))
}
