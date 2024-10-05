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

func InitDatabase() {

	_, err := os.Stat("steam.db")
	if err != nil {
		fmt.Println("Initializing Steam Database...")

		if _, err = os.Create("steam.db"); err != nil {
			msg := fmt.Sprintf("Error: %v", err)
			fmt.Println(msg)
			panic(err)
		}

		db, err := CreateConnection()
		if err != nil {
			panic(err)
		}
		defer db.Close()

		var query string

		fmt.Println("Creating Steam Users Games Table...")
		query = `
			CREATE TABLE IF NOT EXISTS steamusergames (
		        id INTEGER PRIMARY KEY AUTOINCREMENT,
				Appid INTEGER,
		        PlaytimeForever INTEGER,
		        PlaytimeWindowsForever INTEGER,
		        PlaytimeMacForever INTEGER,
		        PlaytimeLinuxForever INTEGER,
		        PlaytimeDeckForever INTEGER,
		        RtimeLastPlayed INTEGER,
				Playtime2Weeks INTEGER,
				LastUpdated TEXT,
				UNIQUE(Appid)
		    );
			`
		_, err = db.Exec(query)
		if err != nil {
			panic(err)
		}

		fmt.Println("Creating Steam App Details Games Table...")
		query = `
			CREATE TABLE IF NOT EXISTS steamappdetails (
		        id INTEGER PRIMARY KEY AUTOINCREMENT,
				Appid INTEGER,
		        Type TEXT,
				Name TEXT,
				SteamAppid INTEGER,
		        RequiredAge,
				IsFree INTEGER,
				DetailedDescription TEXT,
				AboutTheGame TEXT,
				ShortDescription TEXT,
				SupportedLanguages TEXT,
				HeaderImage TEXT,
				CapsuleImage TEXT,
				CapsuleImagev5 TEXT,
				Developers TEXT,
				Publishers TEXT,
				Windows INTEGER,
				Mac INTEGER,
				Linux INTEGER,
				Categories TEXT,
				Genres TEXT,
				ReleaseDate TEXT,
				Background TEXT,
				UNIQUE(Appid, SteamAppid)
		    );
			`
		_, err = db.Exec(query)
		if err != nil {
			panic(err)
		}

		fmt.Println("Creating Steam Achievements Table...")
		query = `
			CREATE TABLE IF NOT EXISTS steamachievements (
		        id INTEGER PRIMARY KEY AUTOINCREMENT,
		        Appid INTEGER,
		        Name TEXT,
				DisplayName TEXT,
				Hidden INTEGER,
				Description TEXT,
				Icon TEXT,
				IconGray TEXT,
				UNIQUE(Appid, Name)
		    );
			`
		_, err = db.Exec(query)
		if err != nil {
			panic(err)
		}

		fmt.Println("Creating Steam User Achievements Table...")
		query = `
			CREATE TABLE IF NOT EXISTS steamuserachievements (
		        id INTEGER PRIMARY KEY AUTOINCREMENT,
		        Appid INTEGER,
		        Apiname TEXT,
		        Achieved INTEGER,
		        Unlocktime INTEGER,
				UNIQUE(Appid, Apiname)
		    );
			`
		_, err = db.Exec(query)
		if err != nil {
			panic(err)
		}
	}
}

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

func GetSteamAppDetailsAppidDB(id int) (bool, error) {
	fmt.Println("Database: GetSteamAppDetailsAppidDB")
	db, err := CreateConnection()
	if err != nil {
		return false, err
	}
	defer db.Close()

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM steamappdetails WHERE SteamAppid = ?)`
	err = db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		if err != sql.ErrNoRows {
			return false, err
		}
		return false, nil
	}
	return exists, nil
}

func GetSteamAchievementsAppidDB(id int) (bool, error) {
	fmt.Println("Database: GetSteamAchievementsAppidDB")
	db, err := CreateConnection()
	if err != nil {
		return false, err
	}
	defer db.Close()

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM steamachievements WHERE Appid = ?)`
	err = db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		if err != sql.ErrNoRows {
			return false, err
		}
		return false, nil
	}
	return exists, nil
}

func ExistSteamUserAchievementsAppidDB(id int) (bool, error) {
	fmt.Println("Database: GetSteamUserAchievementsAppidDB")
	db, err := CreateConnection()
	if err != nil {
		return false, err
	}
	defer db.Close()

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM steamuserachievements WHERE Appid = ?)`
	err = db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		if err != sql.ErrNoRows {
			return false, err
		}
		return false, nil
	}
	return exists, nil
}

func GetSteamUserGamesDB(id int) (bool, error) {
	fmt.Println("Database: GetSteamUserGamesDB")
	db, err := CreateConnection()
	if err != nil {
		return false, err
	}
	defer db.Close()

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM steamusergames WHERE Appid = ?)`
	err = db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		if err != sql.ErrNoRows {
			return false, err
		}
		return false, nil
	}
	return exists, nil
}

func GetLibraryDB(filter string) (model.Library, error) {
	fmt.Println("Database: GetLibraryDB")
	db, err := CreateConnection()
	if err != nil {
		return model.Library{}, err
	}
	defer db.Close()

	var query string
	if filter == "" {
		query = `
		SELECT 
		sad.Appid, 
		sad.Name, 
		IFNULL(sad.HeaderImage, '') AS HeaderImage,
		sug.RtimeLastPlayed,
		(SELECT COUNT(*) FROM steamachievements WHERE Appid = sad.Appid) AS TotalAchivements,
		(SELECT COUNT(*) FROM steamuserachievements WHERE Appid = sad.Appid AND Achieved = 1) AS TotalAchieved
		FROM steamappdetails AS sad
		JOIN steamusergames AS sug ON sug.Appid = sad.Appid
		WHERE sad.Appid IS NOT NULL 
		AND sad.Appid != 0 
	`
	} else {
		query = `
		SELECT 
		sad.Appid, 
		sad.Name, 
		IFNULL(sad.HeaderImage, '') AS HeaderImage,
		sug.RtimeLastPlayed,
		(SELECT COUNT(*) FROM steamachievements WHERE Appid = sad.Appid) AS TotalAchivements,
		(SELECT COUNT(*) FROM steamuserachievements WHERE Appid = sad.Appid AND Achieved = 1) AS TotalAchieved
		FROM steamappdetails AS sad
		JOIN steamusergames AS sug ON sug.Appid = sad.Appid
		WHERE sad.Appid IS NOT NULL 
		AND sad.Appid != 0
		AND sad.Name LIKE '%' || ? || '%'
	`

	}
	rows, err := db.Query(query, filter)
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
			&card.HeaderImage,
			&card.RtimeLastPlayed,
			&card.TotalAchivements,
			&card.TotalAchieved,
		)
		if err != nil {
			return model.Library{}, err
		}

		unlockTimeInt, err := strconv.ParseInt(card.RtimeLastPlayed, 10, 64)
		if err != nil {
			unlockTimeInt = 0
		}

		if unlockTimeInt != 0 {
			card.RtimeLastPlayed = time.Unix(unlockTimeInt, 0).Format(time.RFC1123)
		}

		library.Cards = append(library.Cards, card)
	}
	return library, nil
}

// Causes computer to brick, do not use
func GetSteamUserAchievementsAppidDB(id int) (model.AchivementDetails, error) {
	fmt.Println("Database: GetSteamUserAchievementsAppidDB")
	db, err := CreateConnection()
	if err != nil {
		return model.AchivementDetails{}, err
	}
	defer db.Close()

	query := `
	SELECT
	sa.Name,
	sa.DisplayName,
	sa.Hidden,
	sa.Description,
	sa.Icon,
	sa.IconGray,
	sua.Achieved,
	sua.UnlockTime
	FROM steamachievements AS sa
	JOIN steamuserachievements AS sua ON sua.Appid = sua.Appid
	WHERE sa.Appid = ?`
	rows, err := db.Query(query, id)
	if err != nil {
		return model.AchivementDetails{}, err
	}
	defer rows.Close()

	var achievementDetails model.AchivementDetails

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
			return model.AchivementDetails{}, err
		}

		unlockTimeInt, err := strconv.ParseInt(achievement.Unlocktime, 10, 64)
		if err != nil {
			unlockTimeInt = 0
		}

		if unlockTimeInt != 0 {
			achievement.Unlocktime = time.Unix(unlockTimeInt, 0).Format(time.RFC1123)
		}

		achievementDetails.Achievements = append(achievementDetails.Achievements, achievement)
	}
	return achievementDetails, nil
}

func GetGameDetailsDB(id int) (model.GameData, error) {
	fmt.Println("Database: GetGameDetailsDB")
	db, err := CreateConnection()
	if err != nil {
		return model.GameData{}, err
	}
	defer db.Close()

	var query = `
		SELECT 
			sad.Appid,
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
			sad.Background,
			IFNULL(sug.LastUpdated, "") AS LastUpdated
		FROM steamappdetails as sad
		JOIN steamusergames as sug ON sug.Appid = sad.Appid
		WHERE sad.Appid = ?
	`

	var gameData model.GameData
	var developers string
	var publishers string
	var categories string
	var genres string
	err = db.QueryRow(query, id).Scan(
		&gameData.AppDetails.AppID,
		&gameData.AppDetails.SteamAppID,
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
		&gameData.Achievements.LastUpdated,
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

		gameData.Achievements.Achievements = append(gameData.Achievements.Achievements, achievement)
	}

	return gameData, nil
}

func InsertSteamUserGamesDB(data model.Games, lastUpdated string) error {
	fmt.Println("Database: InsertSteamUserGamesDB")

	db, err := CreateConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	var query string = `
	INSERT INTO steamusergames (
	Appid,
	PlaytimeForever,
	PlaytimeWindowsForever,
	PlaytimeMacForever,
	PlaytimeLinuxForever,
	PlaytimeDeckForever,
	RtimeLastPlayed,
	Playtime2Weeks,
	LastUpdated
	)
	VALUES (?,?,?,?,?,?,?,?);
	`
	_, err = db.Exec(query,
		data.Appid,
		data.PlaytimeForever,
		data.PlaytimeWindowsForever,
		data.PlaytimeMacForever,
		data.PlaytimeLinuxForever,
		data.PlaytimeDeckForever,
		data.RtimeLastPlayed,
		data.Playtime2Weeks,
		lastUpdated,
	)
	if err != nil {
		return err
	}
	return nil
}

func InsertSteamAchievementsDB(data model.AchievementsApi, id int) error {
	fmt.Println("Database: InsertSteamAchievementsDB")

	db, err := CreateConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
		INSERT INTO steamachievements (
			Appid,
			Name,
			DisplayName,
			Hidden,
			Description,
			Icon,
			IconGray
		) VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(Appid, Name) DO UPDATE SET
			DisplayName=excluded.DisplayName,
			Hidden=excluded.Hidden,
			Description=excluded.Description,
			Icon=excluded.Icon,
			IconGray=excluded.IconGray
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, achievement := range data.Game.AvailableGameStats.Achievements {
		_, err = stmt.Exec(
			id,
			achievement.Name,
			achievement.DisplayName,
			achievement.Hidden,
			achievement.Description,
			achievement.Icon,
			achievement.IconGray,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func InsertSteamUserAchievementsDB(data model.UserAchievements, id int) error {
	fmt.Println("Database: InsertSteamUserAchievementsDB")

	db, err := CreateConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
		INSERT INTO steamuserachievements (
			Appid,
			Apiname,
			Achieved,
			Unlocktime
		) VALUES (?, ?, ?, ?)
		ON CONFLICT(Appid, Apiname) DO UPDATE SET
			Achieved=excluded.Achieved,
			Unlocktime=excluded.Unlocktime
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, achievement := range data.Playerstats.Achievements {
		_, err = stmt.Exec(
			id,
			achievement.Apiname,
			achievement.Achieved,
			achievement.Unlocktime,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func InsertSteamAppDetailsDB(data model.AppDetailsAPI, id int) error {
	fmt.Println("Database: InsertSteamAppDetailsDB")
	fmt.Println(data.SteamAppid)

	db, err := CreateConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	var query string = `
		INSERT INTO steamappdetails (
			Appid,
			SteamAppid,
			Name,
			AboutTheGame,
			ShortDescription,
			SupportedLanguages,
			HeaderImage,
			Developers,
			Publishers,
			Windows,
			Mac,
			Linux,
			Categories,
			Genres,
			ReleaseDate,
			Background
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`

	var categories []string
	for _, i := range data.Categories {
		categories = append(categories, i.Description)
	}
	var genres []string
	for _, i := range data.Genres {
		genres = append(genres, i.Description)
	}
	_, err = db.Exec(query,
		id,
		data.SteamAppid,
		data.Name,
		data.AboutTheGame,
		data.ShortDescription,
		data.SupportedLanguages,
		data.HeaderImage,
		strings.Join(data.Developers, ","),
		strings.Join(data.Publishers, ","),
		data.Platforms.Windows,
		data.Platforms.Mac,
		data.Platforms.Linux,
		strings.Join(categories, ","),
		strings.Join(genres, ","),
		data.ReleaseDate.Date,
		data.Background,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func UpdateSteamUserGamesLastUpdated(id int, lastUpdated string) error {
	fmt.Println("Database: UpdateSteamUserGamesLastUpdated")
	db, err := CreateConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	var query string = `
		UPDATE steamusergames 
		SET
		LastUpdated = ?
		WHERE Appid = ?
	`

	_, err = db.Exec(query, lastUpdated, id)
	if err != nil {
		return err
	}
	return nil
}
