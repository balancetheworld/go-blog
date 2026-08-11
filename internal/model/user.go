package model

import (
	"github.com/zyj/my-blog/internal/dto"
	"github.com/zyj/my-blog/pkg/constant"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username       string `gorm:"size:32;not null;uniqueIndex"`
	Nickname       string `gorm:"size:64;not null;default:''"`
	AvatarURL      string `gorm:"size:512;not null;default:''"`
	BackgroundURL  string `gorm:"size:512;not null;default:''"`
	PreferredColor string `gorm:"size:32;not null;default:''"`
	Email          string `gorm:"size:254;not null;uniqueIndex"`
	Gender         string `gorm:"size:16;not null;default:''"`
	Role           constant.Role `gorm:"size:16;not null;default:user;check:role IN ('user','editor','admin')"`
	RoleID         *uint         `gorm:"index"`
	CurrentRole    Role          `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	IsRoot         bool          `gorm:"not null;default:false;index"`
	Language       string `gorm:"size:16;not null;default:''"`
	PasswordHash   string `gorm:"size:255;not null" json:"-"`
	ShowIPLocation bool   `gorm:"not null;default:false"`
}

// 给 user 结构体绑定 todto 方法
func (u User) ToDto() dto.UserDto {
	return dto.UserDto{
		ID:        uint64(u.ID),
		Username:  u.Username,
		Nickname:  u.Nickname,
		Avatar:    u.AvatarURL,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (u User) IsAdmin() bool {
	return u.Role == constant.RoleAdmin
}

func (u User) IsEditor() bool {
	return u.Role == constant.RoleEditor
}

func (u User) IsEditorOrAdmin() bool {
	return u.IsEditor() || u.IsAdmin()
}
