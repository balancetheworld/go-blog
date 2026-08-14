package utils

import (
	"testing"

	"github.com/zyj/my-blog/pkg/constant"
)

func TestGenerateAndParseTokenPair(t *testing.T) {
	t.Setenv(constant.EnvKeyJWTSecret, "0123456789abcdef0123456789abcdef")
	t.Setenv(constant.EnvKeyTokenDuration, "3600")
	t.Setenv(constant.EnvKeyRefreshTokenDuration, "604800")

	pair, err := GenerateTokenPair(1, constant.RoleUser, "session-id", false)
	if err != nil {
		t.Fatal(err)
	}

	accessClaims, err := ParseAccessToken(pair.AccessToken)
	if err != nil {
		t.Fatal(err)
	}
	if accessClaims.UserID != 1 || accessClaims.SessionID != "session-id" {
		t.Fatalf("unexpected access claims: %#v", accessClaims)
	}

	refreshClaims, err := ParseRefreshToken(pair.RefreshToken)
	if err != nil {
		t.Fatal(err)
	}
	if refreshClaims.TokenType != TokenTypeRefresh {
		t.Fatalf("unexpected refresh token type: %s", refreshClaims.TokenType)
	}

	if _, err := ParseRefreshToken(pair.AccessToken); err == nil {
		t.Fatal("access token must not be accepted as refresh token")
	}
}

func TestGenerateEmailVerifyCode(t *testing.T) {
	code, err := GenerateEmailVerifyCode()
	if err != nil {
		t.Fatal(err)
	}
	if len(code) != 6 {
		t.Fatalf("unexpected verification code length: %d", len(code))
	}
	for _, value := range code {
		if value < '0' || value > '9' {
			t.Fatalf("verification code contains non-digit: %q", value)
		}
	}
}

func TestDisabledCaptcha(t *testing.T) {
	t.Setenv(constant.EnvKeyMode, string(constant.ModeDev))
	t.Setenv(constant.EnvKeyCaptchaProvider, string(constant.CaptchaDisable))

	valid, err := VerifyCaptcha(t.Context(), "", "")
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Fatal("disabled captcha must allow the request")
	}
}

func TestValidateTokenConfig(t *testing.T) {
	t.Setenv(constant.EnvKeyJWTSecret, "0123456789abcdef0123456789abcdef")
	t.Setenv(constant.EnvKeyTokenDuration, "3600")
	t.Setenv(constant.EnvKeyRefreshTokenDuration, "604800")
	t.Setenv(constant.EnvKeyRefreshTokenDurationWithRemember, "2592000")

	if err := ValidateTokenConfig(); err != nil {
		t.Fatal(err)
	}

	t.Setenv(constant.EnvKeyTokenDuration, "0")
	if err := ValidateTokenConfig(); err == nil {
		t.Fatal("expected zero token duration to be rejected")
	}
}

func TestValidateCaptchaConfig(t *testing.T) {
	t.Setenv(constant.EnvKeyMode, string(constant.ModeProd))
	t.Setenv(constant.EnvKeyCaptchaProvider, string(constant.CaptchaDisable))
	if err := ValidateCaptchaConfig(); err == nil {
		t.Fatal("expected disabled production captcha to be rejected")
	}

	t.Setenv(constant.EnvKeyCaptchaProvider, string(constant.CaptchaTurnstile))
	t.Setenv(constant.EnvKeyCaptchaSiteKey, "site-key")
	t.Setenv(constant.EnvKeyCaptchaSecretKey, "secret-key")
	if err := ValidateCaptchaConfig(); err != nil {
		t.Fatal(err)
	}

	t.Setenv(constant.EnvKeyCaptchaSecretKey, "")
	if err := ValidateCaptchaConfig(); err == nil {
		t.Fatal("expected missing captcha secret to be rejected")
	}
}
