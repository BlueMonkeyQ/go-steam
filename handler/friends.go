package handler

import (
	"fmt"
	"go-steam/views"

	"github.com/labstack/echo"
)

func GetFriends(c echo.Context) error {
	fmt.Println("Endpoint: GetFriends")
	return render(c, views.FriendsPage())
}
