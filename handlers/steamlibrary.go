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

func GetSteamGameData(c echo.Context, appid int) error {
	url := fmt.Sprintf("http://store.steampowered.com/api/appdetails?appids=%d", appid)
	resp, err := http.Get(url)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, "Error fetching AppDetails")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, "Error Reading Response Body")
	}

	var ap src.AppDetails
	if err := json.Unmarshal(body, &ap); err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, "Error Unmarshaling AppDetails Json Data")
	}

	err = src.InsertSteamAppDetailsDB(ap)
	if strings.Contains(err.Error(), "UNIQUE") {
		msg := fmt.Sprintf("Appid #%d Already Exist", ap.Num240.Data.SteamAppid)
		fmt.Println(msg)
	} else if err != nil {
		c.Logger().Error(err)
		msg := fmt.Sprintf("Error Inserting Appid #%d into SteamAppDetails table", ap.Num240.Data.SteamAppid)
		return c.String(http.StatusInternalServerError, msg)
	}
	return nil

}
