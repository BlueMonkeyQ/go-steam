package handler

import (
	"fmt"
	"go-steam/services"
	"go-steam/views"
	"net/http"

	"github.com/labstack/echo"
)

func GetDetailsPage(c echo.Context) error {
	param := getParam(c, "AppID")
	id, err := stringToInt(param)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	fmt.Printf("Endpoint: GetDetailsPage: %d \n", id)

	data, err := services.GetDetailsPage(id)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	return render(c, views.DetailPageBase(data))
}

func UpdateAchievements(c echo.Context) error {
	param := getParam(c, "AppID")
	id, err := stringToInt(param)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	fmt.Printf("Endpoint: UpdateAchievements: %d \n", id)

	services.UpdateAchievements(id)

	data, err := services.GetDetailsPage(id)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusNotFound, "Unbale to get GetDetailsPage")
	}
	return render(c, views.AchievementTable(data.Achievements))
}
