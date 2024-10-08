package services

import (
	"encoding/json"
	"fmt"
	"go-steam/db"
	"go-steam/model"
	"go-steam/util"
	"io"
	"net/http"
	"strings"
)

func GetFriends() ([]model.Player, error) {
	fmt.Println("Endpoint: GetFriends")
	if err := ValidateSettings(); err != nil {
		return []model.Player{}, err
	}

	steamid := util.GetSteamId()
	data, err := db.GetFriendsDB(steamid)
	if err != nil {
		return []model.Player{}, err
	}
	return data, nil
}

func UpdateUser(steamids []string) error {
	fmt.Println("Endpoint: UpdateUser")
	if err := ValidateSettings(); err != nil {
		return err
	}

	param := strings.Join(steamids, ",")
	steamkey := util.GetSteamKey()
	url := fmt.Sprintf("https://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002/?key=%s&steamids=%s", steamkey, param)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var gps model.GetPlayerSummariesAPI
	if err := json.Unmarshal(body, &gps); err != nil {
		return err
	}

	if err := db.InsertSteamUsersDB(gps.Response.Players); err != nil {
		return err
	}

	return nil
}

func UpdateFriends() error {
	fmt.Println("Endpoint: UpdateFriends")
	if err := ValidateSettings(); err != nil {
		fmt.Println(err)
		return err
	}

	steamkey := util.GetSteamKey()
	steamid := util.GetSteamId()
	url := fmt.Sprintf("https://api.steampowered.com/ISteamUser/GetFriendList/v0001/?key=%s&steamid=%s&relationship=friend", steamkey, steamid)
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

	if err := db.InsertSteamFriendsDB(gfla.Friendslist.Friends, steamid); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
