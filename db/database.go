package db

import (
	"database/sql"
	"fmt"
	"go-steam/model"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Creates a connector to the Steam Sqlite3 Database
func CreateConnection() (*sql.DB, error) {

	_, err := os.Stat("steam.db")
	if err != nil {
		msg := fmt.Sprintf("Warning: %v", err)
		fmt.Println(msg)
		return nil, err
	}

	connector, err := sql.Open("sqlite3", "steam.db")
	if err != nil {
		msg := fmt.Sprintf("Warning: %v", err)
		fmt.Println(msg)
		return nil, err
	}

	return connector, nil
}

func GetLibraryDB() (model.Library, error) {
	fmt.Println("Database: GetLibraryDB")
	db, err := CreateConnection()
	if err != nil {
		return model.Library{}, err
	}
	defer db.Close()

	var query = `
		SELECT SteamAppid, Name, CapsuleImage
		FROM steamappdetails
	`
	rows, err := db.Query(query)
	if err != nil {
		return model.Library{}, err
	}
	defer rows.Close()

	var library model.Library

	for rows.Next() {
		var card model.LibraryCard
		err = rows.Scan(
			&card.AppID,
			&card.Name,
			&card.CapsuleImage,
		)
		if err != nil {
			return model.Library{}, err
		}
		library.Cards = append(library.Cards, card)
	}
	return library, nil
}

func GetGameDB(id int) (model.GameData, error) {
	fmt.Println("Database: GetGameDB")
	db, err := CreateConnection()
	if err != nil {
		return model.GameData{}, err
	}
	defer db.Close()

	var query = `
		SELECT 
			sad.SteamAppid,
			sad.Name,
			sad.AboutTheGame,
			sad.ShortDescription,
			sad.SupportedLanguages,
			sad.HeaderImage,
			sad.Developers,
			sad.Publishers,
			sad.Windows,
			sad.Mac,
			sad.Linux,
			sad.Categories,
			sad.Genres,
			sad.ReleaseDate,
			sad.Background
		FROM steamappdetails as sad
		WHERE sad.SteamAppid = ?
	`

	var gameData model.GameData
	var developers string
	var publishers string
	var categories string
	var genres string
	err = db.QueryRow(query, id).Scan(
		&gameData.AppDetails.AppID,
		&gameData.AppDetails.Name,
		&gameData.AppDetails.AboutTheGame,
		&gameData.AppDetails.ShortDescripton,
		&gameData.AppDetails.SupportedLanguages,
		&gameData.AppDetails.HeaderImage,
		&developers,
		&publishers,
		&gameData.AppDetails.Windows,
		&gameData.AppDetails.Mac,
		&gameData.AppDetails.Linux,
		&categories,
		&genres,
		&gameData.AppDetails.ReleaseDate,
		&gameData.AppDetails.Background,
	)
	if err != nil {
		if err != sql.ErrNoRows {
			return model.GameData{}, err
		}
	}

	gameData.AppDetails.Developers = strings.Split(developers, ",")
	gameData.AppDetails.Publishers = strings.Split(publishers, ",")
	gameData.AppDetails.Categories = strings.Split(categories, ",")
	gameData.AppDetails.Genres = strings.Split(genres, ",")

	query = `
	SELECT
		sa.Name,
		sa.DisplayName,
		sa.Hidden,
		sa.Description,
		sa.Icon,
		sa.IconGray,
		sua.Achieved,
		sua.Unlocktime
	FROM steamachievements as sa
	JOIN steamuserachievements as sua ON sua.Apiname = sa.Name
	WHERE sa.Appid = ?
	`
	rows, err := db.Query(query, id)
	if err != nil {
		return model.GameData{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var achievement model.Achievement
		err = rows.Scan(
			&achievement.Name,
			&achievement.DisplayName,
			&achievement.Hidden,
			&achievement.Description,
			&achievement.Icon,
			&achievement.IconGray,
			&achievement.Achieved,
			&achievement.Unlocktime,
		)
		if err != nil {
			return model.GameData{}, err
		}

		unlockTimeInt, err := strconv.ParseInt(achievement.Unlocktime, 10, 64)
		if err != nil {
			unlockTimeInt = 0
		}

		if unlockTimeInt != 0 {
			achievement.Unlocktime = time.Unix(unlockTimeInt, 0).Format(time.RFC1123)
		}

		gameData.Achievements = append(gameData.Achievements, achievement)
	}

	return gameData, nil
}
