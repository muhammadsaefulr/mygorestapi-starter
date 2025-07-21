package database

import (
	"context"

	"github.com/muhammadsaefulr/NimeStreamAPI/config"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/utils"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func ConnectRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisHost,
		Password: config.RedisPassword,
		DB:       0,
	})

	_, err := client.Ping(Ctx).Result()
	if err != nil {
		utils.Log.Errorf("Failed to connect to redis: %+v", err)
	}

	return client
}
