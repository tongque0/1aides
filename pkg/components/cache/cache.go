package cache

import (
	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

func InitRedis(addr string, password string, db int) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     addr,     // Redis地址
		Password: password, // 密码
		DB:       db,       // 默认数据库
	})
}

func GetRedisClient() *redis.Client {
	return redisClient
}

// func GetModelConfig(userId string) (modhub.Model, error) {
// 	ctx := context.Background()
// 	key := "modelConfig:" + userId
// 	val, err := redisClient.Get(ctx, key).Result()
// 	if err == redis.Nil {
// 		// 键不存在
// 	} else if err != nil {
// 		return modhub.Model{}, err
// 	} else {
// 		// 假设model数据以JSON格式存储
// 		var model modhub.Model
// 		if err := json.Unmarshal([]byte(val), &model); err != nil {
// 			return modhub.Model{}, err
// 		}
// 		return model, nil
// 	}

// 	// 如果Redis中没有找到，从数据库加载，然后保存到Redis
// 	model, err := loadModelFromDB(userId) // 从DB加载
// 	if err != nil {
// 		return modhub.Model{}, err
// 	}

// 	data, err := json.Marshal(model)
// 	if err != nil {
// 		return modhub.Model{}, err
// 	}

// 	redisClient.Set(ctx, key, data, 30*time.Minute) // 缓存30分钟
// 	return model, nil
// }
