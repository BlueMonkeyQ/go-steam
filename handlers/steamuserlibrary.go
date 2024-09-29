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
		break

	}
	return nil
}

func GetSteamUserLibrary(c echo.Context) error {
	fmt.Println("Endpoint: GetSteamUserLibrary")
	games, err := src.GetSteamUserGamesDB()
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, "Error fetching games")
	}
	fmt.Printf("Returning #%d Games\n", len(games))
	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"games": games,
	})
}
