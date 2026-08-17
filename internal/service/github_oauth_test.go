package service

import (
	"testing"

	"github.com/zyj/my-blog/internal/model"
	"github.com/zyj/my-blog/internal/repo"
	"github.com/zyj/my-blog/pkg/constant"
)

func TestGithubLoginCreatesAndReusesUser(t *testing.T) {
	_, ctx := setupEmailVerifyServiceTest(t, "github-login.db")
	profile := githubProfile{
		ID:        123456,
		Login:     "github-user",
		Name:      "GitHub User",
		AvatarURL: "https://avatars.githubusercontent.com/u/123456",
		Email:     "GitHub-User@Example.com",
	}

	first, err := loginGithubUser(ctx, profile, "192.0.2.1", "service-test")
	if err != nil {
		t.Fatal(err)
	}
	if first.User.Username != "github-123456" || first.User.Email != "github-user@example.com" {
		t.Fatalf("unexpected first login user: %#v", first.User)
	}

	user, err := repo.GetUserByGithubID(ctx, profile.ID)
	if err != nil {
		t.Fatal(err)
	}
	if user.GithubID == nil || *user.GithubID != profile.ID {
		t.Fatalf("unexpected github id: %#v", user.GithubID)
	}

	second, err := loginGithubUser(ctx, profile, "192.0.2.2", "service-test")
	if err != nil {
		t.Fatal(err)
	}
	if second.User.ID != first.User.ID {
		t.Fatalf("expected same user, got %d and %d", first.User.ID, second.User.ID)
	}
}

func TestGithubLoginRejectsExistingEmail(t *testing.T) {
	_, ctx := setupEmailVerifyServiceTest(t, "github-email-conflict.db")
	memberRole, err := repo.GetRoleByCode(ctx, constant.RoleCodeMember)
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.CreateUser(ctx, &model.User{
		Username:     "existing-email-user",
		Email:        "existing@example.com",
		PasswordHash: "password-hash",
		Role:         constant.RoleUser,
		RoleID:       &memberRole.ID,
	}); err != nil {
		t.Fatal(err)
	}

	_, err = loginGithubUser(ctx, githubProfile{
		ID:    654321,
		Login: "another-user",
		Email: "existing@example.com",
	}, "192.0.2.3", "service-test")
	requireServiceStatus(t, err, 409)
}
