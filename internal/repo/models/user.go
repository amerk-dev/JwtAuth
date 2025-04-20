package models

type User struct {
	Id    int64  `gorm:"primary_key"`
	Guid  string `json:"guid" gorm:"uniqueIndex"` // UUID
	Ip    string `json:"ip"`
	Email string `json:"email"`
}
