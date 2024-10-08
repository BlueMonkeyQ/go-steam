package services

import (
	"fmt"
	"go-steam/util"
	"net/http"
)

func ValidateSettings() error {
	if err := util.ValidateSettings(); err != nil {
		return err
	}

	steamkey := util.GetSteamKey()
	steamid := util.GetSteamId()
	url := fmt.Sprintf("http://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key=%s&steamid=%s&format=json", steamkey, steamid)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

func UpdateSettingsSteamkey(steamkey string) error {
	fmt.Println("Endpoint: UpdateSettingsSteamkey")
	if err := util.UpdateSteamKey(steamkey); err != nil {
		return err
	}
	return nil
}

func UpdateSettingsSteamid(steamid string) error {
	fmt.Println("Endpoint: UpdateSettingsSteamid")
	if err := util.UpdateSteamId(steamid); err != nil {
		return err
	}
	return nil
}

func GetSettingsSteamkey() string {
	fmt.Println("Endpoint: GetSettingsSteamkey")
	return util.GetSteamKey()
}

func GetSettingsSteamid() string {
	fmt.Println("Endpoint: GetSettingsSteamid")
	return util.GetSteamId()
}
