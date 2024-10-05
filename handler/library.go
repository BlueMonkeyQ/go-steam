package handler

import (
	"fmt"
	"go-steam/db"
	"go-steam/services"
	"go-steam/views"
	"strings"
	"time"

	"github.com/labstack/echo"
)

func GetLibraryFiltered(c echo.Context) error {
	title := c.QueryParam("filter")
	fmt.Printf("Endpoint: ShowLibraryFiltered: %s \n", title)

	data := services.GetLibrary(title)
	fmt.Printf("Returning #%d Games \n", len(data.Cards))
	return render(c, views.LibraryCards(data))
}

func GetLibrary(c echo.Context) error {
	fmt.Println("Endpoint: GetLibrary")

	data := services.GetLibrary("")
	fmt.Printf("Returning #%d Games \n", len(data.Cards))
	return render(c, views.ShowLibrary(data))
}

func UpdateLibrary(c echo.Context) error {
	fmt.Println("Endpoint: UpdateLibrary")
	getOwnedGames, err := services.GetSteamUserGames()
	if err != nil {
		fmt.Printf("Fail: %s", err.Error())
	}
	fmt.Printf("#%d Games\n", getOwnedGames.Response.GameCount)

	timestamp := time.Now().Local().Format(time.RFC850)
	for i, game := range getOwnedGames.Response.Games {
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
	data := services.GetLibrary("")
	fmt.Printf("Returning #%d Games \n", len(data.Cards))
	return render(c, views.LibraryCards(data))
}
