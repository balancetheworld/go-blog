package model

import (
	"github.com/zyj/my-blog/internal/dto"
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	PostID   uint   `gorm:"not null;index"`
	Post     Post   `gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	AuthorID uint   `gorm:"not null;index"`
	Author   User   `gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Content  string `gorm:"type:text;not null"`
}

func (c Comment) ToDto() dto.CommentResponse {
	return dto.CommentResponse{
		ID:        uint64(c.ID),
		PostID:    uint64(c.PostID),
		Content:   c.Content,
		Author:    c.Author.ToDto(),
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
