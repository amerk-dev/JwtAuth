package repo

import (
	"JwtAuth/internal/repo/models"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
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

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (db *DB) StoreRefreshToken(guid, hashedToken, ip, email string) error {
	var user models.User

	res := db.pg.First(&user, "guid = ?", guid)
	if res.Error != nil {
		user = models.User{
			Guid:         guid,
			RefreshToken: hashedToken,
			Ip:           ip,
			Email:        email,
		}

		res = db.pg.Create(&user)
		if res.Error != nil {
			return errors.New("Failed to store refresh token")
		}
	} else {
		user.RefreshToken = hashedToken
		if ip != "" {
			user.Ip = ip
		}

		res = db.pg.Save(&user)
		if res.Error != nil {
			return errors.New("Failed to update refresh token")
		}
	}
	return nil
}
