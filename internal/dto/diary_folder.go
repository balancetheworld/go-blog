package dto

import (
	"time"

	"github.com/zyj/my-blog/pkg/constant"
)

type CreateDiaryFolderRequest struct {
	Name           string                  `json:"name" vd:"len($) >= 1 && len($) <= 100"`
	Slug           string                  `json:"slug" vd:"len($) <= 100"`
	Description    string                  `json:"description" vd:"len($) <= 500"`
	Cover          string                  `json:"cover" vd:"len($) <= 500"`
	Sort           int                     `json:"sort"`
	Visibility     constant.PostVisibility `json:"visibility" vd:"$ == '' || in($, 'public', 'roles', 'private')"`
	VisibleRoleIDs []uint                  `json:"visible_role_ids"`
}

type UpdateDiaryFolderRequest struct {
	Name           *string                  `json:"name" vd:"$ == nil || (len($) >= 1 && len($) <= 100)"`
	Slug           *string                  `json:"slug" vd:"$ == nil || len($) <= 100"`
	Description    *string                  `json:"description" vd:"$ == nil || len($) <= 500"`
	Cover          *string                  `json:"cover" vd:"$ == nil || len($) <= 500"`
	Sort           *int                     `json:"sort"`
	Visibility     *constant.PostVisibility `json:"visibility" vd:"$ == nil || in($, 'public', 'roles', 'private')"`
	VisibleRoleIDs *[]uint                  `json:"visible_role_ids"`
}

type ListDiaryFoldersRequest struct {
	All bool `query:"all"`
}

type DiaryFolderResponse struct {
	ID           uint64                  `json:"id"`
	Name         string                  `json:"name"`
	Slug         string                  `json:"slug"`
	Description  string                  `json:"description"`
	Cover        string                  `json:"cover"`
	Sort         int                     `json:"sort"`
	Visibility   constant.PostVisibility `json:"visibility"`
	VisibleRoles []RoleOptionResponse    `json:"visible_roles"`
	CreatedAt    time.Time               `json:"created_at"`
	UpdatedAt    time.Time               `json:"updated_at"`
}

type ListDiaryFoldersResponse struct {
	Items []DiaryFolderResponse `json:"items"`
}
