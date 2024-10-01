package models

type User struct {
	ID                int    `json:"-"`
	UserName          string `json:"username"`
	EncryptedPassword string `json:"-"`
	Password          string `json:"password,omitempty"`
	Token             string `json:"token,omitempty"`
}

type ParsedToken struct {
	UserName string `json:"username"`
}
