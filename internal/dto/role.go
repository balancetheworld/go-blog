package dto

import (
	"time"

	"github.com/zyj/my-blog/pkg/constant"
)

type RoleOptionResponse struct {
	ID          uint   `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RoleResponse struct {
	ID            uint      `json:"id"`
	Code          string    `json:"code"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	IsSystem      bool      `json:"is_system"`
	IsDefault     bool      `json:"is_default"`
	IsRequestable bool      `json:"is_requestable"`
	Enabled       bool      `json:"enabled"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreateRoleRequest struct {
	Code          string `json:"code" vd:"len($) >= 1 && len($) <= 64"`
	Name          string `json:"name" vd:"len($) >= 1 && len($) <= 64"`
	Description   string `json:"description" vd:"len($) <= 512"`
	IsRequestable bool   `json:"is_requestable"`
	Enabled       *bool  `json:"enabled"`
}

type UpdateRoleRequest struct {
	Name          *string `json:"name" vd:"$ == nil || (len($) >= 1 && len($) <= 64)"`
	Description   *string `json:"description" vd:"$ == nil || len($) <= 512"`
	IsRequestable *bool   `json:"is_requestable"`
	Enabled       *bool   `json:"enabled"`
}

type ListRolesRequest struct {
	Page     int    `query:"page" vd:"$ == 0 || $ >= 1"`
	PageSize int    `query:"page_size" vd:"$ == 0 || ($ >= 1 && $ <= 100)"`
	Keyword  string `query:"keyword"`
}

type ListRolesResponse struct {
	Items    []RoleResponse `json:"items"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}

type RoleApplicationUserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

type RoleApplicationResponse struct {
	ID            uint                           `json:"id"`
	User          RoleApplicationUserResponse    `json:"user"`
	RequestedRole RoleOptionResponse             `json:"requested_role"`
	Status        constant.RoleApplicationStatus `json:"status"`
	ReviewerID    *uint                          `json:"reviewer_id"`
	ReviewedAt    *time.Time                     `json:"reviewed_at"`
	RejectReason  string                         `json:"reject_reason"`
	CreatedAt     time.Time                      `json:"created_at"`
}

type ListRoleApplicationsRequest struct {
	Page     int                            `query:"page" vd:"$ == 0 || $ >= 1"`
	PageSize int                            `query:"page_size" vd:"$ == 0 || ($ >= 1 && $ <= 100)"`
	Status   constant.RoleApplicationStatus `query:"status" vd:"$ == '' || in($, 'pending', 'approved', 'rejected')"`
}

type ListRoleApplicationsResponse struct {
	Items    []RoleApplicationResponse `json:"items"`
	Total    int64                     `json:"total"`
	Page     int                       `json:"page"`
	PageSize int                       `json:"page_size"`
}

type RejectRoleApplicationRequest struct {
	Reason string `json:"reason" vd:"len($) >= 1 && len($) <= 512"`
}
