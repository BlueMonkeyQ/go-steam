package src

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

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
				UNIQUE(SteamAppid)
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

func InitSteamDatabase() error {
	db, err := createConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	var query = `
	SELECT id FROM steamappdetails
	`
	var id int
	err = db.QueryRow(query).Scan(&id)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	}
	return nil
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
	fmt.Println("Database: InsertSteamUserGamesDB")

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
	return nil
}

func InsertSteamAchievementsDB(appId int, data Achievements) error {
	fmt.Println("Database: InsertSteamAchievementsDB")

	db, err := createConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	var query string = `
	INSERT INTO steamachievements (
	Appid,
	Name,
	DisplayName,
	Hidden,
	Description,
	Icon,
	IconGray
	)
	VALUES (?,?,?,?,?,?,?);
	`
	_, err = db.Exec(query,
		appId,
		data.Name,
		data.DisplayName,
		data.Hidden,
		data.Description,
		data.Icon,
		data.Icongray,
	)
	if err != nil {
		return err
	}
	return nil
}

func InsertSteamAppDetailsDB(id int, data AppDetails) error {
	fmt.Println("Database: InsertSteamAppDetailsDB")
	msg := fmt.Sprintf("Inserting Appid #%d", id)
	fmt.Println(msg)
	db, err := createConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	var categories []string
	for _, i := range data.Data.Categories {
		categories = append(categories, i.Description)
	}
	var genres []string
	for _, i := range data.Data.Genres {
		genres = append(genres, i.Description)
	}

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
		data.Data.Type,
		data.Data.Name,
		data.Data.SteamAppid,
		data.Data.RequiredAge,
		data.Data.IsFree,
		data.Data.DetailedDescription,
		data.Data.AboutTheGame,
		data.Data.ShortDescription,
		data.Data.SupportedLanguages,
		data.Data.HeaderImage,
		data.Data.CapsuleImage,
		data.Data.CapsuleImagev5,
		strings.Join(data.Data.Developers, ","),
		strings.Join(data.Data.Publishers, ","),
		data.Data.Platforms.Windows,
		data.Data.Platforms.Mac,
		data.Data.Platforms.Linux,
		strings.Join(categories, ","),
		strings.Join(genres, ","),
		data.Data.ReleaseDate.Date,
		data.Data.Background,
	)
	if err != nil {
		return err
	}
	return nil
}

func InsertSteamUserAchivementsDB(id int, data UserAchievements) error {
	fmt.Println("Database: InsertSteamUserAchivementsDB")

	db, err := createConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	var query string = `
	INSERT INTO steamuserachievements (
	Appid,
	Apiname,
	Achieved,
	Unlocktime
	)
	VALUES (?,?,?,?);
	`
	_, err = db.Exec(query,
		id,
		data.Apiname,
		data.Achieved,
		data.Unlocktime,
	)
	if err != nil {
		return err
	}
	return nil
}

func UpdateSteamAppDetailsDB(id int, data AppDetails) error {
	fmt.Println("Database: UpdateSteamAppDetailsDB")
	db, err := createConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	var categories []string
	for _, i := range data.Data.Categories {
		categories = append(categories, i.Description)
	}
	var genres []string
	for _, i := range data.Data.Genres {
		genres = append(genres, i.Description)
	}

	var query string = `
	UPDATE steamappdetails
	SET
		Type = ?,
		Name = ?,
		RequiredAge = ?,
		IsFree = ?,
		DetailedDescription = ?,
		AboutTheGame = ?,
		ShortDescription = ?,
		SupportedLanguages = ?,
		HeaderImage = ?,
		CapsuleImage = ?,
		CapsuleImagev5 = ?,
		Developers = ?,
		Publishers = ?,
		Windows = ?,
		Mac = ?,
		Linux = ?,
		Categories = ?,
		Genres = ?,
		ReleaseDate = ?,
		Background = ?
	WHERE
		SteamAppid = ?
	`
	_, err = db.Exec(query,
		data.Data.Type,
		data.Data.Name,
		data.Data.RequiredAge,
		data.Data.IsFree,
		data.Data.DetailedDescription,
		data.Data.AboutTheGame,
		data.Data.ShortDescription,
		data.Data.SupportedLanguages,
		data.Data.HeaderImage,
		data.Data.CapsuleImage,
		data.Data.CapsuleImagev5,
		strings.Join(data.Data.Developers, ","),
		strings.Join(data.Data.Publishers, ","),
		data.Data.Platforms.Windows,
		data.Data.Platforms.Mac,
		data.Data.Platforms.Linux,
		strings.Join(categories, ","),
		strings.Join(genres, ","),
		data.Data.ReleaseDate.Date,
		data.Data.Background,
		data.Data.SteamAppid,
	)
	if err != nil {
		return err
	}
	return nil
}

func UpdateSteamUserAchivementsDB(id int, data UserAchievements) error {
	fmt.Println("Database: UpdateSteamUserAchivementsDB")
	db, err := createConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	var query string = `
	UPDATE steamuserachievements
	SET
		Achieved = ?,
		Unlocktime = ?
	WHERE
		Appid = ?
	AND
		Apiname = ?
	`
	_, err = db.Exec(query,
		data.Achieved,
		data.Unlocktime,
		id,
		data.Apiname,
	)
	if err != nil {
		return err
	}
	return nil
}

func ExistSteamAppDetailsDBId(id int) (bool, error) {
	fmt.Println("Database: ExistSteamAppDetailsDBId")
	db, err := createConnection()
	if err != nil {
		return false, err
	}
	defer db.Close()

	var query string = `
	SELECT id FROM steamappdetails WHERE SteamAppid = ?;
	`
	var exist int
	err = db.QueryRow(query, id).Scan(&exist)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}

func ExistSteamAchievementsId(id int) (bool, error) {
	fmt.Println("Database: ExistSteamAchievementsId")
	db, err := createConnection()
	if err != nil {
		return false, err
	}
	defer db.Close()

	var query string = `
	SELECT Appid FROM steamuserachievements WHERE Appid = ?;
	`
	var exist int
	err = db.QueryRow(query, id).Scan(&exist)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}

func ExistSteamUserAchivementsDB(id int, apiname string) (bool, error) {
	fmt.Println("Database: ExistSteamUserAchivementsDB")
	db, err := createConnection()
	if err != nil {
		return false, err
	}
	defer db.Close()

	var query string = `
	SELECT Appid FROM steamuserachievements WHERE Appid = ? AND Apiname = ?;
	`
	var exist int
	err = db.QueryRow(query, id, apiname).Scan(&exist)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}

func GetSteamAppDetailsDBId(id int) (AppDetails, error) {
	fmt.Println("Database: GetSteamAppDetailsDBId")
	db, err := createConnection()
	if err != nil {
		return AppDetails{}, err
	}
	defer db.Close()

	var query string = `
	SELECT * FROM steamappdetails WHERE SteamAppid = ?;
	`
	var appDetails AppDetails
	err = db.QueryRow(query, id).Scan(&appDetails)
	if err != nil {
		return AppDetails{}, err
	}

	return appDetails, nil
}

func GetSteamUserGamesDB() ([]Games, error) {
	fmt.Println("Database: GetSteamUserGamesDB")
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

func GetSteamUserLibrary() ([]Library, error) {
	fmt.Println("Database: GetSteamUserGamesDB")
	db, err := createConnection()
	if err != nil {
		return []Library{}, err
	}
	defer db.Close()

	var query string = `
	SELECT 
		steamusergames.app_id,
		steamappdetails.Type,
		steamappdetails.Name,
		steamappdetails.RequiredAge,
		steamappdetails.IsFree,
		steamappdetails.DetailedDescription,
		steamappdetails.AboutTheGame,
		steamappdetails.ShortDescription,
		steamappdetails.SupportedLanguages,
		steamappdetails.HeaderImage,
		steamappdetails.CapsuleImage,
		steamappdetails.CapsuleImagev5,
		steamappdetails.Developers,
		steamappdetails.Publishers,
		steamappdetails.Windows,
		steamappdetails.Mac,
		steamappdetails.Linux,
		steamappdetails.Categories,
		steamappdetails.Genres,
		steamappdetails.ReleaseDate,
		steamappdetails.Background
	FROM 
		steamusergames
	JOIN 
		steamappdetails 
	ON 
		steamusergames.app_id = steamappdetails.SteamAppid;
	`
	rows, err := db.Query(query)
	if err != nil {
		return []Library{}, err
	}
	defer rows.Close()

	var library []Library

	for rows.Next() {
		var entry struct {
			AppID               int
			Type                string
			Name                string
			RequiredAge         int
			IsFree              int
			DetailedDescription string
			AboutTheGame        string
			ShortDescription    string
			SupportedLanguages  string
			HeaderImage         string
			CapsuleImage        string
			CapsuleImagev5      string
			Developers          string
			Publishers          string
			Windows             int
			Mac                 int
			Linux               int
			Categories          string
			Genres              string
			ReleaseDate         string
			Background          string
		}
		err = rows.Scan(
			&entry.AppID,
			&entry.Type,
			&entry.Name,
			&entry.RequiredAge,
			&entry.IsFree,
			&entry.DetailedDescription,
			&entry.AboutTheGame,
			&entry.ShortDescription,
			&entry.SupportedLanguages,
			&entry.HeaderImage,
			&entry.CapsuleImage,
			&entry.CapsuleImagev5,
			&entry.Developers,
			&entry.Publishers,
			&entry.Windows,
			&entry.Mac,
			&entry.Linux,
			&entry.Categories,
			&entry.Genres,
			&entry.ReleaseDate,
			&entry.Background,
		)
		if err != nil {
			return []Library{}, err
		}
		library = append(library, entry)
	}

	return library, nil
}

func GetSteamUserLibraryAppid(id int) (Entry, error) {
	fmt.Println("Database: GetSteamUserLibraryAppid")
	db, err := createConnection()
	if err != nil {
		return Entry{}, err
	}
	defer db.Close()

	/*

	 */

	var query string = `
	SELECT 
		steamusergames.app_id,
		steamusergames.playtime_forever,
		steamusergames.playtime_windows_forever,
		steamusergames.playtime_mac_forever,
		steamusergames.playtime_linux_forever,
		steamusergames.playtime_deck_forever,
		steamusergames.rtime_last_played,
		steamusergames.playtime_2weeks,
		steamappdetails.Type,
		steamappdetails.Name,
		steamappdetails.SteamAppid,
		steamappdetails.RequiredAge,
		steamappdetails.IsFree,
		steamappdetails.DetailedDescription,
		steamappdetails.AboutTheGame,
		steamappdetails.ShortDescription,
		steamappdetails.SupportedLanguages,
		steamappdetails.HeaderImage,
		steamappdetails.CapsuleImage,
		steamappdetails.CapsuleImagev5,
		steamappdetails.Developers,
		steamappdetails.Publishers,
		steamappdetails.Windows,
		steamappdetails.Mac,
		steamappdetails.Linux,
		steamappdetails.Categories,
		steamappdetails.Genres,
		steamappdetails.ReleaseDate,
		steamappdetails.Background
	FROM 
		steamusergames
	JOIN 
		steamappdetails ON steamusergames.app_id = steamappdetails.SteamAppid
	WHERE
		steamusergames.app_id = ?;
	`

	var entry Entry
	err = db.QueryRow(query, id).Scan(
		&entry.Appid,
		&entry.PlaytimeForever,
		&entry.PlaytimeWindowsForever,
		&entry.PlaytimeMacForever,
		&entry.PlaytimeLinuxForever,
		&entry.PlaytimeDeckForever,
		&entry.RtimeLastPlayed,
		&entry.Playtime2Weeks,
		&entry.Type,
		&entry.Name,
		&entry.SteamAppid,
		&entry.RequiredAge,
		&entry.IsFree,
		&entry.DetailedDescription,
		&entry.AboutTheGame,
		&entry.ShortDescription,
		&entry.SupportedLanguages,
		&entry.HeaderImage,
		&entry.CapsuleImage,
		&entry.CapsuleImagev5,
		&entry.Developers,
		&entry.Publishers,
		&entry.Windows,
		&entry.Mac,
		&entry.Linux,
		&entry.Categories,
		&entry.Genres,
		&entry.ReleaseDate,
		&entry.Background,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return Entry{}, nil
		}
		fmt.Println(err)
		return Entry{}, err
	}

	return entry, nil
}

func GetSteamUserLibraryAchievements(id int) ([]AppidAchivements, error) {
	fmt.Println("Database: GetSteamUserLibraryAchievements")
	db, err := createConnection()
	if err != nil {
		return []AppidAchivements{}, err
	}
	defer db.Close()

	var query string = `
	SELECT 
		steamuserachievements.Appid,
		steamuserachievements.Apiname,
		steamuserachievements.Achieved,
		steamuserachievements.Unlocktime,
		steamachievements.Name,
		steamachievements.DisplayName,
		steamachievements.Hidden,
		steamachievements.Description,
		steamachievements.Icon,
		steamachievements.IconGray
	FROM 
		steamuserachievements
	JOIN 
		steamachievements 
	ON 
		steamuserachievements.Appid = steamachievements.Appid
	AND 
		steamuserachievements.Apiname = steamachievements.Name;
	`
	rows, err := db.Query(query)
	if err != nil {
		return []AppidAchivements{}, err
	}
	defer rows.Close()

	var achivements []AppidAchivements
	for rows.Next() {
		var entry AppidAchivements
		err = rows.Scan(
			&entry.UserAchievements.Apiname,
			&entry.UserAchievements.Achieved,
			&entry.UserAchievements.Unlocktime,
			&entry.Achievements.Name,
			&entry.Achievements.Defaultvalue,
			&entry.Achievements.DisplayName,
			&entry.Achievements.Hidden,
			&entry.Achievements.Description,
			&entry.Achievements.Icon,
			&entry.Achievements.Icongray,
		)
		if err != nil {
			return []AppidAchivements{}, err
		}
		achivements = append(achivements, entry)
	}

	return achivements, nil
}
