package handlers

import (
	"encoding/json"
	"fmt"
	"go-steam/src"
	"io"
	"net/http"

	"github.com/labstack/echo"
)

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
	var appDetails src.AppDetails
	if err := json.Unmarshal(body, &appDetails); err != nil {
		fmt.Println("Error Unmarshaling Json Data")
		return c.String(http.StatusInternalServerError, err.Error())
	}
	fmt.Println(appDetails)

	exist, err := src.ExistSteamAppDetailsDBId(id)
	if err != nil {
		msg := fmt.Sprintf("Error Checking if Appid #%d Exist in SteamAppDetails table", id)
		fmt.Println(msg)
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, msg)
	}
	if exist {
		fmt.Println("Updating AppDetails")
		if err := src.UpdateSteamAppDetailsDB(appDetails); err != nil {
			msg := fmt.Sprintf("Error Inserting Appid #%d into SteamAppDetails table", id)
			fmt.Println(msg)
			c.Logger().Error(err)
			return c.String(http.StatusInternalServerError, msg)
		}
	} else {
		fmt.Println("Getting AppDetails")
		if err := src.InsertSteamAppDetailsDB(appDetails); err != nil {
			msg := fmt.Sprintf("Error Inserting Appid #%d into SteamAppDetails table", id)
			fmt.Println(msg)
			c.Logger().Error(err)
			return c.String(http.StatusInternalServerError, msg)
		}
	}
	fmt.Println("Successfull")
	return nil
}
