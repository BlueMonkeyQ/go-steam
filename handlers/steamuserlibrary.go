package handlers

import (
	"encoding/json"
	"fmt"
	"go-steam/src"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
)

func UpdateSteamUserLibrary(c echo.Context) error {
	fmt.Println("Endpoint: UpdateSteamUserLibrary")
	resp, err := http.Get("http://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key=14EB214CEC3F1701FD192885D330990F&steamid=76561198050437739&format=json")
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, "Error fetching GetOwnedGames")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, "Error Reading Response Body")
	}

	var sul src.SteamUserLibrary
	if err := json.Unmarshal(body, &sul); err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, "Error Unmarshaling Json Data")
	}

	fmt.Printf("#%d Games\n", sul.Response.GameCount)

	for i, game := range sul.Response.Games {
		fmt.Printf("%d Appid #%d\n", i, game.Appid)
		err = src.InsertSteamUserGamesDB(game)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE") {
				msg := fmt.Sprintf("Appid #%d Already Exist", game.Appid)
				fmt.Println(msg)
			} else {
				msg := fmt.Sprintf("Error Inserting Appid #%d into SteamUserGames table", game.Appid)
				fmt.Println(msg)
				c.Logger().Error(err)
				return c.String(http.StatusInternalServerError, msg)
			}
		}
		if err := GetSteamAppDetail(c, game.Appid); err != nil {
			c.Logger().Error(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		if err := GetSteamAppAchievements(c, game.Appid); err != nil {
			c.Logger().Error(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
		break

	}
	return nil
}

func UpdateSteamUserLibraryAchievements(c echo.Context) error {
	param := c.Param("AppID")
	fmt.Printf("Endpoint: UpdateSteamUserLibraryAchievements/%s\n", param)

	appId, err := strconv.Atoi(param)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusBadRequest, "Invalid AppID")
	}

	url := fmt.Sprintf("http://api.steampowered.com/ISteamUserStats/GetPlayerAchievements/v0001/?appid=%d&key=14EB214CEC3F1701FD192885D330990F&steamid=76561198050437739", appId)
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

	var sua src.SteamUserAchivements
	if err := json.Unmarshal(body, &sua); err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, "Error Unmarshaling Json Data")
	}

	for i, achievement := range sua.Playerstats.Achievements {
		fmt.Printf("%d Appid #%d & Apiname: %s\n", i, appId, achievement.Apiname)
		exist, err := src.ExistSteamUserAchivementsDB(appId, achievement.Apiname)
		if err != nil {
			msg := fmt.Sprintf("Error Checking if Appid: #%d & Apiname: %s Exist in SteamUserAchievements table", appId, achievement.Apiname)
			fmt.Println(msg)
			c.Logger().Error(err)
			return c.String(http.StatusInternalServerError, msg)
		}
		if exist {
			fmt.Println("Updating SteamUserAchivement")
			if err := src.UpdateSteamUserAchivementsDB(appId, achievement); err != nil {
				msg := fmt.Sprintf("Error Updating Appid #%d & Apiname: %s into SteamUserAchievements table", appId, achievement.Apiname)
				fmt.Println(msg)
				c.Logger().Error(err)
				return c.String(http.StatusInternalServerError, msg)
			}
		} else {
			fmt.Println("Inserting SteamUserAchivement")
			if err := src.InsertSteamUserAchivementsDB(appId, achievement); err != nil {
				msg := fmt.Sprintf("Error Inserting Appid #%d & Apiname: %s into SteamUserAchievements table", appId, achievement.Apiname)
				fmt.Println(msg)
				c.Logger().Error(err)
				return c.String(http.StatusInternalServerError, msg)
			}
		}

	}
	return nil

}

func GetSteamUserLibrary(c echo.Context) error {
	fmt.Println("Endpoint: GetSteamUserLibrary")
	library, err := src.GetSteamUserLibrary()
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, "Error fetching Library")
	}
	fmt.Printf("Returning #%d Games\n", len(library))
	return c.Render(http.StatusOK, "getSteamUserLibrary.html", map[string]interface{}{
		"library": library,
	})
}

func GetSteamUserLibraryAppid(c echo.Context) error {
	param := c.Param("AppID")
	fmt.Printf("Endpoint: GetSteamUserLibrary/%s\n", param)

	appId, err := strconv.Atoi(param)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusBadRequest, "Invalid AppID")
	}
	game, err := src.GetSteamUserLibraryAppid(appId)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, "Error fetching Game")
	}

	fmt.Println(game)

	fmt.Printf("Returning Library Game #%d\n", appId)
	return c.Render(http.StatusOK, "getSteamUserLibraryAppid.html", map[string]interface{}{
		"game": game,
	})
}

func GetSteamUserLibraryAchievements(c echo.Context) error {
	param := c.Param("AppID")
	fmt.Printf("Endpoint: GetSteamUserLibrary/%s\n", param)

	appId, err := strconv.Atoi(param)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusBadRequest, "Invalid AppID")
	}
	fmt.Println("Endpoint: GetSteamUserLibraryAchievements")
	achivements, err := src.GetSteamUserLibraryAchievements(appId)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, "Error fetching achivements")
	}
	fmt.Printf("Returning #%d Achivements\n", len(achivements))
	return c.Render(http.StatusOK, "getSteamUserLibraryAppid.html", map[string]interface{}{
		"achivements": achivements,
	})
}
