package main

import (
	"go-steam/handler"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	l := handler.Library{}
	e.GET("/", l.ShowLibrary)
	e.Logger.Fatal(e.Start(":8000"))
}
