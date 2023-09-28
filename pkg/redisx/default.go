package redisx

import (
	"eGame-demo-back-office-api/configs"

	"github.com/go-redis/redis"
)

var redisClient *redis.Client

func Init() error {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     configs.App.Redis.Addr,
		Password: configs.App.Redis.Password,
		DB:       configs.App.Redis.Db,
	})

	err := redisClient.Ping().Err()
	if err != nil {
		return err
	}
	return nil
}

func InitLoginRdis() error {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     configs.Login.Addr,
		Password: configs.Login.Password,
		DB:       configs.Login.Db,
	})
	err := redisClient.Ping().Err()
	if err != nil {
		return err
	}
	return nil
}

// 获取redis客户端
func GetRedisClient() *redis.Client {
	return redisClient
}
