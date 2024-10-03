package src

type Games struct {
	Appid                  int `json:"appid"`
	PlaytimeForever        int `json:"playtime_forever"`
	PlaytimeWindowsForever int `json:"playtime_windows_forever"`
	PlaytimeMacForever     int `json:"playtime_mac_forever"`
	PlaytimeLinuxForever   int `json:"playtime_linux_forever"`
	PlaytimeDeckForever    int `json:"playtime_deck_forever"`
	RtimeLastPlayed        int `json:"rtime_last_played"`
	PlaytimeDisconnected   int `json:"playtime_disconnected"`
	Playtime2Weeks         int `json:"playtime_2weeks,omitempty"`
}
type SteamUserLibrary struct {
	Response struct {
		GameCount int     `json:"game_count"`
		Games     []Games `json:"games"`
	} `json:"response"`
}

type AppDetails struct {
	Data struct {
		Type                string   `json:"type"`
		Name                string   `json:"name"`
		SteamAppid          int      `json:"steam_appid"`
		RequiredAge         int      `json:"required_age"`
		IsFree              bool     `json:"is_free,"`
		DetailedDescription string   `json:"detailed_description"`
		AboutTheGame        string   `json:"about_the_game"`
		ShortDescription    string   `json:"short_description"`
		SupportedLanguages  string   `json:"supported_languages"`
		HeaderImage         string   `json:"header_image"`
		CapsuleImage        string   `json:"capsule_image"`
		CapsuleImagev5      string   `json:"capsule_imagev5"`
		Developers          []string `json:"developers"`
		Publishers          []string `json:"publishers"`
		Platforms           struct {
			Windows bool `json:"windows"`
			Mac     bool `json:"mac"`
			Linux   bool `json:"linux"`
		} `json:"platforms"`
		Categories []struct {
			Description string `json:"description"`
		} `json:"categories"`
		Genres []struct {
			Description string `json:"description"`
		} `json:"genres"`
		ReleaseDate struct {
			Date string `json:"date"`
		} `json:"release_date"`
		Background string `json:"background"`
	} `json:"data,omitempty"`
}

type Library struct {
	AppID               int
	Type                string
	Name                string
	RequiredAge         int
	IsFree              int
	DetailedDescription string
	AboutTheGame        string
	ShortDescription    string
	SupportedLanguages  string
	HeaderImage         string
	CapsuleImage        string
	CapsuleImagev5      string
	Developers          string
	Publishers          string
	Windows             int
	Mac                 int
	Linux               int
	Categories          string
	Genres              string
	ReleaseDate         string
	Background          string
}

type Entry struct {
	Appid                  int
	PlaytimeForever        int
	PlaytimeWindowsForever int
	PlaytimeMacForever     int
	PlaytimeLinuxForever   int
	PlaytimeDeckForever    int
	RtimeLastPlayed        int
	Playtime2Weeks         int
	Type                   string
	Name                   string
	SteamAppid             int
	RequiredAge            int
	IsFree                 int
	DetailedDescription    string
	AboutTheGame           string
	ShortDescription       string
	SupportedLanguages     string
	HeaderImage            string
	CapsuleImage           string
	CapsuleImagev5         string
	Developers             string
	Publishers             string
	Windows                int
	Mac                    int
	Linux                  int
	Categories             string
	Genres                 string
	ReleaseDate            string
	Background             string
}

type Achievements struct {
	Name         string `json:"name"`
	Defaultvalue int    `json:"defaultvalue"`
	DisplayName  string `json:"displayName"`
	Hidden       int    `json:"hidden"`
	Description  string `json:"description"`
	Icon         string `json:"icon"`
	Icongray     string `json:"icongray"`
}

type SteamAchievements struct {
	Game struct {
		GameName           string `json:"gameName"`
		GameVersion        string `json:"gameVersion"`
		AvailableGameStats struct {
			Achievements []Achievements `json:"achievements"`
			Stats        []struct {
				Name         string `json:"name"`
				Defaultvalue int    `json:"defaultvalue"`
				DisplayName  string `json:"displayName"`
			} `json:"stats"`
		} `json:"availableGameStats"`
	} `json:"game"`
}

type UserAchievements struct {
	Apiname    string `json:"apiname,omitempty"`
	Achieved   int    `json:"achieved,omitempty"`
	Unlocktime int    `json:"unlocktime,omitempty"`
}

type SteamUserAchivements struct {
	Playerstats struct {
		SteamID      string             `json:"steamID,omitempty"`
		GameName     string             `json:"gameName,omitempty"`
		Achievements []UserAchievements `json:"achievements,omitempty"`
		Success      bool               `json:"success,omitempty"`
	} `json:"playerstats"`
}

type AppidAchivements struct {
	UserAchievements UserAchievements
	Achievements     Achievements
}
