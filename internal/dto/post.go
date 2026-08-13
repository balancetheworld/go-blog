package dto

import (
	"time"

	"github.com/zyj/my-blog/pkg/constant"
)

type CreatePostRequest struct {
	Title          string                  `json:"title" vd:"len($) >= 1 && len($) <= 255"`
	DraftContent   string                  `json:"draft_content"`
	Description    string                  `json:"description" vd:"len($) <= 1000"`
	Cover          string                  `json:"cover" vd:"len($) <= 512"`
	Type           string                  `json:"type" vd:"len($) <= 32"`
	CategoryID     *uint                   `json:"category_id"`
	LabelIDs       []uint                  `json:"label_ids"`
	IsPrivate      bool                    `json:"is_private"`
	Top            bool                    `json:"top"`
	Publish        bool                    `json:"publish"`
	Visibility     constant.PostVisibility `json:"visibility" vd:"$ == '' || in($, 'public', 'roles', 'private')"`
	VisibleRoleIDs []uint                  `json:"visible_role_ids"`
}

type UpdatePostRequest struct {
	Title          *string                  `json:"title" vd:"$ == nil || (len($) >= 1 && len($) <= 255)"`
	DraftContent   *string                  `json:"draft_content"`
	Description    *string                  `json:"description" vd:"$ == nil || len($) <= 1000"`
	Cover          *string                  `json:"cover" vd:"$ == nil || len($) <= 512"`
	Type           *string                  `json:"type" vd:"$ == nil || len($) <= 32"`
	CategoryID     *uint                    `json:"category_id"`
	LabelIDs       *[]uint                  `json:"label_ids"`
	IsPrivate      *bool                    `json:"is_private"`
	Top            *bool                    `json:"top"`
	Publish        *bool                    `json:"publish"`
	Visibility     *constant.PostVisibility `json:"visibility" vd:"$ == nil || in($, 'public', 'roles', 'private')"`
	VisibleRoleIDs *[]uint                  `json:"visible_role_ids"`
}

type PostListRequest struct {
	Page       int    `query:"page" vd:"$ == 0 || $ >= 1"`
	PageSize   int    `query:"page_size" vd:"$ == 0 || ($ >= 1 && $ <= 100)"`
	Keyword    string `query:"keyword"`
	Type       string `query:"type"`
	CategoryID uint   `query:"category_id"`
	LabelID    uint   `query:"label_id"`
	AuthorID   uint   `query:"author_id"`
	Status     string `query:"status" vd:"$ == '' || in($, 'published', 'draft', 'all')"`
	Sort       string `query:"sort" vd:"$ == '' || in($, 'latest', 'oldest', 'hot')"`
}

type CategoryResponse struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type LabelResponse struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type CreateLabelRequest struct {
	Name string `json:"name" vd:"len($) >= 1 && len($) <= 64"`
	Slug string `json:"slug" vd:"len($) >= 1 && len($) <= 64"`
}

type UpdateLabelRequest struct {
	Name *string `json:"name" vd:"$ == nil || (len($) >= 1 && len($) <= 64)"`
	Slug *string `json:"slug" vd:"$ == nil || (len($) >= 1 && len($) <= 64)"`
}

type CreateCategoryRequest struct {
	Name        string `json:"name" vd:"len($) >= 1 && len($) <= 64"`
	Slug        string `json:"slug" vd:"len($) >= 1 && len($) <= 64"`
	Description string `json:"description" vd:"len($) <= 512"`
}

type UpdateCategoryRequest struct {
	Name        *string `json:"name" vd:"$ == nil || (len($) >= 1 && len($) <= 64)"`
	Slug        *string `json:"slug" vd:"$ == nil || (len($) >= 1 && len($) <= 64)"`
	Description *string `json:"description" vd:"$ == nil || len($) <= 512"`
}

type PostDetailResponse struct {
	ID           uint64                  `json:"id"`
	Title        string                  `json:"title"`
	Content      string                  `json:"content"`
	DraftContent *string                 `json:"draft_content,omitempty"`
	Description  string                  `json:"description"`
	Cover        string                  `json:"cover"`
	Type         string                  `json:"type"`
	Slug         string                  `json:"slug"`
	CategoryID   *uint                   `json:"category_id"`
	Category     *CategoryResponse       `json:"category,omitempty"`
	Labels       []LabelResponse         `json:"labels"`
	Author       UserDto                 `json:"author"`
	IsPrivate    bool                    `json:"is_private"`
	Visibility   constant.PostVisibility `json:"visibility"`
	VisibleRoles []RoleOptionResponse    `json:"visible_roles"`
	Top          bool                    `json:"top"`
	LikeCount    uint64                  `json:"like_count"`
	Liked        bool                    `json:"liked"`
	CommentCount uint64                  `json:"comment_count"`
	ViewCount    uint64                  `json:"view_count"`
	Heat         float64                 `json:"heat"`
	Status       string                  `json:"status"`
	PublishedAt  *time.Time              `json:"published_at"`
	CreatedAt    time.Time               `json:"created_at"`
	UpdatedAt    time.Time               `json:"updated_at"`
}

type PostLikeResponse struct {
	Liked     bool   `json:"liked"`
	LikeCount uint64 `json:"like_count"`
}

type PostListItemResponse struct {
	ID           uint64                  `json:"id"`
	Title        string                  `json:"title"`
	Content      string                  `json:"content"`
	Description  string                  `json:"description"`
	Cover        string                  `json:"cover"`
	Type         string                  `json:"type"`
	Slug         string                  `json:"slug"`
	Category     *CategoryResponse       `json:"category,omitempty"`
	Labels       []LabelResponse         `json:"labels"`
	Author       UserDto                 `json:"author"`
	IsPrivate    bool                    `json:"is_private"`
	Visibility   constant.PostVisibility `json:"visibility"`
	VisibleRoles []RoleOptionResponse    `json:"visible_roles"`
	Top          bool                    `json:"top"`
	LikeCount    uint64                  `json:"like_count"`
	CommentCount uint64                  `json:"comment_count"`
	ViewCount    uint64                  `json:"view_count"`
	Heat         float64                 `json:"heat"`
	Status       string                  `json:"status"`
	PublishedAt  *time.Time              `json:"published_at"`
	CreatedAt    time.Time               `json:"created_at"`
}

type PostListResponse struct {
	Items    []PostListItemResponse `json:"items"`
	Total    int64                  `json:"total"`
	Page     int                    `json:"page"`
	PageSize int                    `json:"page_size"`
}
