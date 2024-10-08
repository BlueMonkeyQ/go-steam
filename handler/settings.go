package handler

import (
	"go-steam/model"
	"go-steam/services"
	"go-steam/util"
	"go-steam/views"
	"strings"

	"github.com/labstack/echo"
)

func SettingsPage(c echo.Context) error {
	var settings model.Settings
	settings.Steamkey = services.GetSettingsSteamkey()
	settings.Steamid = services.GetSettingsSteamid()

	steamkey := c.QueryParam("steamkey")
	if strings.Compare(steamkey, "") != 0 {
		if err := services.UpdateSettingsSteamkey(steamkey); err != nil {
			return util.Render(c, views.SettingsPage(model.Settings{}))
		}
		settings.Steamkey = steamkey
		if err := services.ValidateSettings(); err != nil {
			settings.Valid = "False"
		} else {
			settings.Valid = "True"
		}
		return util.Render(c, views.SettingsInfo(settings))

	}

	steamid := c.QueryParam("steamid")
	if strings.Compare(steamid, "") != 0 {
		if err := services.UpdateSettingsSteamid(steamid); err != nil {
			return util.Render(c, views.SettingsPage(model.Settings{}))
		}
		settings.Steamid = steamid
		if err := services.ValidateSettings(); err != nil {
			settings.Valid = "False"
		} else {
			settings.Valid = "True"
		}
		return util.Render(c, views.SettingsInfo(settings))

	}

	return util.Render(c, views.SettingsPage(settings))
}
