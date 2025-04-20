package domain

type User struct {
	Id    int64  `json:"id"`
	Guid  string `json:"guid"`
	Ip    string `json:"ip"`
	Email string `json:"email"`
}
