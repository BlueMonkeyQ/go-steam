package services

import (
	"fmt"
	"go-steam/db"
	"go-steam/model"
)

func GetLibrary() model.Library {
	data, err := db.GetLibraryDB()
	if err != nil {
		fmt.Println(err)
		return model.Library{}
	}
	return data
}
