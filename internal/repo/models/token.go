package models

import "time"

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshToken struct {
	ID        int    `gorm:"primaryKey"`
	UserGuid  string `gorm:"type:uuid;not null"`
	TokenHash string `gorm:"not null"`
	AccessJTI string `gorm:"not null"`
	IPAddress string `gorm:"not null"`
	CreatedAt time.Time
	Revoked   bool `gorm:"default:false"`
}
