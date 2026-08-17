package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/zyj/my-blog/internal/dto"
	"github.com/zyj/my-blog/internal/model"
	"github.com/zyj/my-blog/internal/repo"
	"github.com/zyj/my-blog/pkg/constant"
	"github.com/zyj/my-blog/pkg/errs"
	"github.com/zyj/my-blog/pkg/utils"
)

const (
	githubOAuthAuthorizeURL = "https://github.com/login/oauth/authorize"
	githubOAuthTokenURL     = "https://github.com/login/oauth/access_token"
	githubAPIUserURL        = "https://api.github.com/user"
	githubAPIEmailsURL      = "https://api.github.com/user/emails"
)

var githubOAuthHTTPClient = &http.Client{Timeout: 10 * time.Second}

type githubOAuthConfig struct {
	clientID     string
	clientSecret string
	redirectURL  string
}

type githubOAuthTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type githubProfile struct {
	ID        uint64 `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	Email     string `json:"email"`
}

type githubEmail struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

func GetGithubAuthorizationURL(ctx context.Context) (string, error) {
	config, err := getGithubOAuthConfig()
	if err != nil {
		return "", err
	}

	state, err := utils.GenerateOAuthState()
	if err != nil {
		return "", err
	}
	saved, err := repo.SaveGithubOAuthState(ctx, state)
	if err != nil {
		return "", err
	}
	if !saved {
		return "", errors.New("save github oauth state failed")
	}

	values := url.Values{
		"client_id":    {config.clientID},
		"redirect_uri": {config.redirectURL},
		"scope":        {"read:user user:email"},
		"state":        {state},
	}
	return githubOAuthAuthorizeURL + "?" + values.Encode(), nil
}

func CompleteGithubOAuthLogin(
	ctx context.Context,
	code string,
	state string,
	userIP string,
	userAgent string,
) (dto.UserAuthResponse, error) {
	if strings.TrimSpace(code) == "" || strings.TrimSpace(state) == "" {
		return dto.UserAuthResponse{}, errs.NewBadRequest(
			http.StatusBadRequest,
			"invalid github oauth callback",
		)
	}

	valid, err := repo.ConsumeGithubOAuthState(ctx, state)
	if err != nil {
		return dto.UserAuthResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"consume github oauth state failed",
		)
	}
	if !valid {
		return dto.UserAuthResponse{}, errs.NewUnauthorized(
			http.StatusUnauthorized,
			"invalid github oauth state",
		)
	}

	config, err := getGithubOAuthConfig()
	if err != nil {
		return dto.UserAuthResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"github oauth configuration is incomplete",
		)
	}

	accessToken, err := exchangeGithubOAuthCode(ctx, config, code)
	if err != nil {
		return dto.UserAuthResponse{}, errs.NewUnauthorized(
			http.StatusUnauthorized,
			"github oauth authorization failed",
		)
	}
	profile, err := getGithubProfile(ctx, accessToken)
	if err != nil {
		return dto.UserAuthResponse{}, errs.NewUnauthorized(
			http.StatusUnauthorized,
			"get github profile failed",
		)
	}

	return loginGithubUser(ctx, profile, userIP, userAgent)
}

func loginGithubUser(
	ctx context.Context,
	profile githubProfile,
	userIP string,
	userAgent string,
) (dto.UserAuthResponse, error) {
	if profile.ID == 0 || strings.TrimSpace(profile.Email) == "" {
		return dto.UserAuthResponse{}, errs.NewUnauthorized(
			http.StatusUnauthorized,
			"github account does not have a verified email",
		)
	}

	user, err := repo.GetUserByGithubID(ctx, profile.ID)
	if err == nil {
		return createGithubSession(ctx, user, userIP, userAgent)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.UserAuthResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"get github user failed",
		)
	}

	if !utils.GetAsBool(constant.EnvKeyEnableRegister, true) {
		return dto.UserAuthResponse{}, errs.NewForbidden(
			http.StatusForbidden,
			"registration is disabled",
		)
	}

	username := fmt.Sprintf("github-%d", profile.ID)
	usernameExists, err := repo.CheckUsernameExists(ctx, username)
	if err != nil {
		return dto.UserAuthResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"check github username failed",
		)
	}
	if usernameExists {
		return dto.UserAuthResponse{}, errs.NewConflict(
			http.StatusConflict,
			"github account cannot be registered",
		)
	}

	email := strings.ToLower(strings.TrimSpace(profile.Email))
	emailExists, err := repo.CheckEmailExists(ctx, email)
	if err != nil {
		return dto.UserAuthResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"check github email failed",
		)
	}
	if emailExists {
		return dto.UserAuthResponse{}, errs.NewConflict(
			http.StatusConflict,
			"github email is already registered",
		)
	}

	password, err := utils.GenerateSessionID()
	if err != nil {
		return dto.UserAuthResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"generate github password failed",
		)
	}
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return dto.UserAuthResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"hash github password failed",
		)
	}

	githubID := profile.ID
	verifiedAt := time.Now().UTC()
	user = model.User{
		Username:        username,
		Nickname:        githubNickname(profile),
		AvatarURL:       strings.TrimSpace(profile.AvatarURL),
		Email:           email,
		EmailVerifiedAt: &verifiedAt,
		GithubID:        &githubID,
		PasswordHash:    string(passwordHash),
		Role:            constant.RoleUser,
	}

	sessionID, err := utils.GenerateSessionID()
	if err != nil {
		return dto.UserAuthResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"generate session failed",
		)
	}
	session := model.Session{
		SessionID: sessionID,
		UserIP:    userIP,
		UserAgent: userAgent,
	}

	var tokens utils.TokenPair
	var tokenErr error
	if err := repo.RegisterUserWithSession(
		ctx,
		&user,
		&session,
		nil,
		func() error {
			tokens, tokenErr = utils.GenerateTokenPair(
				user.ID,
				user.Role,
				sessionID,
				false,
			)
			return tokenErr
		},
	); err != nil {
		if tokenErr != nil {
			return dto.UserAuthResponse{}, errs.NewInternalServer(
				http.StatusInternalServerError,
				"generate token failed",
			)
		}
		return dto.UserAuthResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"register github user failed",
		)
	}

	return dto.UserAuthResponse{
		User:          toUserPrivateResponse(user),
		AccessToken:   tokens.AccessToken,
		RefreshToken:  tokens.RefreshToken,
		AccessMaxAge:  tokens.AccessMaxAge,
		RefreshMaxAge: tokens.RefreshMaxAge,
	}, nil
}

func createGithubSession(
	ctx context.Context,
	user model.User,
	userIP string,
	userAgent string,
) (dto.UserAuthResponse, error) {
	sessionID, err := utils.GenerateSessionID()
	if err != nil {
		return dto.UserAuthResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"generate session failed",
		)
	}
	tokens, err := utils.GenerateTokenPair(user.ID, user.Role, sessionID, false)
	if err != nil {
		return dto.UserAuthResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"generate token failed",
		)
	}
	if err := repo.CreateSession(ctx, &model.Session{
		UserID:    user.ID,
		SessionID: sessionID,
		UserIP:    userIP,
		UserAgent: userAgent,
	}); err != nil {
		return dto.UserAuthResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"create session failed",
		)
	}

	return dto.UserAuthResponse{
		User:          toUserPrivateResponse(user),
		AccessToken:   tokens.AccessToken,
		RefreshToken:  tokens.RefreshToken,
		AccessMaxAge:  tokens.AccessMaxAge,
		RefreshMaxAge: tokens.RefreshMaxAge,
	}, nil
}

func getGithubOAuthConfig() (githubOAuthConfig, error) {
	config := githubOAuthConfig{
		clientID:     strings.TrimSpace(utils.Get(constant.EnvKeyGithubClientID)),
		clientSecret: strings.TrimSpace(utils.Get(constant.EnvKeyGithubClientSecret)),
		redirectURL:  strings.TrimSpace(utils.Get(constant.EnvKeyGithubRedirectURL)),
	}
	if config.clientID == "" || config.clientSecret == "" || config.redirectURL == "" {
		return githubOAuthConfig{}, errors.New("github oauth configuration is incomplete")
	}
	return config, nil
}

func exchangeGithubOAuthCode(
	ctx context.Context,
	config githubOAuthConfig,
	code string,
) (string, error) {
	body, err := json.Marshal(map[string]string{
		"client_id":     config.clientID,
		"client_secret": config.clientSecret,
		"code":          code,
		"redirect_uri":  config.redirectURL,
	})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		githubOAuthTokenURL,
		bytes.NewReader(body),
	)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	response, err := githubOAuthHTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github token response status: %d", response.StatusCode)
	}

	var result githubOAuthTokenResponse
	if err := json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(&result); err != nil {
		return "", err
	}
	if strings.TrimSpace(result.AccessToken) == "" {
		return "", errors.New("github access token is empty")
	}
	return result.AccessToken, nil
}

func getGithubProfile(ctx context.Context, accessToken string) (githubProfile, error) {
	profile, err := getGithubProfileResponse(ctx, accessToken)
	if err != nil {
		return githubProfile{}, err
	}
	emails, err := getGithubEmails(ctx, accessToken)
	if err != nil {
		return githubProfile{}, err
	}
	for _, email := range emails {
		if email.Primary && email.Verified {
			profile.Email = email.Email
			return profile, nil
		}
	}
	for _, email := range emails {
		if email.Verified {
			profile.Email = email.Email
			return profile, nil
		}
	}
	return githubProfile{}, errors.New("github verified email is unavailable")
}

func getGithubProfileResponse(ctx context.Context, accessToken string) (githubProfile, error) {
	request, err := githubAPIRequest(ctx, githubAPIUserURL, accessToken)
	if err != nil {
		return githubProfile{}, err
	}
	response, err := githubOAuthHTTPClient.Do(request)
	if err != nil {
		return githubProfile{}, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return githubProfile{}, fmt.Errorf("github user response status: %d", response.StatusCode)
	}

	var profile githubProfile
	if err := json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(&profile); err != nil {
		return githubProfile{}, err
	}
	if profile.ID == 0 {
		return githubProfile{}, errors.New("github user id is empty")
	}
	return profile, nil
}

func getGithubEmails(ctx context.Context, accessToken string) ([]githubEmail, error) {
	request, err := githubAPIRequest(ctx, githubAPIEmailsURL, accessToken)
	if err != nil {
		return nil, err
	}
	response, err := githubOAuthHTTPClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github emails response status: %d", response.StatusCode)
	}

	var emails []githubEmail
	if err := json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(&emails); err != nil {
		return nil, err
	}
	return emails, nil
}

func githubAPIRequest(ctx context.Context, endpoint string, accessToken string) (*http.Request, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("Authorization", "Bearer "+accessToken)
	request.Header.Set("User-Agent", "my-blog")
	return request, nil
}

func githubNickname(profile githubProfile) string {
	if value := strings.TrimSpace(profile.Name); value != "" {
		return value
	}
	if value := strings.TrimSpace(profile.Login); value != "" {
		return value
	}
	return fmt.Sprintf("GitHub %d", profile.ID)
}
