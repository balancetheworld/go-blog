package utils

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/zyj/my-blog/pkg/constant"
)

var captchaClient = &http.Client{
	Timeout: 10 * time.Second,
}

func VerifyCaptcha(
	ctx context.Context,
	token string,
	remoteIP string,
) (bool, error) {
	provider := constant.CaptchaType(
		Get(
			constant.EnvKeyCaptchaProvider,
			string(constant.CaptchaDisable),
		),
	)
	if provider == constant.CaptchaDisable {
		return true, nil
	}
	if token == "" {
		return false, nil
	}

	endpoint, err := captchaEndpoint(provider)
	if err != nil {
		return false, err
	}
	values := url.Values{
		"secret":   {Get(constant.EnvKeyCaptchaSecretKey)},
		"response": {token},
	}
	if remoteIP != "" {
		values.Set("remoteip", remoteIP)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		endpoint,
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := captchaClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return false, errors.New("captcha provider request failed")
	}

	var result struct {
		Success bool `json:"success"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}

	return result.Success, nil
}

func captchaEndpoint(provider constant.CaptchaType) (string, error) {
	if endpoint := Get(constant.EnvKeyCaptchaURL); endpoint != "" {
		return endpoint, nil
	}

	switch provider {
	case constant.CaptchaTurnstile:
		return "https://challenges.cloudflare.com/turnstile/v0/siteverify", nil
	case constant.CaptchaRecaptcha:
		return "https://www.google.com/recaptcha/api/siteverify", nil
	case constant.CaptchaHcaptcha:
		return "https://api.hcaptcha.com/siteverify", nil
	default:
		return "", errors.New("unsupported captcha provider")
	}
}
