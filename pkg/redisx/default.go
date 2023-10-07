package redisx

import (
	"eGame-demo-back-office-api/configs"
	"encoding/json"

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

// func InitLoginRdis() error {
// 	redisClient = redis.NewClient(&redis.Options{
// 		Addr:     configs.Login.Addr,
// 		Password: configs.Login.Password,
// 		DB:       configs.Login.Db,
// 	})
// 	err := redisClient.Ping().Err()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// 获取redis客户端
func GetRedisClient() *redis.Client {
	return redisClient
}

type UserData struct {
	Groupname   string          `json:"groupname"`
	Permissions map[string]bool `json:"permissions"`
	Token       string          `json:"token"`
	Uid         int             `json:"uid"`
	Username    string          `json:"username"`
}

func GetUserDataFromRedis(token string) (UserData, error) {
	// 使用提供的令牌從Redis中檢索用戶數據
	redisClient := GetRedisClient()
	data, err := redisClient.Get(token).Result()
	if err != nil {
		// 处理从Redis中检索数据时出错的情况
		return UserData{}, err
	}

	// 解析从Redis检索的JSON数据
	var userData UserData
	if err := json.Unmarshal([]byte(data), &userData); err != nil {
		// 处理JSON解析错误的情况
		return UserData{}, err
	}

	return userData, nil
}
