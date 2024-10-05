package util

import (
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/labstack/echo"
)

func Render(c echo.Context, component templ.Component) error {
	return component.Render(c.Request().Context(), c.Response())
}

func GetParam(c echo.Context, param string) string {
	value := c.Param(param)
	return value

}

func StringToInt(value string) (int, error) {
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return intValue, nil
}

func IntToString(value int) string {
	stringValue := strconv.Itoa(value)
	return stringValue
}

func StringToTime(value string) string {
	unlockTimeInt, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		unlockTimeInt = 0
	}
	return time.Unix(unlockTimeInt, 0).Format(time.RFC1123)
}
