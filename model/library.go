package model

type LibraryCard struct {
	AppID        int
	Name         string
	CapsuleImage string
}

type Library struct {
	Cards []LibraryCard
}
