package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/zyj/my-blog/internal/model"
	"github.com/zyj/my-blog/internal/repo"
	"github.com/zyj/my-blog/pkg/constant"
	"github.com/zyj/my-blog/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func EnsureRootUser(ctx context.Context) error {
	adminRole, err := repo.GetRoleByCode(ctx, constant.RoleCodeAdmin)
	if err != nil {
		return fmt.Errorf("get admin role: %w", err)
	}

	rootUser, err := repo.GetRootUser(ctx)
	if err == nil {
		return repo.RestoreRootUser(ctx, rootUser.ID, adminRole.ID)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("get root user: %w", err)
	}

	username := strings.TrimSpace(utils.Get(constant.EnvKeyRootAdminUsername))
	email := strings.TrimSpace(utils.Get(constant.EnvKeyRootAdminEmail))
	password := utils.Get(constant.EnvKeyRootAdminPassword)
	if username == "" || email == "" || password == "" {
		return errors.New("root admin username, email and password are required")
	}
	if len(username) < 3 || len(username) > 32 {
		return errors.New("root admin username length must be between 3 and 32")
	}
	if len(password) < 8 || len(password) > 72 {
		return errors.New("root admin password length must be between 8 and 72")
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return fmt.Errorf("hash root admin password: %w", err)
	}

	rootUser = model.User{
		Username:     username,
		Nickname:     username,
		Email:        email,
		PasswordHash: string(passwordHash),
		Role:         constant.RoleAdmin,
		RoleID:       &adminRole.ID,
		IsRoot:       true,
	}

	if err := repo.CreateUser(ctx, &rootUser); err != nil {
		return fmt.Errorf("create root user: %w", err)
	}

	return nil
}
