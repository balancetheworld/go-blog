package model

import "time"

type PostLike struct {
	PostID    uint `gorm:"primaryKey"`
	Post      Post `gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserID    uint `gorm:"primaryKey"`
	User      User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt time.Time
}
