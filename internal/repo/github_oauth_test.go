package repo

import (
	"context"
	"testing"
)

func TestGithubOAuthStateIsConsumedOnce(t *testing.T) {
	server := setupEmailVerifyTest(t)
	ctx := context.Background()
	state := "github-state-test"

	saved, err := SaveGithubOAuthState(ctx, state)
	if err != nil {
		t.Fatal(err)
	}
	if !saved {
		t.Fatal("expected oauth state to be saved")
	}

	saved, err = SaveGithubOAuthState(ctx, state)
	if err != nil {
		t.Fatal(err)
	}
	if saved {
		t.Fatal("expected duplicate oauth state to be rejected")
	}

	valid, err := ConsumeGithubOAuthState(ctx, state)
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Fatal("expected oauth state to be valid")
	}
	if server.Exists(githubOAuthStateKey(state)) {
		t.Fatal("expected oauth state to be deleted")
	}

	valid, err = ConsumeGithubOAuthState(ctx, state)
	if err != nil {
		t.Fatal(err)
	}
	if valid {
		t.Fatal("expected oauth state to be single-use")
	}
}
