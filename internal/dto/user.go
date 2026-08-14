package dto

import (
	"time"

	"github.com/zyj/my-blog/pkg/constant"
)

type CreateUserRequest struct {
	Username string `json:"username" vd:"len($) >= 3 && len($) <= 32"`
	Email    string `json:"email" vd:"email($)"`
	Password string `json:"password" vd:"len($) >= 8 && len($) <= 72"`
	Nickname string `json:"nickname,omitempty" vd:"len($) <= 64"`
}

type UpdateUserRequest struct {
	Username *string `json:"username,omitempty" vd:"$ == nil || (len($) >= 3 && len($) <= 32)"`
	Email    *string `json:"email,omitempty" vd:"$ == nil || email($)"`
	Nickname *string `json:"nickname,omitempty" vd:"$ == nil || len($) <= 64"`
	Avatar   *string `json:"avatar,omitempty"`
}

type UpdateUserRoleRequest struct {
	Role constant.Role `json:"role" vd:"in($, 'user', 'editor', 'admin')"`
}

type ListUsersRequest struct {
	Page     int `query:"page" vd:"$ == 0 || $ >= 1"`
	PageSize int `query:"page_size" vd:"$ == 0 || ($ >= 1 && $ <= 100)"`
}

type UserResponse struct {
	ID        uint64        `json:"id"`
	Username  string        `json:"username"`
	Nickname  string        `json:"nickname"`
	Avatar    string        `json:"avatar"`
	Role      constant.Role `json:"role"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type UserPrivateResponse struct {
	UserResponse
	Email string `json:"email"`
}

type UserListResponse struct {
	Items    []UserResponse `json:"items"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}

type UserDto = UserResponse

type UserLoginReq struct {
	Account   string `json:"account" vd:"len($) >= 3"`
	Password  string `json:"password" vd:"len($) >= 8 && len($) <= 72"`
	Remember  bool   `json:"remember_me"`
	UserIP    string `json:"-"`
	UserAgent string `json:"-"`
}

type UserAuthResponse struct {
	User          UserPrivateResponse `json:"user"`
	AccessToken   string              `json:"-"`
	RefreshToken  string              `json:"-"`
	AccessMaxAge  int                 `json:"-"`
	RefreshMaxAge int                 `json:"-"`
}

type UserRegisterReq struct {
	Username        string `json:"username" vd:"len($) >= 3 && len($) <= 32"`
	Email           string `json:"email" vd:"email($)"`
	Password        string `json:"password" vd:"len($) >= 8 && len($) <= 72"`
	Nickname        string `json:"nickname,omitempty" vd:"len($) <= 64"`
	Code            string `json:"email_code,omitempty"`
	Remember        bool   `json:"remember_me"`
	UserIP          string `json:"-"`
	UserAgent       string `json:"-"`
	RequestedRoleID *uint  `json:"requested_role_id,omitempty" vd:"$ == nil || $ > 0"`
}

type VerifyEmailReq struct {
	Email  string `json:"email" vd:"email($)"`
	UserIP string `json:"-"`
}

type UpdatePasswordReq struct {
	OldPassword string `json:"old_password" vd:"len($) >= 8 && len($) <= 72"`
	NewPassword string `json:"new_password" vd:"len($) >= 8 && len($) <= 72"`
}

type ResetPasswordReq struct {
	Email       string `json:"email" vd:"email($)"`
	Code        string `json:"code" vd:"len($) == 6"`
	NewPassword string `json:"new_password" vd:"len($) >= 8 && len($) <= 72"`
}

type UpdateEmailReq struct {
	Email string `json:"email" vd:"email($)"`
	Code  string `json:"code" vd:"len($) == 6"`
}

type GetUserByUsernameReq struct {
	Username string `query:"username" vd:"len($) >= 3 && len($) <= 32"`
}

type UserRoleApplicationResponse struct {
	ID            uint                           `json:"id"`
	RequestedRole RoleOptionResponse             `json:"requested_role"`
	Status        constant.RoleApplicationStatus `json:"status"`
	RejectReason  string                         `json:"reject_reason"`
	ReviewedAt    *time.Time                     `json:"reviewed_at"`
	CreatedAt     time.Time                      `json:"created_at"`
}

type LoginUserResponse struct {
	User            *UserResponse                `json:"user"`
	Role            constant.Role                `json:"role"`
	Identity        *RoleOptionResponse          `json:"identity"`
	RoleApplication *UserRoleApplicationResponse `json:"role_application"`
}

type CaptchaConfigResponse struct {
	Provider constant.CaptchaType `json:"provider"`
	SiteKey  string               `json:"site_key"`
}
