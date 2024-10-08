package util

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/labstack/echo"
	"gopkg.in/ini.v1"
)

func Render(c echo.Context, component templ.Component) error {
	return component.Render(c.Request().Context(), c.Response())
}

func GetParam(c echo.Context, param string) string {
	value := c.Param(param)
	return value

}

func StringToInt(value string) (int, error) {
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return intValue, nil
}

func IntToString(value int) string {
	stringValue := strconv.Itoa(value)
	return stringValue
}

func StringToTime(value string) string {
	unlockTimeInt, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		unlockTimeInt = 0
	}
	return time.Unix(unlockTimeInt, 0).Format(time.RFC1123)
}

func InitConfig() error {
	_, err := os.Stat(".config.ini")
	if err != nil {
		fmt.Println("Initializing config.ini...")
		config := ini.Empty()

		// Create sections
		sec_api, err := config.NewSection("API")
		if err != nil {
			msg := fmt.Sprintf("Warning: %v\n", err)
			fmt.Println(msg)
			panic(err)
		}

		// Create Keys
		if _, err := sec_api.NewKey("Steamkey", ""); err != nil {
			msg := fmt.Sprintf("Warning: %v\n", err)
			fmt.Println(msg)
			panic(err)
		}
		if _, err := sec_api.NewKey("Steamid", ""); err != nil {
			msg := fmt.Sprintf("Warning: %v\n", err)
			fmt.Println(msg)
			panic(err)
		}

		// Create file
		if err := config.SaveTo(".config.ini"); err != nil {
			msg := fmt.Sprintf("Warning: %v\n", err)
			fmt.Println(msg)
			panic(err)
		}
	}
	return nil
}

func getConfig() *ini.File {
	// Returns config file

	config, err := ini.Load(".config.ini")
	if err != nil {
		msg := fmt.Sprintf("Warning: %v\n", err)
		fmt.Println(msg)
		panic(err)
	}

	return config
}

func GetSteamKey() string {
	config := getConfig()
	return config.Section("API").Key("Steamkey").String()
}

func GetSteamId() string {
	config := getConfig()
	return config.Section("API").Key("Steamid").String()
}

func UpdateSteamKey(steamkey string) error {
	config := getConfig()
	config.Section("API").Key("Steamkey").SetValue(steamkey)
	if err := config.SaveTo(".config.ini"); err != nil {
		msg := fmt.Sprintf("Warning: %v\n", err)
		fmt.Println(msg)
	}
	return nil
}

func UpdateSteamId(steamid string) error {
	config := getConfig()
	config.Section("API").Key("Steamid").SetValue(steamid)
	if err := config.SaveTo(".config.ini"); err != nil {
		msg := fmt.Sprintf("Warning: %v\n", err)
		fmt.Println(msg)
	}
	return nil
}

func ValidateSettings() error {
	steamkey := GetSteamKey()
	steamid := GetSteamId()

	if strings.Compare(steamkey, "") == 0 || strings.Compare(steamid, "") == 0 {
		return errors.New("invalid")
	}
	return nil
}
