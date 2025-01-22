package userHandler

import (
	"github.com/rodrigo462003/FlickMeter/store"
	"gorm.io/gorm"
)

type UserHandler struct {
	us *store.UserStore
}

func NewUserHandler(d *gorm.DB) *UserHandler {
	return &UserHandler{us: store.NewUserStore(d)}
}
