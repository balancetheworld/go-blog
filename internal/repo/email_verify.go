package repo

import (
	"context"
	"crypto/subtle"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	emailVerifyTTL      = 10 * time.Minute
	emailVerifyCooldown = time.Minute
)

var saveEmailVerifyCodeScript = redis.NewScript(`
if redis.call("EXISTS", KEYS[2]) == 1 then
    return 0
end
redis.call("SET", KEYS[1], ARGV[1], "EX", ARGV[2])
redis.call("SET", KEYS[2], "1", "EX", ARGV[3])
return 1
`)

var consumeEmailVerifyCodeScript = redis.NewScript(`
local value = redis.call("GET", KEYS[1])
if not value then
    return 0
end
if value == ARGV[1] then
    redis.call("DEL", KEYS[1])
    return 1
end
return 0
`)

func SaveEmailVerifyCode(ctx context.Context, email, code string) (bool, error) {
	client, err := getRedis()
	if err != nil {
		return false, err
	}

	result, err := saveEmailVerifyCodeScript.Run(
		ctx,
		client,
		[]string{
			emailVerifyKey(email),
			emailVerifyCooldownKey(email),
		},
		code,
		int64(emailVerifyTTL/time.Second),
		int64(emailVerifyCooldown/time.Second),
	).Int64()
	if err != nil {
		return false, err
	}

	return result == 1, nil
}

func VerifyEmailCode(ctx context.Context, email, code string) (bool, error) {
	client, err := getRedis()
	if err != nil {
		return false, err
	}

	storedCode, err := client.Get(ctx, emailVerifyKey(email)).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return subtle.ConstantTimeCompare([]byte(storedCode), []byte(code)) == 1, nil
}

func ConsumeEmailVerifyCode(ctx context.Context, email, code string) (bool, error) {
	client, err := getRedis()
	if err != nil {
		return false, err
	}

	result, err := consumeEmailVerifyCodeScript.Run(
		ctx,
		client,
		[]string{emailVerifyKey(email)},
		code,
	).Int64()
	if err != nil {
		return false, err
	}

	return result == 1, nil
}

func emailVerifyKey(email string) string {
	return "email_verify:" + strings.ToLower(strings.TrimSpace(email))
}

func emailVerifyCooldownKey(email string) string {
	return "email_verify_cooldown:" + strings.ToLower(strings.TrimSpace(email))
}
