package db

import (
	"1aides/pkg/log/zlog"
	"context"
	"fmt"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var (
	MongoDB *mongo.Database
	once    sync.Once
)

// NewMongoDB 初始化MongoDB客户端并设置全局MongoDB对象
func NewMongoDB() {
	once.Do(func() {
		// 从环境变量中获取MongoDB用户名、密码和地址
		mongoUser := getEnv("MONGO_USER", "tongque")
		mongoPassword := getEnv("MONGO_PASSWORD", "Y2hlbjA0MTY=")
		mongoHost := getEnv("MONGO_HOST", "localhost:27017")

		// 构建MongoDB连接URI
		uri := fmt.Sprintf("mongodb://%s:%s@%s", mongoUser, mongoPassword, mongoHost)

		// 设置MongoDB客户端选项
		clientOptions := options.Client().ApplyURI(uri)

		// 连接到MongoDB
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			zlog.Fatal("无法连接到MongoDB", zap.Error(err))
		}

		// 检查连接
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			zlog.Fatal("无法连接到MongoDB", zap.Error(err))
		}

		zlog.Info("成功连接到MongoDB", zap.String("URI", uri))

		// 初始化全局的MongoDB对象
		MongoDB = client.Database("aides")
	})
}

// GetMongoDB 返回全局的MongoDB实例
func GetMongoDB() *mongo.Database {
	if MongoDB == nil {
		NewMongoDB()
	}
	return MongoDB
}

// getEnv 获取环境变量值，如果未设置则返回默认值
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
