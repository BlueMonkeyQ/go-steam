package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-steam/db"
	"go-steam/model"
	"io"
	"net/http"
	"strings"
	"time"
)

func GetLibrary(filter string) model.Library {
	data, err := db.GetLibraryDB(filter)
	if err != nil {
		fmt.Println(err)
		return model.Library{}
	}
	return data
}

func UpdateLibrary(data model.GetOwnedGamesAPI) {
	timestamp := time.Now().Local().Format(time.RFC850)
	for i, game := range data.Response.Games {
		fmt.Printf("#%d AppID: %d \n", i, game.Appid)

		// Steam User Game
		exist, err := db.GetSteamUserGamesDB(game.Appid)
		if err != nil {
			fmt.Printf("Fail: %s \n", err.Error())
			break
		}
		if !exist {
			err := db.InsertSteamUserGamesDB(game, timestamp)
			if err != nil {
				if strings.Contains(err.Error(), "UNIQUE") {
					fmt.Println("Warning: Already Exist")
				} else {
					fmt.Printf("Fail: %s", err.Error())
					break
				}
			} else {
				fmt.Println("Pass: Inserted")
			}

		} else {
			fmt.Printf("Info: %d Already Exist\n", game.Appid)
			continue
		}

		// Steam App Details
		exist, err = db.GetSteamAppDetailsAppidDB(game.Appid)
		if err != nil {
			fmt.Printf("Fail: %s \n", err.Error())
			break
		}

		if !exist {
			appDetails, err := GetSteamAppDetail(game.Appid)
			if err != nil {
				if !strings.Contains(err.Error(), "False") {
					fmt.Printf("Fail: %s \n", err.Error())
					break
				}
			}

			err = db.InsertSteamAppDetailsDB(appDetails, game.Appid)
			if err != nil {
				if strings.Contains(err.Error(), "UNIQUE") {
					fmt.Println("Warning: Already Exist")
				} else {
					fmt.Printf("Fail: %s", err.Error())
					break
				}
				continue
			} else {
				fmt.Println("Pass: Inserted")
			}
		} else {
			fmt.Printf("Info: %d Already Exist\n", game.Appid)
		}

		// Steam Achievements
		exist, err = db.GetSteamAchievementsAppidDB(game.Appid)
		if err != nil {
			fmt.Printf("Fail: %s \n", err.Error())
			break
		}

		if !exist {
			achievements, err := GetSteamAchievements(game.Appid)
			if err != nil {
				if !strings.Contains(err.Error(), "False") {
					fmt.Printf("Fail: %s \n", err.Error())
					break
				}
			}

			err = db.InsertSteamAchievementsDB(achievements, game.Appid)
			if err != nil {
				if strings.Contains(err.Error(), "UNIQUE") {
					fmt.Println("Warning: Already Exist")
				} else {
					fmt.Printf("Fail: %s", err.Error())
					break
				}
				continue
			} else {
				fmt.Println("Pass: Inserted")
			}
		} else {
			fmt.Printf("Info: %d Already Exist\n", game.Appid)
		}

		// Steam User Achievements
		exist, err = db.ExistSteamUserAchievementsAppidDB(game.Appid)
		if err != nil {
			fmt.Printf("Fail: %s \n", err.Error())
			break
		}

		if !exist {
			err = db.UpdateSteamUserGamesLastUpdated(game.Appid, timestamp)
			if err != nil {
				fmt.Printf("Fail: %s \n", err.Error())
			}
			userAchievements, err := GetSteamUserAchievements(game.Appid)
			if err != nil {
				if !strings.Contains(err.Error(), "False") {
					fmt.Printf("Fail: %s \n", err.Error())
					break
				}
			}

			err = db.InsertSteamUserAchievementsDB(userAchievements, game.Appid)
			if err != nil {
				if strings.Contains(err.Error(), "UNIQUE") {
					fmt.Println("Warning: Already Exist")
				} else {
					fmt.Printf("Fail: %s", err.Error())
					break
				}
				continue
			} else {
				fmt.Println("Pass: Inserted")
			}
		} else {
			fmt.Printf("Info: %d Already Exist\n", game.Appid)
		}
	}
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

func GetSteamGlobalAchievements(id int) (model.GlobalAchievementsAPI, error) {
	fmt.Println("Endpoint: GetSteamUserAchievements")
	url := fmt.Sprintf("https://api.steampowered.com/ISteamUserStats/GetGlobalAchievementPercentagesForApp/v0002/?gameid=%d&format=json", id)
	resp, err := http.Get(url)
	if err != nil {
		return model.GlobalAchievementsAPI{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.GlobalAchievementsAPI{}, err
	}

	var achievements model.GlobalAchievementsAPI
	if err := json.Unmarshal(body, &achievements); err != nil {
		return model.GlobalAchievementsAPI{}, err
	}
	return achievements, nil

}
