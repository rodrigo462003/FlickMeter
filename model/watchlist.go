package model

type Watchlist struct {
	UserID  uint `gorm:"not null;index;uniqueIndex:idx_user_movie"`
	MovieID uint `gorm:"not null;index;uniqueIndex:idx_user_movie"`
}
