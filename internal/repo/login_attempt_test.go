package repo

import (
	"context"
	"fmt"
	"testing"
)

func TestLoginAttemptAccountLimitAndClear(t *testing.T) {
	server := setupEmailVerifyTest(t)
	ctx := context.Background()
	account := "LoginUser@example.com"
	userIP := "192.0.2.10"

	for attempt := 1; attempt <= loginAccountFailureLimit; attempt++ {
		blocked, err := RecordLoginFailure(ctx, account, userIP)
		if err != nil {
			t.Fatal(err)
		}
		if blocked != (attempt == loginAccountFailureLimit) {
			t.Fatalf("unexpected blocked state at attempt %d", attempt)
		}
	}

	blocked, err := IsLoginAttemptBlocked(ctx, " loginuser@EXAMPLE.com ", userIP)
	if err != nil {
		t.Fatal(err)
	}
	if !blocked {
		t.Fatal("expected normalized account to be blocked")
	}
	if ttl := server.TTL(loginAccountFailureKey(account)); ttl != loginFailureWindow {
		t.Fatalf("expected account failure TTL %s, got %s", loginFailureWindow, ttl)
	}

	if err := ClearLoginFailures(ctx, account, userIP); err != nil {
		t.Fatal(err)
	}
	blocked, err = IsLoginAttemptBlocked(ctx, account, userIP)
	if err != nil {
		t.Fatal(err)
	}
	if blocked {
		t.Fatal("expected login failures to be cleared")
	}
}

func TestLoginAttemptIPLimitAndWindow(t *testing.T) {
	server := setupEmailVerifyTest(t)
	ctx := context.Background()
	userIP := "192.0.2.20"

	for attempt := 1; attempt <= loginIPFailureLimit; attempt++ {
		blocked, err := RecordLoginFailure(
			ctx,
			fmt.Sprintf("login-user-%d", attempt),
			userIP,
		)
		if err != nil {
			t.Fatal(err)
		}
		if blocked != (attempt == loginIPFailureLimit) {
			t.Fatalf("unexpected blocked state at attempt %d", attempt)
		}
	}

	blocked, err := IsLoginAttemptBlocked(ctx, "another-user", userIP)
	if err != nil {
		t.Fatal(err)
	}
	if !blocked {
		t.Fatal("expected IP to be blocked")
	}

	server.FastForward(loginFailureWindow)
	blocked, err = IsLoginAttemptBlocked(ctx, "another-user", userIP)
	if err != nil {
		t.Fatal(err)
	}
	if blocked {
		t.Fatal("expected login failure window to expire")
	}
}
