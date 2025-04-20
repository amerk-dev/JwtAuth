package repo

import (
	"JwtAuth/internal/repo/models"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"time"
)

type DB struct {
	pg *gorm.DB
}

func NewDB(pg *gorm.DB) *DB {
	return &DB{pg}
}

func InitDB() (*gorm.DB, error) {
	host := os.Getenv("POSTGRES_HOST")
	port := 5432
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	sslMode := "disable"
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslMode)
	fmt.Println(dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.User{}, &models.RefreshToken{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (db *DB) StoreRefreshToken(guid, tokenHash, accessJTI, ip string) error {
	rt := models.RefreshToken{
		UserGuid:  guid,
		TokenHash: tokenHash,
		AccessJTI: accessJTI,
		IPAddress: ip,
		CreatedAt: time.Now(),
		Revoked:   false,
	}
	return db.pg.Create(&rt).Error
}

func (db *DB) FindRefreshToken(raw string) (*models.RefreshToken, error) {
	var tokens []models.RefreshToken
	if err := db.pg.Where("revoked = false").Find(&tokens).Error; err != nil {
		return nil, err
	}
	for _, t := range tokens {
		if err := bcrypt.CompareHashAndPassword([]byte(t.TokenHash), []byte(raw)); err == nil {
			return &t, nil
		}
	}
	return nil, errors.New("refresh token not found")
}

func (db *DB) UpdateRefreshToken(oldID int, guid, newHash, newAccessJTI, ip string) error {
	tx := db.pg.Begin()

	if err := tx.Model(&models.RefreshToken{}).Where("id = ?", oldID).Update("revoked", true).Error; err != nil {
		tx.Rollback()
		return err
	}

	newToken := models.RefreshToken{
		UserGuid:  guid,
		TokenHash: newHash,
		AccessJTI: newAccessJTI,
		IPAddress: ip,
		CreatedAt: time.Now(),
		Revoked:   false,
	}

	if err := tx.Create(&newToken).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
