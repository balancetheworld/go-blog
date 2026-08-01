package service

import (
	"context"
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/zyj/my-blog/internal/dto"
	"github.com/zyj/my-blog/internal/model"
	"github.com/zyj/my-blog/internal/repo"
	"github.com/zyj/my-blog/pkg/constant"
	"github.com/zyj/my-blog/pkg/errs"
)

func ListUsers(ctx context.Context, req dto.ListUsersRequest) (dto.UserListResponse, error) {
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	users, total, err := repo.ListUsers(ctx, (req.Page-1)*req.PageSize, req.PageSize)
	if err != nil {
		return dto.UserListResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"list users failed",
		)
	}

	items := make([]dto.UserResponse, 0, len(users))
	for _, user := range users {
		items = append(items, toUserResponse(user))
	}

	return dto.UserListResponse{
		Items:    items,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

func GetUser(ctx context.Context, id uint64) (dto.UserResponse, error) {
	user, err := repo.GetUserByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.UserResponse{}, errs.NewNotFound(
			http.StatusNotFound,
			"user not found",
		)
	}
	if err != nil {
		return dto.UserResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"get user failed",
		)
	}

	return toUserResponse(user), nil
}

func CreateUser(ctx context.Context, req dto.CreateUserRequest) (dto.UserPrivateResponse, error) {
	// 1. 参数校验已经在 controller 层完成，这里不需要重复校验
	// 2. 业务逻辑处理：调用 repo 层，保存用户数据到数据库
	// 3. 返回结果给 controller 层，controller 层再返回给前端
	exists, err := repo.UserExists(ctx, req.Username, req.Email)
	if err != nil {
		//空的**响应 DTO 结构体**。规范：哪怕出错，也要把出参 DTO 占位返回，避免上层调用方结构体解析报错。
		return dto.UserPrivateResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"query user exists failed",
		)
	}
	if exists {
		return dto.UserPrivateResponse{}, errs.NewConflict(
			http.StatusConflict,
			"user already exists",
		)
	}
	//使用 bcrypt 算法给前端传来的明文密码生成加密哈希，自动加盐，安全保存到数据库，防止密码泄露。
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.UserPrivateResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"hash password failed",
		)
	}

	user := model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(passwordHash),
		Nickname:     req.Nickname,
		Role:         constant.RoleUser,
	}

	if err := repo.CreateUser(ctx, &user); err != nil {
		return dto.UserPrivateResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"create user failed",
		)
	}
	return toUserPrivateResponse(user), nil
}

func UpdateUser(ctx context.Context, id uint64, req dto.UpdateUserRequest) (dto.UserPrivateResponse, error) {
	user, err := repo.GetUserByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.UserPrivateResponse{}, errs.NewNotFound(
			http.StatusNotFound,
			"user not found",
		)
	}
	if err != nil {
		return dto.UserPrivateResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"get user failed",
		)
	}

	if req.Username != nil {
		user.Username = *req.Username
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Nickname != nil {
		user.Nickname = *req.Nickname
	}
	if req.Avatar != nil {
		user.Avatar = *req.Avatar
	}

	exists, err := repo.UserExistsExcept(ctx, id, user.Username, user.Email)
	if err != nil {
		return dto.UserPrivateResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"query user exists failed",
		)
	}
	if exists {
		return dto.UserPrivateResponse{}, errs.NewConflict(
			http.StatusConflict,
			"username or email already exists",
		)
	}

	if err := repo.UpdateUser(ctx, &user); err != nil {
		return dto.UserPrivateResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"update user failed",
		)
	}

	return toUserPrivateResponse(user), nil
}

func DeleteUser(ctx context.Context, id uint64) error {
	rowsAffected, err := repo.DeleteUser(ctx, id)
	if err != nil {
		return errs.NewInternalServer(
			http.StatusInternalServerError,
			"delete user failed",
		)
	}
	if rowsAffected == 0 {
		return errs.NewNotFound(
			http.StatusNotFound,
			"user not found",
		)
	}

	return nil
}

// 转换函数确保数据库模型中的密码不会进入响应
// 专门写转换函数，**隔离数据库 model 和前端返回 DTO，彻底阻止密码、敏感字段泄露给前端
// model.User 是数据库模型   dto.UserResp 是前端响应 dto
func toUserResponse(user model.User) dto.UserResponse {
	return dto.UserResponse{
		ID: user.ID,
		Username:  user.Username,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func toUserPrivateResponse(user model.User) dto.UserPrivateResponse {
	return dto.UserPrivateResponse{
                UserResponse: toUserResponse(user),
                Email:        user.Email,
        }
}
