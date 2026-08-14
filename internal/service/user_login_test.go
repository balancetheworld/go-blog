package service

import (
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/zyj/my-blog/internal/dto"
	"github.com/zyj/my-blog/internal/model"
	"github.com/zyj/my-blog/internal/repo"
	"github.com/zyj/my-blog/pkg/constant"
)

func TestUserLoginReturnsTooManyRequestsAtAccountLimit(t *testing.T) {
	_, ctx := setupEmailVerifyServiceTest(t, "login-attempt-limit.db")
	memberRole, err := repo.GetRoleByCode(ctx, constant.RoleCodeMember)
	if err != nil {
		t.Fatal(err)
	}
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte("correct-password"),
		bcrypt.DefaultCost,
	)
	if err != nil {
		t.Fatal(err)
	}
	user := model.User{
		Username:     "login-attempt-limit",
		Email:        "login-attempt-limit@example.com",
		PasswordHash: string(passwordHash),
		Role:         constant.RoleUser,
		RoleID:       &memberRole.ID,
	}
	if err := repo.CreateUser(ctx, &user); err != nil {
		t.Fatal(err)
	}

	req := dto.UserLoginReq{
		Account:  user.Username,
		Password: "wrong-password",
		UserIP:   "192.0.2.30",
	}
	for attempt := 1; attempt <= 5; attempt++ {
		_, err := UserLogin(ctx, &req)
		if attempt < 5 {
			requireServiceStatus(t, err, 401)
		} else {
			requireServiceStatus(t, err, 429)
		}
	}
}
