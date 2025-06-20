package database

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var CtxDb = context.Background()

func CreateClientRedis(dbNo int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("DB_ADDR"),
		Password: os.Getenv("DB_PASS"),
		DB:       dbNo,
	})
	return rdb
}
