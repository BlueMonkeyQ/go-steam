package services

import (
	"encoding/json"
	"fmt"
	"go-steam/db"
	"go-steam/model"
	"io"
	"net/http"
	"strings"
)

func GetFriends() []model.Player {
	fmt.Println("Endpoint: GetFriends")
	data, err := db.GetFriendsDB(76561198050437739)
	if err != nil {
		fmt.Println(err)
		return []model.Player{}
	}
	return data
}

func UpdateUser(steamids []string) error {
	fmt.Println("Endpoint: UpdateUser")

	param := strings.Join(steamids, ",")
	url := fmt.Sprintf("https://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002/?key=14EB214CEC3F1701FD192885D330990F&steamids=%s", param)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	var gps model.GetPlayerSummariesAPI
	if err := json.Unmarshal(body, &gps); err != nil {
		fmt.Println(err)
		return err
	}

	if err := db.InsertSteamUsersDB(gps.Response.Players); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func UpdateFriends() error {
	fmt.Println("Endpoint: UpdateFriends")
	resp, err := http.Get("https://api.steampowered.com/ISteamUser/GetFriendList/v0001/?key=14EB214CEC3F1701FD192885D330990F&steamid=76561198050437739&relationship=friend")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	var gfla model.GetFriendListAPI
	if err := json.Unmarshal(body, &gfla); err != nil {
		fmt.Println(err)
		return err
	}

	var steamids []string
	for _, i := range gfla.Friendslist.Friends {
		steamids = append(steamids, i.Steamid)
	}

	if err := UpdateUser(steamids); err != nil {
		fmt.Println(err)
		return err
	}

	if err := db.InsertSteamFriendsDB(gfla.Friendslist.Friends, 76561198050437739); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
