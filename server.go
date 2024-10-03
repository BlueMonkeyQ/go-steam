package main

import (
	"go-steam/handlers"
	"go-steam/src"
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	src.InitDatabase()
	src.InitSteamDatabase()

	e := echo.New()

	// Pre compiled Templates
	t := &Template{
		templates: template.Must(template.ParseGlob("*.html")),
	}
	e.Renderer = t
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", nil)
	})
	e.GET("/updateSteamUserLibrary", handlers.UpdateSteamUserLibrary)
	e.GET("/getSteamUserLibrary", handlers.GetSteamUserLibrary)
	e.GET("/getSteamUserLibrary/:AppID", handlers.GetSteamUserLibraryAppid)
	e.GET("/getSteamUserLibraryAchievements/:AppID", handlers.GetSteamUserLibraryAchievements)
	e.GET("/updateSteamUserLibraryAchievements/:AppID", handlers.UpdateSteamUserLibraryAchievements)
	e.Logger.Fatal(e.Start(":8000"))
}
