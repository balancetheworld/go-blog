package repo

import (
	"context"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	emailVerifyTTL      = 10 * time.Minute
	emailVerifyCooldown = time.Minute
	emailVerifyLockTTL  = 15 * time.Minute
	emailVerifyIPWindow = 10 * time.Minute
	emailVerifyMaxTries = 5
	emailVerifyIPLimit  = 10
)

var saveEmailVerifyCodeScript = redis.NewScript(`
if redis.call("EXISTS", KEYS[4]) == 1 then
    return 0
end
if redis.call("EXISTS", KEYS[2]) == 1 then
    return 0
end
local requests = redis.call("INCR", KEYS[3])
if requests == 1 then
    redis.call("EXPIRE", KEYS[3], ARGV[4])
end
if requests > tonumber(ARGV[5]) then
    return -1
end
redis.call("SET", KEYS[1], ARGV[1], "EX", ARGV[2])
redis.call("SET", KEYS[2], "1", "EX", ARGV[3])
return 1
`)

var reserveEmailVerifyCodeScript = redis.NewScript(`
if redis.call("EXISTS", KEYS[3]) == 1 then
    return -1
end
local value = redis.call("GET", KEYS[1])
if not value then
    return 0
end
if value == ARGV[1] then
    local ttl = redis.call("PTTL", KEYS[1])
    if ttl <= 0 then
        return 0
    end
    local reserved = redis.call("SET", KEYS[4], ARGV[2], "PX", ttl, "NX")
    if not reserved then
        return -1
    end
    return 1
end
local attempts = redis.call("INCR", KEYS[2])
if attempts == 1 then
    redis.call("EXPIRE", KEYS[2], ARGV[3])
end
if attempts >= tonumber(ARGV[4]) then
    redis.call("DEL", KEYS[1], KEYS[2])
    redis.call("SET", KEYS[3], "1", "EX", ARGV[5])
    return -1
end
return 0
`)

var commitEmailVerifyCodeScript = redis.NewScript(`
if redis.call("GET", KEYS[4]) ~= ARGV[1] then
    return 0
end
redis.call("DEL", KEYS[1], KEYS[2], KEYS[4])
return 1
`)

var releaseEmailVerifyCodeScript = redis.NewScript(`
if redis.call("GET", KEYS[4]) ~= ARGV[1] then
    return 0
end
redis.call("DEL", KEYS[4])
return 1
`)

func SaveEmailVerifyCode(ctx context.Context, email, purpose, code, userIP string) (bool, error) {
	client, err := getRedis()
	if err != nil {
		return false, err
	}

	result, err := saveEmailVerifyCodeScript.Run(
		ctx,
		client,
		[]string{
			emailVerifyKey(purpose, email),
			emailVerifyCooldownKey(purpose, email),
			emailVerifyIPKey(userIP),
			emailVerifyReservationKey(purpose, email),
		},
		code,
		int64(emailVerifyTTL/time.Second),
		int64(emailVerifyCooldown/time.Second),
		int64(emailVerifyIPWindow/time.Second),
		emailVerifyIPLimit,
	).Int64()
	if err != nil {
		return false, err
	}

	return result == 1, nil
}

func ReserveEmailVerifyCode(ctx context.Context, email, purpose, code, token string) (bool, error) {
	client, err := getRedis()
	if err != nil {
		return false, err
	}

	result, err := reserveEmailVerifyCodeScript.Run(
		ctx,
		client,
		[]string{
			emailVerifyKey(purpose, email),
			emailVerifyAttemptKey(purpose, email),
			emailVerifyLockKey(purpose, email),
			emailVerifyReservationKey(purpose, email),
		},
		code,
		token,
		int64(emailVerifyTTL/time.Second),
		emailVerifyMaxTries,
		int64(emailVerifyLockTTL/time.Second),
	).Int64()
	if err != nil {
		return false, err
	}

	return result == 1, nil
}

func CommitEmailVerifyCode(ctx context.Context, email, purpose, token string) (bool, error) {
	client, err := getRedis()
	if err != nil {
		return false, err
	}

	result, err := commitEmailVerifyCodeScript.Run(
		ctx,
		client,
		[]string{
			emailVerifyKey(purpose, email),
			emailVerifyAttemptKey(purpose, email),
			emailVerifyLockKey(purpose, email),
			emailVerifyReservationKey(purpose, email),
		},
		token,
	).Int64()
	if err != nil {
		return false, err
	}

	return result == 1, nil
}

func ReleaseEmailVerifyCode(ctx context.Context, email, purpose, token string) (bool, error) {
	client, err := getRedis()
	if err != nil {
		return false, err
	}

	result, err := releaseEmailVerifyCodeScript.Run(
		ctx,
		client,
		[]string{
			emailVerifyKey(purpose, email),
			emailVerifyAttemptKey(purpose, email),
			emailVerifyLockKey(purpose, email),
			emailVerifyReservationKey(purpose, email),
		},
		token,
	).Int64()
	if err != nil {
		return false, err
	}

	return result == 1, nil
}

func emailVerifyKey(purpose, email string) string {
	return "email_verify:" + strings.ToLower(strings.TrimSpace(purpose)) + ":" + strings.ToLower(strings.TrimSpace(email))
}

func emailVerifyCooldownKey(purpose, email string) string {
	return "email_verify_cooldown:" + strings.ToLower(strings.TrimSpace(purpose)) + ":" + strings.ToLower(strings.TrimSpace(email))
}

func emailVerifyAttemptKey(purpose, email string) string {
	return "email_verify_attempt:" + strings.ToLower(strings.TrimSpace(purpose)) + ":" + strings.ToLower(strings.TrimSpace(email))
}

func emailVerifyLockKey(purpose, email string) string {
	return "email_verify_lock:" + strings.ToLower(strings.TrimSpace(purpose)) + ":" + strings.ToLower(strings.TrimSpace(email))
}

func emailVerifyReservationKey(purpose, email string) string {
	return "email_verify_reservation:" + strings.ToLower(strings.TrimSpace(purpose)) + ":" + strings.ToLower(strings.TrimSpace(email))
}

func emailVerifyIPKey(userIP string) string {
	return "email_verify_ip:" + strings.TrimSpace(userIP)
}
