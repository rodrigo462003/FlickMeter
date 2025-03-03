package db

import (
	"database/sql"
	"fmt"

	"github.com/rodrigo462003/FlickMeter/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(connStr string) *gorm.DB {
	sqlDB, err := sql.Open("pgx", connStr)
	if err != nil {
		panic(fmt.Sprintln("sqlDB error: ", err))
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	gormDB.AutoMigrate(&model.User{}, &model.VerificationCode{}, &model.Auth{})

	return gormDB
}
