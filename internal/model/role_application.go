package model

import (
	"time"

	"github.com/zyj/my-blog/pkg/constant"
	"gorm.io/gorm"
)

type RoleApplication struct {
	gorm.Model
	UserID          uint                           `gorm:"not null;index"`
	User            User                           `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RequestedRoleID uint                           `gorm:"not null;index"`
	RequestedRole   Role                           `gorm:"foreignKey:RequestedRoleID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Status          constant.RoleApplicationStatus `gorm:"size:16;not null;index"`
	ReviewerID      *uint
	ReviewedAt      *time.Time
	RejectReason    string `gorm:"size:512;not null;default:''"`
}
