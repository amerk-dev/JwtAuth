package models

type User struct {
	Id           int64  `gorm:"primary_key"`
	Guid         string `json:"guid"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Ip           string `json:"ip"`
	Email        string `json:"email"`
}
