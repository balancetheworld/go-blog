package model

import (
	"time"

	"github.com/zyj/my-blog/pkg/constant"
	"gorm.io/gorm"
)

type PostBase struct {
	Title        string                  `gorm:"size:255;not null"`
	Content      string                  `gorm:"type:text;not null;default:''"`
	DraftContent string                  `gorm:"type:text;not null;default:''"`
	Description  string                  `gorm:"size:1000;not null;default:''"`
	Cover        string                  `gorm:"size:512;not null;default:''"`
	Type         string                  `gorm:"size:32;not null;default:article"`
	Slug         string                  `gorm:"size:255;not null;uniqueIndex"`
	CategoryID   *uint                   `gorm:"index"`
	IsPrivate    bool                    `gorm:"not null;default:false;index"`
	Visibility   constant.PostVisibility `gorm:"size:16;not null;default:public;index"`
	Top          bool                    `gorm:"not null;default:false;index"`
	LikeCount    uint64                  `gorm:"not null;default:0"`
	CommentCount uint64                  `gorm:"not null;default:0"`
	ViewCount    uint64                  `gorm:"not null;default:0"`
	Heat         float64                 `gorm:"not null;default:0"`
	PublishedAt  *time.Time
}

type Post struct {
	gorm.Model
	PostBase
	AuthorID     uint      `gorm:"not null;index"`
	Author       User      `gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Category     *Category `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Labels       []Label   `gorm:"many2many:post_labels;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	VisibleRoles []Role    `gorm:"many2many:post_visible_roles;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (p Post) CalculateHeat() float64 {
	return float64(p.ViewCount) + float64(p.LikeCount)*3 + float64(p.CommentCount)*5
}
