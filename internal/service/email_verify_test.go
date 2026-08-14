package service

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/zyj/my-blog/internal/dto"
	"github.com/zyj/my-blog/internal/model"
	"github.com/zyj/my-blog/internal/repo"
	"github.com/zyj/my-blog/pkg/constant"
)

func setupEmailVerifyServiceTest(t *testing.T, databaseName string) (*miniredis.Miniredis, context.Context) {
	t.Helper()

	_ = repo.CloseRedis()
	_ = repo.CloseDatabase()
	server := miniredis.RunT(t)
	t.Setenv(constant.EnvKeyDBDriver, "sqlite")
	t.Setenv(
		constant.EnvKeyDBSQLitePath,
		filepath.Join(t.TempDir(), databaseName),
	)
	t.Setenv(constant.EnvKeyRedisAddr, server.Addr())
	t.Setenv(constant.EnvKeyRedisPassword, "")
	t.Setenv(constant.EnvKeyRedisDB, "0")
	t.Setenv(constant.EnvKeyEnableEmailVerify, "true")
	t.Setenv(constant.EnvKeyEnableRegister, "true")
	t.Setenv(constant.EnvKeyJWTSecret, "0123456789abcdef0123456789abcdef")
	t.Setenv(constant.EnvKeyTokenDuration, "3600")
	t.Setenv(constant.EnvKeyRefreshTokenDuration, "604800")
	t.Setenv(constant.EnvKeyRefreshTokenDurationWithRemember, "2592000")

	ctx := context.Background()
	if err := repo.InitDatabase(); err != nil {
		t.Fatal(err)
	}
	if err := repo.InitRedis(ctx); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = repo.CloseRedis()
		_ = repo.CloseDatabase()
	})

	return server, ctx
}

func saveServiceVerificationCode(t *testing.T, ctx context.Context, email, userIP string) {
	t.Helper()

	saved, err := repo.SaveEmailVerifyCode(ctx, email, "123456", userIP)
	if err != nil {
		t.Fatal(err)
	}
	if !saved {
		t.Fatal("expected verification code to be saved")
	}
}

func registrationRequest(username, email, userIP string) dto.UserRegisterReq {
	return dto.UserRegisterReq{
		Username:  username,
		Email:     email,
		Password:  "password123",
		Nickname:  username,
		Code:      "123456",
		UserIP:    userIP,
		UserAgent: "service-test",
	}
}

func TestResetPasswordReleasesVerificationCodeAfterDatabaseFailure(t *testing.T) {
	_, ctx := setupEmailVerifyServiceTest(t, "reset-password.db")

	email := "reset-password-release@example.com"
	userIP := "reset-password-release"

	memberRole, err := repo.GetRoleByCode(ctx, constant.RoleCodeMember)
	if err != nil {
		t.Fatal(err)
	}
	user := model.User{
		Username:     "reset-password-release",
		Email:        email,
		PasswordHash: "existing-password-hash",
		Role:         constant.RoleUser,
		RoleID:       &memberRole.ID,
	}
	if err := repo.CreateUser(ctx, &user); err != nil {
		t.Fatal(err)
	}

	saveServiceVerificationCode(t, ctx, email, userIP)

	callbackName := "test:fail_reset_password_update"
	callbackRegistered := true
	t.Cleanup(func() {
		if callbackRegistered {
			_ = repo.GetDB().Callback().Update().Remove(callbackName)
		}
	})
	if err := repo.GetDB().Callback().Update().Before("gorm:update").Register(
		callbackName,
		func(tx *gorm.DB) {
			if tx.Statement.Table == "users" {
				tx.AddError(errors.New("forced password update failure"))
			}
		},
	); err != nil {
		t.Fatal(err)
	}

	req := dto.ResetPasswordReq{
		Email:       email,
		Code:        "123456",
		NewPassword: "new-password",
	}
	err = ResetPassword(ctx, &req)
	requireServiceStatus(t, err, 500)

	if err := repo.GetDB().Callback().Update().Remove(callbackName); err != nil {
		t.Fatal(err)
	}
	callbackRegistered = false

	if err := ResetPassword(ctx, &req); err != nil {
		t.Fatal(err)
	}

	updatedUser, err := repo.GetUserByEmail(ctx, email)
	if err != nil {
		t.Fatal(err)
	}
	if err := bcrypt.CompareHashAndPassword(
		[]byte(updatedUser.PasswordHash),
		[]byte(req.NewPassword),
	); err != nil {
		t.Fatal(err)
	}

	reserved, err := repo.ReserveEmailVerifyCode(ctx, email, "123456", "replay-token")
	if err != nil {
		t.Fatal(err)
	}
	if reserved {
		t.Fatal("expected committed verification code to be unavailable")
	}
}

func TestRegisterReleasesVerificationCodeAfterDatabaseFailure(t *testing.T) {
	_, ctx := setupEmailVerifyServiceTest(t, "register-database-failure.db")
	email := "register-database-failure@example.com"
	req := registrationRequest("register-db-failure", email, "register-db-failure")
	saveServiceVerificationCode(t, ctx, email, req.UserIP)

	callbackName := "test:fail_register_user_create"
	callbackRegistered := true
	t.Cleanup(func() {
		if callbackRegistered {
			_ = repo.GetDB().Callback().Create().Remove(callbackName)
		}
	})
	if err := repo.GetDB().Callback().Create().Before("gorm:create").Register(
		callbackName,
		func(tx *gorm.DB) {
			if tx.Statement.Table == "users" {
				tx.AddError(errors.New("forced user create failure"))
			}
		},
	); err != nil {
		t.Fatal(err)
	}

	_, err := UserRegister(ctx, &req)
	requireServiceStatus(t, err, 500)

	if err := repo.GetDB().Callback().Create().Remove(callbackName); err != nil {
		t.Fatal(err)
	}
	callbackRegistered = false

	if _, err := UserRegister(ctx, &req); err != nil {
		t.Fatal(err)
	}
	reserved, err := repo.ReserveEmailVerifyCode(ctx, email, "123456", "replay-token")
	if err != nil {
		t.Fatal(err)
	}
	if reserved {
		t.Fatal("expected successful registration to consume verification code")
	}
}

func TestRegisterRollsBackUserAndSessionAfterTokenFailure(t *testing.T) {
	_, ctx := setupEmailVerifyServiceTest(t, "register-token-failure.db")
	email := "register-token-failure@example.com"
	req := registrationRequest("register-token-failure", email, "register-token-failure")
	saveServiceVerificationCode(t, ctx, email, req.UserIP)
	t.Setenv(constant.EnvKeyJWTSecret, "short")

	_, err := UserRegister(ctx, &req)
	requireServiceStatus(t, err, 500)

	var userCount int64
	if err := repo.GetDB().Model(&model.User{}).
		Where("username = ?", req.Username).
		Count(&userCount).
		Error; err != nil {
		t.Fatal(err)
	}
	if userCount != 0 {
		t.Fatalf("expected user rollback, got %d users", userCount)
	}
	var sessionCount int64
	if err := repo.GetDB().Model(&model.Session{}).
		Where("session_id <> ?", "").
		Count(&sessionCount).
		Error; err != nil {
		t.Fatal(err)
	}
	if sessionCount != 0 {
		t.Fatalf("expected session rollback, got %d sessions", sessionCount)
	}

	t.Setenv(constant.EnvKeyJWTSecret, "0123456789abcdef0123456789abcdef")
	if _, err := UserRegister(ctx, &req); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateEmailReleasesVerificationCodeAfterDatabaseFailure(t *testing.T) {
	_, ctx := setupEmailVerifyServiceTest(t, "update-email-failure.db")
	memberRole, err := repo.GetRoleByCode(ctx, constant.RoleCodeMember)
	if err != nil {
		t.Fatal(err)
	}
	user := model.User{
		Username:     "update-email-failure",
		Email:        "old-email@example.com",
		PasswordHash: "existing-password-hash",
		Role:         constant.RoleUser,
		RoleID:       &memberRole.ID,
	}
	if err := repo.CreateUser(ctx, &user); err != nil {
		t.Fatal(err)
	}

	newEmail := "new-email@example.com"
	saveServiceVerificationCode(t, ctx, newEmail, "update-email-failure")
	callbackName := "test:fail_update_email"
	callbackRegistered := true
	t.Cleanup(func() {
		if callbackRegistered {
			_ = repo.GetDB().Callback().Update().Remove(callbackName)
		}
	})
	if err := repo.GetDB().Callback().Update().Before("gorm:update").Register(
		callbackName,
		func(tx *gorm.DB) {
			if tx.Statement.Table == "users" {
				tx.AddError(errors.New("forced email update failure"))
			}
		},
	); err != nil {
		t.Fatal(err)
	}

	req := dto.UpdateEmailReq{Email: newEmail, Code: "123456"}
	err = UpdateEmail(ctx, user.ID, &req)
	requireServiceStatus(t, err, 500)

	if err := repo.GetDB().Callback().Update().Remove(callbackName); err != nil {
		t.Fatal(err)
	}
	callbackRegistered = false

	if err := UpdateEmail(ctx, user.ID, &req); err != nil {
		t.Fatal(err)
	}
	updated, err := repo.GetUserByID(ctx, uint64(user.ID))
	if err != nil {
		t.Fatal(err)
	}
	if updated.Email != newEmail {
		t.Fatalf("expected email %s, got %s", newEmail, updated.Email)
	}
}

func TestRegisterSucceedsWhenVerificationCommitFails(t *testing.T) {
	server, ctx := setupEmailVerifyServiceTest(t, "register-commit-failure.db")
	email := "register-commit-failure@example.com"
	req := registrationRequest("register-commit-failure", email, "register-commit-failure")
	saveServiceVerificationCode(t, ctx, email, req.UserIP)

	callbackName := "test:fail_verification_commit"
	if err := repo.GetDB().Callback().Create().After("gorm:create").Register(
		callbackName,
		func(tx *gorm.DB) {
			if tx.Statement.Table == "sessions" {
				server.SetError("LOADING Redis is loading the dataset in memory")
			}
		},
	); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		server.SetError("")
		_ = repo.GetDB().Callback().Create().Remove(callbackName)
	})

	result, err := UserRegister(ctx, &req)
	if err != nil {
		t.Fatal(err)
	}
	if result.User.Username != req.Username || result.AccessToken == "" || result.RefreshToken == "" {
		t.Fatalf("unexpected registration result: %#v", result)
	}
}
