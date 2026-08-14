package repo

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	loginAccountFailureLimit = 5
	loginIPFailureLimit      = 20
	loginFailureWindow       = 15 * time.Minute
)

var recordLoginFailureScript = redis.NewScript(`
local accountAttempts = redis.call("INCR", KEYS[1])
if accountAttempts == 1 then
    redis.call("EXPIRE", KEYS[1], ARGV[1])
end
local ipAttempts = redis.call("INCR", KEYS[2])
if ipAttempts == 1 then
    redis.call("EXPIRE", KEYS[2], ARGV[1])
end
if accountAttempts >= tonumber(ARGV[2]) or ipAttempts >= tonumber(ARGV[3]) then
    return 1
end
return 0
`)

func IsLoginAttemptBlocked(ctx context.Context, account, userIP string) (bool, error) {
	client, err := getRedis()
	if err != nil {
		return false, err
	}

	values, err := client.MGet(
		ctx,
		loginAccountFailureKey(account),
		loginIPFailureKey(userIP),
	).Result()
	if err != nil {
		return false, err
	}

	return loginFailureCount(values[0]) >= loginAccountFailureLimit ||
		loginFailureCount(values[1]) >= loginIPFailureLimit, nil
}

func RecordLoginFailure(ctx context.Context, account, userIP string) (bool, error) {
	client, err := getRedis()
	if err != nil {
		return false, err
	}

	result, err := recordLoginFailureScript.Run(
		ctx,
		client,
		[]string{
			loginAccountFailureKey(account),
			loginIPFailureKey(userIP),
		},
		int64(loginFailureWindow/time.Second),
		loginAccountFailureLimit,
		loginIPFailureLimit,
	).Int64()
	if err != nil {
		return false, err
	}

	return result == 1, nil
}

func ClearLoginFailures(ctx context.Context, account, userIP string) error {
	client, err := getRedis()
	if err != nil {
		return err
	}

	return client.Del(
		ctx,
		loginAccountFailureKey(account),
		loginIPFailureKey(userIP),
	).Err()
}

func loginFailureCount(value any) int {
	text, ok := value.(string)
	if !ok {
		return 0
	}

	count := 0
	for _, digit := range text {
		if digit < '0' || digit > '9' {
			return 0
		}
		count = count*10 + int(digit-'0')
	}
	return count
}

func loginAccountFailureKey(account string) string {
	return "login_failure:account:" + loginFailureKeyPart(account)
}

func loginIPFailureKey(userIP string) string {
	return "login_failure:ip:" + loginFailureKeyPart(userIP)
}

func loginFailureKeyPart(value string) string {
	digest := sha256.Sum256([]byte(strings.ToLower(strings.TrimSpace(value))))
	return hex.EncodeToString(digest[:])
}
