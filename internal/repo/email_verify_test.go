package repo

import (
	"context"
	"fmt"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

const testEmailPurpose = "register"

func setupEmailVerifyTest(t *testing.T) *miniredis.Miniredis {
	t.Helper()

	server := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{
		Addr: server.Addr(),
	})
	previousClient := redisClient
	redisClient = client
	t.Cleanup(func() {
		_ = client.Close()
		redisClient = previousClient
	})

	return server
}

func TestEmailVerifyLocksAfterFifthInvalidAttempt(t *testing.T) {
	server := setupEmailVerifyTest(t)
	ctx := context.Background()
	email := "email-verify-lock@example.com"

	saved, err := SaveEmailVerifyCode(ctx, email, testEmailPurpose, "123456", "lock-test-ip")
	if err != nil {
		t.Fatal(err)
	}
	if !saved {
		t.Fatal("expected verification code to be saved")
	}

	for attempt := 1; attempt <= emailVerifyMaxTries; attempt++ {
		valid, err := ReserveEmailVerifyCode(
			ctx,
			email,
			testEmailPurpose,
			"654321",
			fmt.Sprintf("invalid-token-%d", attempt),
		)
		if err != nil {
			t.Fatal(err)
		}
		if valid {
			t.Fatalf("expected attempt %d to fail", attempt)
		}
		if attempt < emailVerifyMaxTries && server.Exists(emailVerifyLockKey(testEmailPurpose, email)) {
			t.Fatalf("expected no lock before attempt %d", emailVerifyMaxTries)
		}
	}

	if !server.Exists(emailVerifyLockKey(testEmailPurpose, email)) {
		t.Fatal("expected email verification to be locked")
	}
	if ttl := server.TTL(emailVerifyLockKey(testEmailPurpose, email)); ttl != emailVerifyLockTTL {
		t.Fatalf("expected lock TTL %s, got %s", emailVerifyLockTTL, ttl)
	}
	if server.Exists(emailVerifyKey(testEmailPurpose, email)) {
		t.Fatal("expected verification code to be deleted after lock")
	}
	if server.Exists(emailVerifyAttemptKey(testEmailPurpose, email)) {
		t.Fatal("expected attempt counter to be deleted after lock")
	}

	valid, err := ReserveEmailVerifyCode(ctx, email, testEmailPurpose, "123456", "locked-token")
	if err != nil {
		t.Fatal(err)
	}
	if valid {
		t.Fatal("expected correct code to fail while locked")
	}
}

func TestEmailVerifyCodeIsConsumedOnce(t *testing.T) {
	setupEmailVerifyTest(t)
	ctx := context.Background()
	email := "email-verify-consume@example.com"

	saved, err := SaveEmailVerifyCode(ctx, email, testEmailPurpose, "123456", "consume-test-ip")
	if err != nil {
		t.Fatal(err)
	}
	if !saved {
		t.Fatal("expected verification code to be saved")
	}

	reserved, err := ReserveEmailVerifyCode(ctx, email, testEmailPurpose, "123456", "consume-token")
	if err != nil {
		t.Fatal(err)
	}
	if !reserved {
		t.Fatal("expected verification code to be reserved")
	}

	committed, err := CommitEmailVerifyCode(ctx, email, testEmailPurpose, "consume-token")
	if err != nil {
		t.Fatal(err)
	}
	if !committed {
		t.Fatal("expected verification code to be committed")
	}

	reserved, err = ReserveEmailVerifyCode(ctx, email, testEmailPurpose, "123456", "replay-token")
	if err != nil {
		t.Fatal(err)
	}
	if reserved {
		t.Fatal("expected committed verification code to be unavailable")
	}

	committed, err = CommitEmailVerifyCode(ctx, email, testEmailPurpose, "consume-token")
	if err != nil {
		t.Fatal(err)
	}
	if committed {
		t.Fatal("expected second commit to fail")
	}
}

func TestEmailVerifyRejectsEleventhIPRequest(t *testing.T) {
	server := setupEmailVerifyTest(t)
	ctx := context.Background()
	userIP := "ip-limit-test"

	for request := 1; request <= emailVerifyIPLimit; request++ {
		email := fmt.Sprintf("email-verify-ip-%d@example.com", request)
		saved, err := SaveEmailVerifyCode(ctx, email, testEmailPurpose, "123456", userIP)
		if err != nil {
			t.Fatal(err)
		}
		if !saved {
			t.Fatalf("expected request %d to be accepted", request)
		}
	}

	blockedEmail := "email-verify-ip-blocked@example.com"
	saved, err := SaveEmailVerifyCode(ctx, blockedEmail, testEmailPurpose, "123456", userIP)
	if err != nil {
		t.Fatal(err)
	}
	if saved {
		t.Fatal("expected eleventh IP request to be rejected")
	}
	if server.Exists(emailVerifyKey(testEmailPurpose, blockedEmail)) {
		t.Fatal("expected rejected request not to save a code")
	}
	if ttl := server.TTL(emailVerifyIPKey(userIP)); ttl != emailVerifyIPWindow {
		t.Fatalf("expected IP limit TTL %s, got %s", emailVerifyIPWindow, ttl)
	}

	server.FastForward(emailVerifyIPWindow)
	saved, err = SaveEmailVerifyCode(
		ctx,
		"email-verify-ip-window-reset@example.com",
		testEmailPurpose,
		"123456",
		userIP,
	)
	if err != nil {
		t.Fatal(err)
	}
	if !saved {
		t.Fatal("expected request after IP window to be accepted")
	}
}

func TestEmailVerifyCooldown(t *testing.T) {
	server := setupEmailVerifyTest(t)
	ctx := context.Background()
	email := "email-verify-cooldown@example.com"
	userIP := "cooldown-test-ip"

	saved, err := SaveEmailVerifyCode(ctx, email, testEmailPurpose, "123456", userIP)
	if err != nil {
		t.Fatal(err)
	}
	if !saved {
		t.Fatal("expected first request to be accepted")
	}
	if ttl := server.TTL(emailVerifyCooldownKey(testEmailPurpose, email)); ttl != emailVerifyCooldown {
		t.Fatalf("expected cooldown TTL %s, got %s", emailVerifyCooldown, ttl)
	}

	saved, err = SaveEmailVerifyCode(ctx, email, testEmailPurpose, "654321", userIP)
	if err != nil {
		t.Fatal(err)
	}
	if saved {
		t.Fatal("expected request during cooldown to be rejected")
	}
	code, err := server.Get(emailVerifyKey(testEmailPurpose, email))
	if err != nil {
		t.Fatal(err)
	}
	if code != "123456" {
		t.Fatalf("expected original code to remain, got %s", code)
	}

	server.FastForward(emailVerifyCooldown)
	saved, err = SaveEmailVerifyCode(ctx, email, testEmailPurpose, "654321", userIP)
	if err != nil {
		t.Fatal(err)
	}
	if !saved {
		t.Fatal("expected request after cooldown to be accepted")
	}
	code, err = server.Get(emailVerifyKey(testEmailPurpose, email))
	if err != nil {
		t.Fatal(err)
	}
	if code != "654321" {
		t.Fatalf("expected new code to be saved, got %s", code)
	}
}

func TestEmailVerifyReservationLifecycle(t *testing.T) {
	server := setupEmailVerifyTest(t)
	ctx := context.Background()
	email := "email-verify-reservation@example.com"
	userIP := "reservation-test-ip"

	saved, err := SaveEmailVerifyCode(ctx, email, testEmailPurpose, "123456", userIP)
	if err != nil {
		t.Fatal(err)
	}
	if !saved {
		t.Fatal("expected verification code to be saved")
	}

	reserved, err := ReserveEmailVerifyCode(ctx, email, testEmailPurpose, "123456", "token-1")
	if err != nil {
		t.Fatal(err)
	}
	if !reserved {
		t.Fatal("expected verification code to be reserved")
	}
	if codeTTL, reservationTTL := server.TTL(emailVerifyKey(testEmailPurpose, email)), server.TTL(emailVerifyReservationKey(testEmailPurpose, email)); codeTTL != reservationTTL {
		t.Fatalf("expected matching TTLs, got code=%s reservation=%s", codeTTL, reservationTTL)
	}

	reserved, err = ReserveEmailVerifyCode(ctx, email, testEmailPurpose, "123456", "token-2")
	if err != nil {
		t.Fatal(err)
	}
	if reserved {
		t.Fatal("expected concurrent reservation to fail")
	}

	server.FastForward(emailVerifyCooldown)
	saved, err = SaveEmailVerifyCode(ctx, email, testEmailPurpose, "654321", userIP)
	if err != nil {
		t.Fatal(err)
	}
	if saved {
		t.Fatal("expected save while reserved to fail")
	}

	released, err := ReleaseEmailVerifyCode(ctx, email, testEmailPurpose, "wrong-token")
	if err != nil {
		t.Fatal(err)
	}
	if released {
		t.Fatal("expected release with wrong token to fail")
	}

	released, err = ReleaseEmailVerifyCode(ctx, email, testEmailPurpose, "token-1")
	if err != nil {
		t.Fatal(err)
	}
	if !released {
		t.Fatal("expected reservation to be released")
	}

	reserved, err = ReserveEmailVerifyCode(ctx, email, testEmailPurpose, "123456", "token-2")
	if err != nil {
		t.Fatal(err)
	}
	if !reserved {
		t.Fatal("expected code to be reusable after release")
	}
}
