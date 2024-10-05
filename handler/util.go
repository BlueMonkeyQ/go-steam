package handler

import (
	"strconv"

	"github.com/a-h/templ"
	"github.com/labstack/echo"
)

func render(c echo.Context, component templ.Component) error {
	return component.Render(c.Request().Context(), c.Response())
}

func getParam(c echo.Context, param string) string {
	value := c.Param(param)
	return value

}

func stringToInt(value string) (int, error) {
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return intValue, nil
}
