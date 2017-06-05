package model

import "golang.org/x/oauth2"

type User struct {
	ID          int
	CharacterID int
	Name        string
	Email       string
}

type Authorization struct {
	UserID int

	*oauth2.Token
}
