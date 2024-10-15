package db

import (
	"database/sql"
	"fmt"
	"go-steam/model"
	"go-steam/util"
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
	}

	db, err := CreateConnection()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var query string

	fmt.Println("Steam Users Games Table...")
	query = `
		CREATE TABLE IF NOT EXISTS steamusergames (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			SteamId INTEGER,
			Appid INTEGER,
			PlaytimeForever INTEGER,
			PlaytimeWindowsForever INTEGER,
			PlaytimeMacForever INTEGER,
			PlaytimeLinuxForever INTEGER,
			PlaytimeDeckForever INTEGER,
			RtimeLastPlayed INTEGER,
			Playtime2Weeks INTEGER,
			LastUpdated TEXT,
			UNIQUE(SteamId, Appid)
		);
		`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}

	fmt.Println("Steam App Details Games Table...")
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

	fmt.Println("Steam Achievements Table...")
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

	fmt.Println("Steam User Achievements Table...")
	query = `
		CREATE TABLE IF NOT EXISTS steamuserachievements (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			Appid INTEGER NOT NULL,
			Apiname TEXT NOT NULL,
			Achieved INTEGER DEFAULT 0,
			Unlocktime INTEGER DEFAULT 0,
			Percentage REAL DEFAULT 0,
			UNIQUE(Appid, Apiname)
		);
		`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}

	fmt.Println("Steam Friends Table...")
	query = `
		CREATE TABLE IF NOT EXISTS steamfriends (
			Userid,
			Steamid,
			Relationship,
			FriendSince,
			UNIQUE(Userid,Steamid)
		);
		`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}

	fmt.Println("Steam Users Table...")
	query = `
		CREATE TABLE IF NOT EXISTS steamusers (
			Steamid,
			Communityvisibilitystate,
			Profilestate,
			Personaname,
			Profileurl,
			Avatar,
			Avatarmedium,
			Avatarfull,
			Avatarhash,
			Lastlogoff,
			Personastate,
			UNIQUE(Steamid)
		);
		`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
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

func GetFilterOption() (model.FilterOptions, error) {
	fmt.Println("Database: GetFilterOption")
	db, err := CreateConnection()
	if err != nil {
		return model.FilterOptions{}, err
	}
	defer db.Close()

	var query = `
	SELECT
	Categories,
	Genres,
	Developers,
	Publishers
	FROM steamappdetails
	`
	rows, err := db.Query(query)
	if err != nil {
		return model.FilterOptions{}, err
	}

	var categoriesList []string
	var categoriesDict = make(map[string]string)

	var genresList []string
	var genresDict = make(map[string]struct{})

	var developersList []string
	var developersDict = make(map[string]string)

	var publishersList []string
	var publishersDict = make(map[string]string)

	for rows.Next() {
		var filter struct {
			Categories string
			Genres     string
			Developers string
			Publishers string
		}
		err = rows.Scan(
			&filter.Categories,
			&filter.Genres,
			&filter.Developers,
			&filter.Publishers,
		)
		if err != nil {
			return model.FilterOptions{}, err
		}

		var categories = strings.Split(filter.Categories, ",")
		var genres = strings.Split(filter.Genres, ",")
		var developers = strings.Split(filter.Developers, ",")
		var publishers = strings.Split(filter.Publishers, ",")

		for _, value := range categories {
			if _, ok := categoriesDict[value]; !ok {
				categoriesDict[value] = ""
				categoriesList = append(categoriesList, value)
			}
		}

		for _, value := range genres {
			if _, ok := genresDict[value]; !ok && value != "" {
				genresDict[value] = struct{}{}
				genresList = append(genresList, value)
			}
		}

		for _, value := range developers {
			if _, ok := developersDict[value]; !ok {
				developersDict[value] = ""
				developersList = append(developersList, value)
			}
		}
		for _, value := range publishers {
			if _, ok := publishersDict[value]; !ok {
				publishersDict[value] = ""
				publishersList = append(publishersList, value)
			}
		}

	}

	var fo model.FilterOptions
	fo.Categories = categoriesList
	fo.Genres = genresList
	fo.Publishers = publishersList
	fo.Developers = developersList
	return fo, nil
}

func GetLibraryDB(title string, genre string) (model.Library, error) {
	fmt.Println("Database: GetLibraryDB")
	db, err := CreateConnection()
	if err != nil {
		return model.Library{}, err
	}
	defer db.Close()

	var query string
	var rows *sql.Rows

	if title != "" {
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
			ORDER BY sug.PlaytimeForever DESC
		`
		rows, err = db.Query(query, title)
	} else if genre != "" {
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
			AND sad.Genres LIKE '%' || ? || '%'
			ORDER BY sug.PlaytimeForever DESC
		`
		rows, err = db.Query(query, genre)
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
			ORDER BY sug.PlaytimeForever DESC
		`
		rows, err = db.Query(query)
	}

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

		card.RtimeLastPlayed = util.StringToTime(card.RtimeLastPlayed)
		library.Cards = append(library.Cards, card)
	}
	return library, nil
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
			IFNULL (sug.PlaytimeForever, 0) AS PlaytimeForever,
			IFNULL (sug.PlaytimeWindowsForever, 0) AS PlaytimeWindowsForever,
			IFNULL (sug.PlaytimeMacForever, 0) AS PlaytimeMacForever,
			IFNULL (sug.PlaytimeLinuxForever, 0) AS PlaytimeLinuxForever,
			IFNULL (sug.RtimeLastPlayed, 0) AS RtimeLastPlayed,
			IFNULL (sug.Playtime2Weeks, 0) AS Playtime2Weeks,
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
		&gameData.AppDetails.PlaytimeForever,
		&gameData.AppDetails.PlaytimeWindowsForever,
		&gameData.AppDetails.PlaytimeMacForever,
		&gameData.AppDetails.PlaytimeLinuxForever,
		&gameData.AppDetails.RtimeLastPlayed,
		&gameData.AppDetails.Playtime2Weeks,
		&gameData.AchivementDetails.LastUpdated,
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

	gameData.AppDetails.PlaytimeForever = (int(gameData.AppDetails.PlaytimeForever / 60))
	gameData.AppDetails.PlaytimeWindowsForever = (int(gameData.AppDetails.PlaytimeWindowsForever / 60))
	gameData.AppDetails.PlaytimeMacForever = (int(gameData.AppDetails.PlaytimeMacForever / 60))
	gameData.AppDetails.Playtime2Weeks = (int(gameData.AppDetails.Playtime2Weeks / 60))
	gameData.AppDetails.PlaytimeLinuxForever = (int(gameData.AppDetails.PlaytimeLinuxForever / 60))

	gameData.AppDetails.RtimeLastPlayed = util.StringToTime(gameData.AppDetails.RtimeLastPlayed)

	query = `
	SELECT
		sa.Name,
		sa.DisplayName,
		sa.Hidden,
		sa.Description,
		sa.Icon,
		sa.IconGray,
		sua.Achieved,
		sua.Unlocktime,
		sua.Percentage
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
			&achievement.Percentage,
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

		gameData.AchivementDetails.Achievements = append(gameData.AchivementDetails.Achievements, achievement)
	}

	return gameData, nil
}

func GetFriendsDB(steamid string) ([]model.Player, error) {
	fmt.Println("Database: GetFriends")
	db, err := CreateConnection()
	if err != nil {
		return []model.Player{}, err
	}
	defer db.Close()

	var query = `
		SELECT
			sf.Steamid,
			sf.FriendSince,
			su.Communityvisibilitystate,
			su.Profilestate,
			su.Personaname,
			su.Profileurl,
			su.Avatar,
			su.Avatarmedium,
			su.Avatarfull,
			su.Lastlogoff,
			su.Personastate
		FROM steamfriends as sf
		JOIN steamusers as su ON su.Steamid = sf.Steamid
		WHERE sf.Userid = ?
	`

	rows, err := db.Query(query, steamid)
	if err != nil {
		return []model.Player{}, err
	}
	defer rows.Close()

	var players []model.Player

	for rows.Next() {
		var player model.Player
		err = rows.Scan(
			&player.Steamid,
			&player.FriendSince,
			&player.Communityvisibilitystate,
			&player.Profilestate,
			&player.Personaname,
			&player.Profileurl,
			&player.Avatar,
			&player.Avatarmedium,
			&player.Avatarfull,
			&player.Lastlogoff,
			&player.Personastate,
		)
		if err != nil {
			return []model.Player{}, err
		}

		player.Lastlogoff = util.StringToTime(player.Lastlogoff)

		players = append(players, player)
	}

	return players, nil
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
	VALUES (?,?,?,?,?,?,?,?,?);
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
			Unlocktime,
			Percentage
		) VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(Appid, Apiname) DO UPDATE SET
			Achieved=excluded.Achieved,
			Unlocktime=excluded.Unlocktime,
			Percentage=excluded.Percentage
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
			achievement.Percentage,
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

func InsertSteamFriendsDB(data []model.FriendAPI, userid string) error {
	fmt.Println("Database: InsertSteamFriendsDB")

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
		INSERT INTO steamfriends (
			Userid,
			Steamid,
			Relationship,
			FriendSince
		)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(Userid, Steamid) DO UPDATE SET
			Relationship=excluded.Relationship,
			FriendSince=excluded.FriendSince
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, friend := range data {
		friendSince := util.StringToTime(util.IntToString(friend.FriendSince))
		_, err = stmt.Exec(
			userid,
			friend.Steamid,
			friend.Relationship,
			friendSince,
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

func InsertSteamUsersDB(data []model.PlayerAPI) error {
	fmt.Println("Database: InsertSteamUsersDB")

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
		INSERT INTO steamusers (
			Steamid,
			Communityvisibilitystate,
			Profilestate,
			Personaname,
			Profileurl,
			Avatar,
			Avatarmedium,
			Avatarfull,
			Avatarhash,
			Lastlogoff,
			Personastate
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(Steamid) DO UPDATE SET
			Communityvisibilitystate=excluded.Communityvisibilitystate,
			Profilestate=excluded.Profilestate,
			Personaname=excluded.Personaname,
			Profileurl=excluded.Profileurl,
			Avatar=excluded.Avatar,
			Avatarmedium=excluded.Avatarmedium,
			Avatarfull=excluded.Avatarfull,
			Avatarhash=excluded.Avatarhash,
			Lastlogoff=excluded.Lastlogoff,
			Personastate=excluded.Personastate
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, user := range data {
		_, err = stmt.Exec(
			user.Steamid,
			user.Communityvisibilitystate,
			user.Profilestate,
			user.Personaname,
			user.Profileurl,
			user.Avatar,
			user.Avatarmedium,
			user.Avatarfull,
			user.Avatarhash,
			user.Lastlogoff,
			user.Personastate,
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

func GetSteamUserAchievements(id int, filter string) ([]model.Achievement, error) {
	fmt.Println("Database: GetSteamUserAchievements")

	db, err := CreateConnection()
	if err != nil {
		return []model.Achievement{}, err
	}
	defer db.Close()
	var query string
	if filter == "All" {
		query = `
		SELECT
			sa.Name,
			sa.DisplayName,
			sa.Hidden,
			sa.Description,
			sa.Icon,
			sa.IconGray,
			sua.Achieved,
			sua.Unlocktime,
			sua.Percentage
		FROM steamachievements as sa
		JOIN steamuserachievements as sua ON sua.Apiname = sa.Name
		WHERE sa.Appid = ?
		`
	} else if filter == "Locked" {
		query = `
		SELECT
			sa.Name,
			sa.DisplayName,
			sa.Hidden,
			sa.Description,
			sa.Icon,
			sa.IconGray,
			sua.Achieved,
			sua.Unlocktime,
			sua.Percentage
		FROM steamachievements as sa
		JOIN steamuserachievements as sua ON sua.Apiname = sa.Name
		WHERE sa.Appid = ?
		AND sua.Achieved == 0
		`
	} else if filter == "Unlocked" {
		query = `
		SELECT
			sa.Name,
			sa.DisplayName,
			sa.Hidden,
			sa.Description,
			sa.Icon,
			sa.IconGray,
			sua.Achieved,
			sua.Unlocktime,
			sua.Percentage
		FROM steamachievements as sa
		JOIN steamuserachievements as sua ON sua.Apiname = sa.Name
		WHERE sa.Appid = ?
		AND sua.Achieved == 1
		`
	}

	rows, err := db.Query(query, id)
	if err != nil {
		return []model.Achievement{}, err
	}
	defer rows.Close()

	var achievements []model.Achievement
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
			&achievement.Percentage,
		)
		if err != nil {
			return []model.Achievement{}, err
		}

		unlockTimeInt, err := strconv.ParseInt(achievement.Unlocktime, 10, 64)
		if err != nil {
			unlockTimeInt = 0
		}

		if unlockTimeInt != 0 {
			achievement.Unlocktime = time.Unix(unlockTimeInt, 0).Format(time.RFC1123)
		}

		achievements = append(achievements, achievement)
	}
	return achievements, nil
}
