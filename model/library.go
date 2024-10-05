package model

type LibraryCard struct {
	AppID            string
	Name             string
	HeaderImage      string
	RtimeLastPlayed  string
	TotalAchivements string
	TotalAchieved    string
}

type Library struct {
	Cards []LibraryCard
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
}

type AppDetails struct {
	AppID              string
	SteamAppID         string
	Name               string
	AboutTheGame       string
	ShortDescripton    string
	SupportedLanguages string
	HeaderImage        string
	Developers         []string
	Publishers         []string
	Windows            bool
	Mac                bool
	Linux              bool
	Categories         []string
	Genres             []string
	ReleaseDate        string
	Background         string
}

type AchivementDetails struct {
	Achievements []Achievement
	LastUpdated  string
}

type GameData struct {
	AppDetails   AppDetails
	Achievements AchivementDetails
}
