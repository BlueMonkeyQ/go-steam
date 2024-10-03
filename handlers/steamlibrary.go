package handlers

import (
	"encoding/json"
	"fmt"
	"go-steam/src"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

// Call the appdetails endpoint
// Because the AppID is used as the key, this key changes values
// Get the data from the key, load into a Json string, then Un load into AppDetails Struct
// If the appid exist in the steamAppDetails table, then update the values; Else insert new record
func GetSteamAppDetail(c echo.Context, id int) error {
	fmt.Println("Endpoint: GetSteamAppDetails")
	url := fmt.Sprintf("http://store.steampowered.com/api/appdetails?appids=%d", id)
	resp, err := http.Get(url)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error fetching AppDetails")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error Reading Response Body")
		return c.String(http.StatusInternalServerError, err.Error())
	}

	var bodyJson map[string]interface{}
	if err := json.Unmarshal(body, &bodyJson); err != nil {
		fmt.Println("Error Unmarshaling Json Data into interface")
		return c.String(http.StatusInternalServerError, err.Error())
	}
	appData, ok := bodyJson[fmt.Sprintf("%d", id)].(map[string]interface{})
	if !ok {
		return c.String(http.StatusInternalServerError, "Error parsing app data")
	}

	jsonData, err := json.Marshal(appData)
	if err != nil {
		fmt.Println("Error Marshling appData into Json")
		return c.String(http.StatusInternalServerError, err.Error())
	}

	var appDetails src.AppDetails
	if err := json.Unmarshal(jsonData, &appDetails); err != nil {
		fmt.Println("Error Unmarshaling Json Data into AppDetails")
		return c.String(http.StatusInternalServerError, err.Error())
	}

	exist, err := src.ExistSteamAppDetailsDBId(id)
	if err != nil {
		msg := fmt.Sprintf("Error Checking if Appid #%d Exist in SteamAppDetails table", id)
		fmt.Println(msg)
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, msg)
	}
	if exist {
		fmt.Println("Updating AppDetails")
		if err := src.UpdateSteamAppDetailsDB(id, appDetails); err != nil {
			msg := fmt.Sprintf("Error Updating Appid #%d into SteamAppDetails table", id)
			fmt.Println(msg)
			c.Logger().Error(err)
			return c.String(http.StatusInternalServerError, msg)
		}
	} else {
		fmt.Println("Inserting AppDetails")
		if err := src.InsertSteamAppDetailsDB(id, appDetails); err != nil {
			msg := fmt.Sprintf("Error Inserting Appid #%d into SteamAppDetails table", id)
			fmt.Println(msg)
			c.Logger().Error(err)
			return c.String(http.StatusInternalServerError, msg)
		}
	}
	fmt.Println("Successfull")
	return nil
}

func GetSteamAppAchievements(c echo.Context, id int) error {
	fmt.Println("Endpoint: GetSteamAppAchievements")
	url := fmt.Sprintf("http://api.steampowered.com/ISteamUserStats/GetSchemaForGame/v0002/?key=14EB214CEC3F1701FD192885D330990F&appid=%d&l=english&format=json", id)
	resp, err := http.Get(url)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error fetching GetSteamAppAchievements")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error Reading Response Body")
		return c.String(http.StatusInternalServerError, err.Error())
	}

	var steamAchievements src.SteamAchievements
	if err := json.Unmarshal(body, &steamAchievements); err != nil {
		fmt.Println("Error Unmarshaling body into SteamAchievements")
		return c.String(http.StatusInternalServerError, err.Error())
	}

	for i, achievement := range steamAchievements.Game.AvailableGameStats.Achievements {
		fmt.Printf("%d Appid #%d\n", i, id)
		err = src.InsertSteamAchievementsDB(id, achievement)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE") {
				msg := fmt.Sprintf("Appid #%d Already Exist", id)
				fmt.Println(msg)
			} else {
				msg := fmt.Sprintf("Error Inserting Appid #%d into SteamAchievements table", id)
				fmt.Println(msg)
				c.Logger().Error(err)
				return c.String(http.StatusInternalServerError, msg)
			}
		}
	}
	return nil
}
