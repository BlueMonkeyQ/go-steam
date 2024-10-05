package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-steam/db"
	"go-steam/model"
	"io"
	"net/http"
)

func GetLibrary(filter string) model.Library {
	data, err := db.GetLibraryDB(filter)
	if err != nil {
		fmt.Println(err)
		return model.Library{}
	}
	return data
}

func GetSteamUserGames() (model.GetOwnedGamesAPI, error) {
	fmt.Println("Endpoint: GetOwnedGames")
	resp, err := http.Get("http://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key=14EB214CEC3F1701FD192885D330990F&steamid=76561198050437739&format=json")
	if err != nil {
		return model.GetOwnedGamesAPI{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.GetOwnedGamesAPI{}, err
	}

	var getOwnedGames model.GetOwnedGamesAPI
	if err := json.Unmarshal(body, &getOwnedGames); err != nil {
		return model.GetOwnedGamesAPI{}, err
	}
	return getOwnedGames, nil
}

func GetSteamAppDetail(id int) (model.AppDetailsAPI, error) {
	fmt.Println("Endpoint: AppDetails")
	url := fmt.Sprintf("http://store.steampowered.com/api/appdetails?appids=%d", id)
	resp, err := http.Get(url)
	if err != nil {
		return model.AppDetailsAPI{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.AppDetailsAPI{}, err
	}
	var bodyJson map[string]interface{}
	if err := json.Unmarshal(body, &bodyJson); err != nil {
		fmt.Println("Error: Unable to Unmarshal body")
		return model.AppDetailsAPI{}, err
	}

	appData, ok := bodyJson[fmt.Sprintf("%d", id)].(map[string]interface{})
	if !ok {
		fmt.Println("KeyError: Appid")
		return model.AppDetailsAPI{}, errors.New("KeyError: Appid")
	}

	success, ok := appData["success"].(bool)
	if !ok {
		fmt.Println("KeyError: Success")
		return model.AppDetailsAPI{}, errors.New("KeyError: Success")
	}
	if !success {
		fmt.Println("Error: Success is False")
		return model.AppDetailsAPI{}, errors.New("False")
	}

	data, ok := appData["data"].(map[string]interface{})
	if !ok {
		fmt.Println("KeyError: Data")
		return model.AppDetailsAPI{}, errors.New("KeyError: Data")
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error: Unable to Marshal data")
		return model.AppDetailsAPI{}, err
	}

	var appDetails model.AppDetailsAPI
	if err := json.Unmarshal(jsonData, &appDetails); err != nil {
		fmt.Println("Error: Unable to Unmarshal jsonData")
		return model.AppDetailsAPI{}, err
	}

	return appDetails, nil
}

func GetSteamAchievements(id int) (model.AchievementsApi, error) {
	fmt.Println("Endpoint: GetSteamAchievements")
	url := fmt.Sprintf("http://api.steampowered.com/ISteamUserStats/GetSchemaForGame/v0002/?key=14EB214CEC3F1701FD192885D330990F&appid=%d&l=english&format=json", id)
	resp, err := http.Get(url)
	if err != nil {
		return model.AchievementsApi{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.AchievementsApi{}, err
	}

	var achievements model.AchievementsApi
	if err := json.Unmarshal(body, &achievements); err != nil {
		return model.AchievementsApi{}, err
	}
	return achievements, nil

}

func GetSteamUserAchievements(id int) (model.UserAchievements, error) {
	fmt.Println("Endpoint: GetSteamUserAchievements")
	url := fmt.Sprintf("http://api.steampowered.com/ISteamUserStats/GetPlayerAchievements/v0001/?appid=%d&key=14EB214CEC3F1701FD192885D330990F&steamid=76561198050437739", id)
	resp, err := http.Get(url)
	if err != nil {
		return model.UserAchievements{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.UserAchievements{}, err
	}

	var achievements model.UserAchievements
	if err := json.Unmarshal(body, &achievements); err != nil {
		return model.UserAchievements{}, err
	}
	return achievements, nil

}
