package handler

import (
	"fmt"
	"go-steam/db"
	"go-steam/services"
	"go-steam/views"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
)

func GetDetailsPage(c echo.Context) error {
	param := c.Param("AppID")
	id, err := strconv.Atoi(param)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusBadRequest, "Invalid AppID")
	}

	fmt.Printf("Endpoint: GetDetailsPage: %d \n", id)

	data, err := db.GetGameDetailsDB(id)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusNotFound, "Unbale to get GetDetailsPage")
	}
	return render(c, views.DetailPageBase(data))
}

func UpdateAchievements(c echo.Context) error {
	param := c.Param("AppID")
	id, err := strconv.Atoi(param)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusBadRequest, "Invalid AppID")
	}
	fmt.Printf("Endpoint: UpdateAchievements: %d \n", id)

	timestamp := time.Now().Local().Format(time.RFC850)
	err = db.UpdateSteamUserGamesLastUpdated(id, timestamp)
	if err != nil {
		fmt.Printf("Fail: %s \n", err.Error())
	}

	userAchievements, err := services.GetSteamUserAchievements(id)
	if err != nil {
		if !strings.Contains(err.Error(), "False") {
			fmt.Printf("Fail: %s \n", err.Error())
		}
	}

	err = db.InsertSteamUserAchievementsDB(userAchievements, id)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			fmt.Println("Warning: Already Exist")
		} else {
			fmt.Printf("Fail: %s", err.Error())
		}
	} else {
		fmt.Println("Pass: Inserted")
	}

	data, err := db.GetGameDetailsDB(id)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusNotFound, "Unbale to get GetDetailsPage")
	}

	data.Achievements.LastUpdated = timestamp
	return render(c, views.AchievementTable(data.Achievements))

	// Causes computer to brick do not use
	// data, err := db.GetSteamUserAchievementsAppidDB(id)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// data.LastUpdated = timestamp
	// return render(c, views.AchievementTable(data))
}
