package domain

type User struct {
	ID           int64  `json:"id"`
	Guid         string `json:"guid"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Ip           string `json:"ip"`
	Email        string `json:"email"`
}
