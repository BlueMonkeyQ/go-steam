package handler

import (
	"fmt"
	"go-steam/model"
	"go-steam/services"
	"go-steam/util"
	"go-steam/views"

	"github.com/labstack/echo"
)

func GetFriends(c echo.Context) error {
	data, err := services.GetFriends()
	if err != nil {
		c.Logger().Error(err)
		return util.Render(c, views.FriendsPage([]model.Player{}))
	}

	fmt.Printf("Returning #%d Friends \n", len(data))
	return util.Render(c, views.FriendsPage(data))
}

func UpdateFriends(c echo.Context) error {
	if err := services.UpdateFriends(); err != nil {
		fmt.Printf("Fail: %s", err.Error())
	}
	data, err := services.GetFriends()
	if err != nil {
		c.Logger().Error(err)
		return util.Render(c, views.FriendyCards([]model.Player{}))
	}

	fmt.Printf("Returning #%d Friends \n", len(data))
	return util.Render(c, views.FriendyCards(data))
}
