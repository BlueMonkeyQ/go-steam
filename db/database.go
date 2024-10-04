package db

import (
	"database/sql"
	"fmt"
	"go-steam/model"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// Creates a connector to the Steam Sqlite3 Database
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

func GetLibraryDB() (model.Library, error) {
	fmt.Println("Database: GetLibraryDB")
	db, err := createConnection()
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
