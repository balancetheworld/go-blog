package model

import (
	"time"

	"github.com/zyj/my-blog/pkg/constant"
	"gorm.io/gorm"
)

type Diary struct {
	gorm.Model
	Title        string                  `gorm:"size:200;not null;default:''"`
	Slug         string                  `gorm:"size:200;not null;default:'';uniqueIndex:idx_diaries_slug,where:slug <> ''"`
	Description  string                  `gorm:"size:500;not null;default:''"`
	Cover        string                  `gorm:"size:500;not null;default:''"`
	Content      string                  `gorm:"type:text;not null;default:''"`
	DraftContent string                  `gorm:"type:text;not null;default:''"`
	AuthorID     uint                    `gorm:"not null;index"`
	Author       User                    `gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	FolderID     *uint                   `gorm:"index"`
	Folder       *DiaryFolder            `gorm:"foreignKey:FolderID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Visibility   constant.PostVisibility `gorm:"size:16;not null;default:public;index"`
	VisibleRoles []Role                  `gorm:"many2many:diary_visible_roles;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ViewCount    uint64                  `gorm:"not null;default:0"`
	LikeCount    uint64                  `gorm:"not null;default:0"`
	CommentCount uint64                  `gorm:"not null;default:0"`
	PublishedAt  *time.Time              `gorm:"index"`
}
