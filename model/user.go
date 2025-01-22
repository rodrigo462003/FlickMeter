package model

import "gorm.io/gorm"

type User struct {
    gorm.Model
    Username   string `gorm:"type:varchar(50);unique;not null"`
    Email      string `gorm:"type:varchar(254);unique;not null"`
    Password   string `gorm:"type:varchar(255);not null"`
}
