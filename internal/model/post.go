package model

import (
	"time"

	"github.com/zyj/my-blog/internal/dto"
	"gorm.io/gorm"
)

type PostBase struct {
	Title        string  `gorm:"size:255;not null"`
	Content      string  `gorm:"type:text;not null;default:''"`
	DraftContent string  `gorm:"type:text;not null;default:''"`
	Description  string  `gorm:"size:1000;not null;default:''"`
	Cover        string  `gorm:"size:512;not null;default:''"`
	Type         string  `gorm:"size:32;not null;default:article"`
	Slug         string  `gorm:"size:255;not null;uniqueIndex"`
	CategoryID   *uint   `gorm:"index"`
	IsPrivate    bool    `gorm:"not null;default:false;index"`
	Top          bool    `gorm:"not null;default:false;index"`
	LikeCount    uint64  `gorm:"not null;default:0"`
	CommentCount uint64  `gorm:"not null;default:0"`
	ViewCount    uint64  `gorm:"not null;default:0"`
	Heat         float64 `gorm:"not null;default:0"`
	PublishedAt  *time.Time
}

type Post struct {
	gorm.Model
	PostBase
	AuthorID uint `gorm:"not null;index"`
	Author   User `gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Category *Category `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Labels   []Label `gorm:"many2many:post_labels;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (p Post) ToDto() dto.PostDetailResponse {
	return dto.PostDetailResponse{
		ID:           uint64(p.ID),
		Title:        p.Title,
		Content:      p.Content,
		Description:  p.Description,
		Cover:        p.Cover,
		Type:         p.Type,
		Slug:         p.Slug,
		CategoryID:   p.CategoryID,
		Category:     p.categoryToDto(),
		Labels:       p.labelsToDto(),
		Author:       p.Author.ToDto(),
		IsPrivate:    p.IsPrivate,
		Top:          p.Top,
		LikeCount:    p.LikeCount,
		CommentCount: p.CommentCount,
		ViewCount:    p.ViewCount,
		Heat:         p.Heat,
		Status:       p.status(),
		PublishedAt:  p.PublishedAt,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}
}

func (p Post) ToDtoWithShortContent() dto.PostListItemResponse {
	content := []rune(p.Content)
	if len(content) > 200 {
		content = content[:200]
	}

	return dto.PostListItemResponse{
		ID:           uint64(p.ID),
		Title:        p.Title,
		Content:      string(content),
		Description:  p.Description,
		Cover:        p.Cover,
		Type:         p.Type,
		Slug:         p.Slug,
		Category:     p.categoryToDto(),
		Labels:       p.labelsToDto(),
		Author:       p.Author.ToDto(),
		IsPrivate:    p.IsPrivate,
		Top:          p.Top,
		LikeCount:    p.LikeCount,
		CommentCount: p.CommentCount,
		ViewCount:    p.ViewCount,
		Heat:         p.Heat,
		Status:       p.status(),
		PublishedAt:  p.PublishedAt,
		CreatedAt:    p.CreatedAt,
	}
}

func (p Post) CalculateHeat() float64 {
	return float64(p.ViewCount) + float64(p.LikeCount)*3 + float64(p.CommentCount)*5
}

func (p Post) status() string {
	if p.Content == "" {
		return "draft"
	}

	return "published"
}

func (p Post) categoryToDto() *dto.CategoryResponse {
	if p.Category == nil {
		return nil
	}

	return &dto.CategoryResponse{
		ID:          uint64(p.Category.ID),
		Name:        p.Category.Name,
		Slug:        p.Category.Slug,
		Description: p.Category.Description,
		CreatedAt:   p.Category.CreatedAt,
		UpdatedAt:   p.Category.UpdatedAt,
	}
}

func (p Post) labelsToDto() []dto.LabelResponse {
	labels := make([]dto.LabelResponse, 0, len(p.Labels))
	for _, label := range p.Labels {
		labels = append(labels, dto.LabelResponse{
			ID:   uint64(label.ID),
			Name: label.Name,
			Slug: label.Slug,
		})
	}

	return labels
}
