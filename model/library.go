package model

type Settings struct {
	Steamkey string
	Steamid  string
	Valid    string
}

type LibraryCard struct {
	AppID            string
	Name             string
	HeaderImage      string
	RtimeLastPlayed  string
	TotalAchivements string
	TotalAchieved    string
}

type Library struct {
	Cards         []LibraryCard
	FilterOptions FilterOptions
}

type Achievement struct {
	Name        string
	DisplayName string
	Hidden      bool
	Description string
	Icon        string
	IconGray    string
	Achieved    bool
	Unlocktime  string
	Percentage  string
}

type AppDetails struct {
	AppID                  string
	SteamAppID             string
	Name                   string
	AboutTheGame           string
	ShortDescripton        string
	SupportedLanguages     string
	HeaderImage            string
	Developers             []string
	Publishers             []string
	Windows                bool
	Mac                    bool
	Linux                  bool
	Categories             []string
	Genres                 []string
	ReleaseDate            string
	Background             string
	PlaytimeForever        int
	PlaytimeWindowsForever int
	PlaytimeMacForever     int
	PlaytimeLinuxForever   int
	RtimeLastPlayed        string
	Playtime2Weeks         int
}

type AchivementDetails struct {
	Achievements []Achievement
	LastUpdated  string
}

type GameData struct {
	AppDetails   AppDetails
	Achievements AchivementDetails
}

type Player struct {
	Steamid                  string
	FriendSince              string
	Communityvisibilitystate int
	Profilestate             int
	Personaname              string
	Profileurl               string
	Avatar                   string
	Avatarmedium             string
	Avatarfull               string
	Lastlogoff               string
	Personastate             int
}

type FilterOptions struct {
	Categories []string
	Genres     []string
	Developers []string
	Publishers []string
}
