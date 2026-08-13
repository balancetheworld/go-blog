package dto

import (
	"time"

	"github.com/zyj/my-blog/pkg/constant"
)

type CreateDiaryRequest struct {
	Title          string                  `json:"title" vd:"len($) <= 200"`
	Slug           string                  `json:"slug" vd:"len($) <= 200"`
	Description    string                  `json:"description" vd:"len($) <= 500"`
	Cover          string                  `json:"cover" vd:"len($) <= 500"`
	FolderID       *uint                   `json:"folder_id"`
	DraftContent   string                  `json:"draft_content" vd:"len($) <= 1000000"`
	Publish        bool                    `json:"publish"`
	Visibility     constant.PostVisibility `json:"visibility" vd:"$ == '' || in($, 'public', 'roles', 'private')"`
	VisibleRoleIDs []uint                  `json:"visible_role_ids"`
}

type UpdateDiaryRequest struct {
	Title          *string                  `json:"title" vd:"$ == nil || len($) <= 200"`
	Slug           *string                  `json:"slug" vd:"$ == nil || len($) <= 200"`
	Description    *string                  `json:"description" vd:"$ == nil || len($) <= 500"`
	Cover          *string                  `json:"cover" vd:"$ == nil || len($) <= 500"`
	FolderID       *uint                    `json:"folder_id"`
	ClearFolder    bool                     `json:"clear_folder"`
	DraftContent   *string                  `json:"draft_content" vd:"$ == nil || len($) <= 1000000"`
	Publish        *bool                    `json:"publish"`
	Visibility     *constant.PostVisibility `json:"visibility" vd:"$ == nil || in($, 'public', 'roles', 'private')"`
	VisibleRoleIDs *[]uint                  `json:"visible_role_ids"`
}

type ListDiariesRequest struct {
	Page     int    `query:"page" vd:"$ == 0 || $ >= 1"`
	PageSize int    `query:"page_size" vd:"$ == 0 || ($ >= 1 && $ <= 100)"`
	Status   string `query:"status" vd:"$ == '' || in($, 'published', 'draft', 'all')"`
	Keyword  string `query:"keyword" vd:"len($) <= 100"`
	FolderID uint   `query:"folder_id"`
}

type DiaryResponse struct {
	ID           uint64                  `json:"id"`
	Title        string                  `json:"title"`
	Slug         string                  `json:"slug"`
	Description  string                  `json:"description"`
	Cover        string                  `json:"cover"`
	Content      string                  `json:"content"`
	DraftContent *string                 `json:"draft_content,omitempty"`
	Author       UserDto                 `json:"author"`
	Folder       *DiaryFolderResponse    `json:"folder"`
	Visibility   constant.PostVisibility `json:"visibility"`
	VisibleRoles []RoleOptionResponse    `json:"visible_roles"`
	ViewCount    uint64                  `json:"view_count"`
	LikeCount    uint64                  `json:"like_count"`
	CommentCount uint64                  `json:"comment_count"`
	Status       string                  `json:"status"`
	PublishedAt  *time.Time              `json:"published_at"`
	CreatedAt    time.Time               `json:"created_at"`
	UpdatedAt    time.Time               `json:"updated_at"`
}

type ListDiariesResponse struct {
	Items    []DiaryResponse `json:"items"`
	Total    int64           `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
}
