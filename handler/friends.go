package handler

import (
	"fmt"
	"go-steam/services"
	"go-steam/util"
	"go-steam/views"

	"github.com/labstack/echo"
)

func GetFriends(c echo.Context) error {
	fmt.Println("Endpoint: GetFriends")
	data := services.GetFriends()
	fmt.Printf("Returning #%d Friends \n", len(data))
	return util.Render(c, views.FriendsPage(data))
}

func UpdateFriends(c echo.Context) error {
	fmt.Println("Endpoint: UpdateFriends")
	if err := services.UpdateFriends(); err != nil {
		fmt.Printf("Fail: %s", err.Error())
	}
	data := services.GetFriends()
	fmt.Printf("Returning #%d Friends \n", len(data))
	return util.Render(c, views.FriendyCards(data))
}
