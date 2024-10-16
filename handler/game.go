package handler

import (
	"fmt"
	"go-steam/model"
	"go-steam/services"
	"go-steam/util"
	"go-steam/views"
	"net/http"

	"github.com/labstack/echo"
)

func GetDetailsPage(c echo.Context) error {
	param := util.GetParam(c, "AppID")
	id, err := util.StringToInt(param)
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

	data.AchivementDetails.AppID = util.IntToString(id)
	return util.Render(c, views.DetailPageBase(data))
}

func GetAchievements(c echo.Context) error {
	param := util.GetParam(c, "AppID")
	id, err := util.StringToInt(param)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	param = util.GetParam(c, "Filter")
	data, err := services.GetAchievements(id, param)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	return util.Render(c, views.AchievementTable(data))
}

func UpdateAchievements(c echo.Context) error {
	param := util.GetParam(c, "AppID")
	id, err := util.StringToInt(param)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	fmt.Printf("Endpoint: UpdateAchievements: %d \n", id)

	if err := services.UpdateAchievements(id); err != nil {
		c.Logger().Error(err)
		return util.Render(c, views.AchievementTable([]model.Achievement{}))
	}

	data, err := services.GetDetailsPage(id)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusNotFound, "Unbale to get GetDetailsPage")
	}
	return util.Render(c, views.AchievementTable(data.AchivementDetails.Achievements))
}

// func FilterAchievements(c echo.Context) error {
// 	param := util.GetParam(c, "filter")
// 	switch param {
// 	case "All":

// 	}
// }
