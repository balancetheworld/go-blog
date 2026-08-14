package repo

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCloseDatabaseIsIdempotent(t *testing.T) {
	database, err := gorm.Open(
		sqlite.Open("file:close_database_test?mode=memory&cache=shared"),
		&gorm.Config{},
	)
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, err := database.DB()
	if err != nil {
		t.Fatal(err)
	}

	previousDB := db
	db = database
	t.Cleanup(func() {
		db = previousDB
	})

	if err := CloseDatabase(); err != nil {
		t.Fatal(err)
	}
	if db != nil {
		t.Fatal("expected database reference to be cleared")
	}
	if err := sqlDB.Ping(); err == nil {
		t.Fatal("expected database connection to be closed")
	}
	if err := CloseDatabase(); err != nil {
		t.Fatal(err)
	}
}

func TestCloseRedisIsIdempotent(t *testing.T) {
	server := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{
		Addr: server.Addr(),
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		t.Fatal(err)
	}

	previousClient := redisClient
	redisClient = client
	t.Cleanup(func() {
		redisClient = previousClient
	})

	if err := CloseRedis(); err != nil {
		t.Fatal(err)
	}
	if redisClient != nil {
		t.Fatal("expected Redis reference to be cleared")
	}
	if err := client.Ping(context.Background()).Err(); err == nil {
		t.Fatal("expected Redis client to be closed")
	}
	if err := CloseRedis(); err != nil {
		t.Fatal(err)
	}
}
