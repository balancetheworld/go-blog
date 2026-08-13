package model

import (
	"github.com/zyj/my-blog/pkg/constant"
	"gorm.io/gorm"
)

type DiaryFolder struct {
	gorm.Model
	Name         string                  `gorm:"size:100;not null"`
	Slug         string                  `gorm:"size:100;not null;uniqueIndex"`
	Description  string                  `gorm:"size:500;not null;default:''"`
	Cover        string                  `gorm:"size:500;not null;default:''"`
	Sort         int                     `gorm:"not null;default:0;index"`
	Visibility   constant.PostVisibility `gorm:"size:16;not null;default:public;index"`
	VisibleRoles []Role                  `gorm:"many2many:diary_folder_visible_roles;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
