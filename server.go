package main

import (
	"go-steam/handlers"
	"go-steam/src"
	"html/template"
	"io"

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
	e := echo.New()

	// Pre compiled Templates
	t := &Template{
		templates: template.Must(template.ParseGlob("*.html")),
	}
	e.Renderer = t

	e.GET("/updateSteamUserLibrary", handlers.UpdateSteamUserLibrary)
	e.GET("/", handlers.GetSteamUserLibrary)
	e.Logger.Fatal(e.Start(":8000"))
}
