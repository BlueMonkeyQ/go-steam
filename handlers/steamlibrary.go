package handlers

import (
	"encoding/json"
	"fmt"
	"go-steam/src"
	"io"
	"net/http"

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
		fmt.Println("Error Unmarshaling Json Data into interface")
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
			msg := fmt.Sprintf("Error Inserting Appid #%d into SteamAppDetails table", id)
			fmt.Println(msg)
			c.Logger().Error(err)
			return c.String(http.StatusInternalServerError, msg)
		}
	} else {
		fmt.Println("Getting AppDetails")
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
