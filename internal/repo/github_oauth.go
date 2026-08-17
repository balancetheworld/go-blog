package repo

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const githubOAuthStateTTL = 10 * time.Minute

func SaveGithubOAuthState(ctx context.Context, state string) (bool, error) {
	client, err := getRedis()
	if err != nil {
		return false, err
	}

	return client.SetNX(
		ctx,
		githubOAuthStateKey(state),
		"1",
		githubOAuthStateTTL,
	).Result()
}

func ConsumeGithubOAuthState(ctx context.Context, state string) (bool, error) {
	client, err := getRedis()
	if err != nil {
		return false, err
	}

	value, err := client.GetDel(ctx, githubOAuthStateKey(state)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, err
	}

	return value == "1", nil
}

func githubOAuthStateKey(state string) string {
	return "github_oauth_state:" + strings.TrimSpace(state)
}
