package src

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// Initalizes the Steam database if it does not exist
// Creates steam.db and all tables
func InitDatabase() {

	_, err := os.Stat("steam.db")
	if err != nil {
		fmt.Println("Initializing Steam Database...")

		if _, err = os.Create("steam.db"); err != nil {
			msg := fmt.Sprintf("Error: %v", err)
			fmt.Println(msg)
			panic(err)
		}

		db, err := createConnection()
		if err != nil {
			panic(err)
		}
		defer db.Close()

		var query string

		fmt.Println("Creating Steam Users Games Table...")
		query = `
			CREATE TABLE IF NOT EXISTS steamusergames (
		        id INTEGER PRIMARY KEY AUTOINCREMENT,
				app_id INTEGER,
		        playtime_forever INTEGER,
		        playtime_windows_forever INTEGER,
		        playtime_mac_forever INTEGER,
		        playtime_linux_forever INTEGER,
		        playtime_deck_forever INTEGER,
		        rtime_last_played INTEGER,
				playtime_2weeks INTEGER,
				UNIQUE(app_id)
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
		        Type TEXT,
				Name TEXT,
				SteamAppid INTEGER,
		        RequiredAge
				IsFree INTEGER,
				DetailedDescription TEXT,
				AboutTheGame TEXT,
				ShortDescription TEXT,
				SupportedLanguages TEXT,
				HeaderImage TEXT,
				CapsuleImage TEXT,
				CapsuleImagev5 TEXT,
				Developers BLOB,
				Publishers BLOB,
				Windows INTEGER,
				Mac INTEGER,
				Linux INTEGER,
				Categories BLOB,
				Genres BLOB,
				ReleaseDate TEXT,
				Background TEXT,
				UNIQUE(SteamAppid)
		    );
			`
		_, err = db.Exec(query)
		if err != nil {
			panic(err)
		}
	}
}

// Creates a connector to the sqlite3 database
func createConnection() (*sql.DB, error) {

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

func InsertSteamUserGamesDB(data Games) error {
	msg := fmt.Sprintf("Inserting Appid #%d into SteamUserGames", data.Appid)
	fmt.Println(msg)

	db, err := createConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	var query string = `
	INSERT INTO steamusergames (
	app_id,
	playtime_forever,
	playtime_windows_forever,
	playtime_mac_forever,
	playtime_linux_forever,
	playtime_deck_forever,
	rtime_last_played,
	playtime_2weeks
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
	)
	if err != nil {
		return err
	}
	fmt.Println("Successfull")
	return nil
}

func GetSteamUserGamesDB() ([]Games, error) {
	db, err := createConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var query string = `
	SELECT 
	app_id,
	playtime_forever,
	playtime_windows_forever,
	playtime_mac_forever,
	playtime_linux_forever,
	playtime_deck_forever,
	rtime_last_played,
	playtime_2weeks
	FROM steamusergames;
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []Games
	for rows.Next() {
		var game Games
		err = rows.Scan(
			&game.Appid,
			&game.PlaytimeForever,
			&game.PlaytimeWindowsForever,
			&game.PlaytimeMacForever,
			&game.PlaytimeLinuxForever,
			&game.PlaytimeDeckForever,
			&game.RtimeLastPlayed,
			&game.Playtime2Weeks,
		)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}

	return games, nil
}

func InsertSteamAppDetailsDB(data AppDetails) error {
	db, err := createConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	var query string = `
	INSERT INTO steamappdetails (
	Type,
	Name,
	SteamAppid,
	RequiredAge,
	IsFree,
	DetailedDescription,
	AboutTheGame,
	ShortDescription,
	SupportedLanguages,
	HeaderImage,
	CapsuleImage,
	CapsuleImagev5,
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
	VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);
	`
	_, err = db.Exec(query,
		data.Num240.Data.Type,
		data.Num240.Data.Name,
		data.Num240.Data.SteamAppid,
		data.Num240.Data.RequiredAge,
		data.Num240.Data.IsFree,
		data.Num240.Data.DetailedDescription,
		data.Num240.Data.AboutTheGame,
		data.Num240.Data.ShortDescription,
		data.Num240.Data.SupportedLanguages,
		data.Num240.Data.HeaderImage,
		data.Num240.Data.CapsuleImage,
		data.Num240.Data.CapsuleImagev5,
		data.Num240.Data.Developers,
		data.Num240.Data.Publishers,
		data.Num240.Data.Platforms.Windows,
		data.Num240.Data.Platforms.Mac,
		data.Num240.Data.Platforms.Linux,
		data.Num240.Data.Categories,
		data.Num240.Data.Genres,
		data.Num240.Data.ReleaseDate,
		data.Num240.Data.Background,
	)
	if err != nil {
		return err
	}
	msg := fmt.Sprintf("Inserting Appid #%d", data.Num240.Data.SteamAppid)
	fmt.Println(msg)
	return nil
}
