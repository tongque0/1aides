package db

import (
	"1aides/pkg/log/zlog"
	"context"
	"fmt"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type Config struct {
	Model   string `bson:"model"`
	ApiKey  string `bson:"apiKey"`
	BaseURL string `bson:"baseURL"`
	Prompt  string `bson:"prompt"`
}

type GPTData struct {
	Type   string `bson:"type"`
	Config Config `bson:"config"`
}

var (
	MongoDB *mongo.Database
	once    sync.Once
)

// NewMongoDB 初始化MongoDB客户端并设置全局MongoDB对象
func NewMongoDB() {
	once.Do(func() {
		// 从环境变量中获取MongoDB用户名、密码和地址
		mongoUser := getEnv("MONGO_USER", "aides")
		mongoPassword := getEnv("MONGO_PASSWORD", "dGhpcyBpcyBhaWRlcw==")
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

		// 检查并初始化数据表和数据
		ensureData(MongoDB)
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

// ensureData 检查并初始化数据表和数据
func ensureData(db *mongo.Database) {
	collectionName := "models"
	collection := db.Collection(collectionName)

	// 检查集合中的文档数量
	count, err := collection.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		zlog.Fatal("无法获取数据表信息", zap.Error(err))
	}

	// 如果集合为空，初始化数据
	if count < 1 {
		initialData := GPTData{
			Type: "GPT",
			Config: Config{
				Model:   "gpt-4o-mini",
				ApiKey:  "sk-29sMaKDD5aBgDtyx02014694972846Cc8c8b9fEb18192532",
				BaseURL: "https://prime.zetatechs.com/v1",
				Prompt:  "你的身份是一位微信消息机器人，你的开发者是同阙。你可以回复任何你想回复的内容，但是要有逻辑。",
			},
		}
		_, err := collection.InsertOne(context.TODO(), initialData)
		if err != nil {
			zlog.Fatal("无法初始化数据", zap.Error(err))
		}
		zlog.Info("数据表已初始化", zap.String("collection", collectionName))
	}
}
