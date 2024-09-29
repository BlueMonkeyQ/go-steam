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
	Num240 struct {
		Success bool `json:"success"`
		Data    struct {
			Type                string   `json:"type"`
			Name                string   `json:"name"`
			SteamAppid          int      `json:"steam_appid"`
			RequiredAge         int      `json:"required_age"`
			IsFree              bool     `json:"is_free"`
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
				ID          int    `json:"id"`
				Description string `json:"description"`
			} `json:"categories"`
			Genres []struct {
				ID          string `json:"id"`
				Description string `json:"description"`
			} `json:"genres"`
			ReleaseDate struct {
				Date string `json:"date"`
			} `json:"release_date"`
			Background string `json:"background"`
		} `json:"data"`
	} `json:"240"`
}
