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

type Library struct{}

func (l Library) ShowLibrary(c echo.Context) error {
	fmt.Println("Endpoint: ShowLibrary")
	data := services.GetLibrary()
	fmt.Printf("Returning #%d Games \n", len(data.Cards))
	return render(c, views.ShowLibrary(data))
}

func (l Library) GetLibrary(c echo.Context) error {
	fmt.Println("Endpoint: GetLibrary")
	getOwnedGames, err := services.GetSteamUserGames()
	if err != nil {
		fmt.Printf("Fail: %s", err.Error())
	}

	fmt.Printf("#%d Games\n", getOwnedGames.Response.GameCount)
	for i, game := range getOwnedGames.Response.Games {
		// if game.Appid != 10190 {
		// 	continue
		// }
		fmt.Printf("#%d AppID: %d \n", i, game.Appid)

		// Steam User Game
		exist, err := db.GetSteamUserGamesDB(game.Appid)
		if err != nil {
			fmt.Printf("Fail: %s \n", err.Error())
			break
		}
		if !exist {
			err := db.InsertSteamUserGamesDB(game)
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
		}

		// Steam App Details
		exist, err = db.GetSteamAppDetailsAppidDB(game.Appid)
		if err != nil {
			fmt.Printf("Fail: %s \n", err.Error())
			break
		}

		if !exist {
			appDetails, err := services.GetSteamAppDetail(game.Appid)
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
			achievements, err := services.GetSteamAchievements(game.Appid)
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
		exist, err = db.GetSteamUserAchievementsAppidDB(game.Appid)
		if err != nil {
			fmt.Printf("Fail: %s \n", err.Error())
			break
		}

		if !exist {
			userAchievements, err := services.GetSteamUserAchievements(game.Appid)
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
	return nil

}

func (l Library) ShowGame(c echo.Context) error {
	param := c.Param("AppID")
	id, err := strconv.Atoi(param)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusBadRequest, "Invalid AppID")
	}

	data := services.GetGame(id)
	return render(c, views.ShowGame(data))
}

func (l Library) UpdateAchievements(c echo.Context) error {
	param := c.Param("AppID")
	id, err := strconv.Atoi(param)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusBadRequest, "Invalid AppID")
	}

	timestamp := time.Now().Format(time.RFC3339)
	fmt.Printf("Current Timestamp: %s\n", timestamp)

	fmt.Printf("Endpoint: UpdateAchievements; Param: %d \n", id)
	return nil
}
