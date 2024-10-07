package services

import (
	"fmt"
	"go-steam/db"
	"go-steam/model"
	"math"
	"strings"
	"time"
)

func GetDetailsPage(id int) (model.GameData, error) {
	data, err := db.GetGameDetailsDB(id)
	if err != nil {
		return model.GameData{}, err
	}
	return data, nil
}

func UpdateAchievements(id int) {
	timestamp := time.Now().Local().Format(time.RFC850)
	err := db.UpdateSteamUserGamesLastUpdated(id, timestamp)
	if err != nil {
		fmt.Printf("Fail: %s \n", err.Error())
	}

	userAchievements, err := GetSteamUserAchievements(id)
	if err != nil {
		if !strings.Contains(err.Error(), "False") {
			fmt.Printf("Fail: %s \n", err.Error())
		}
	}
	globalAchievements, err := GetSteamGlobalAchievements(id)
	if err != nil {
		if !strings.Contains(err.Error(), "False") {
			fmt.Printf("Fail: %s \n", err.Error())
		}
	}

	for i := range userAchievements.Playerstats.Achievements {
		name := userAchievements.Playerstats.Achievements[i].Apiname

		for _, global := range globalAchievements.Achievementpercentages.Achievements {
			if name == global.Name {
				percent := math.Floor(global.Percent*100) / 100
				userAchievements.Playerstats.Achievements[i].Percentage = percent
				continue
			}
		}
	}

	err = db.InsertSteamUserAchievementsDB(userAchievements, id)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			fmt.Println("Warning: Already Exist")
		} else {
			fmt.Printf("Fail: %s", err.Error())
		}
	} else {
		fmt.Println("Pass: Inserted")
	}
}
