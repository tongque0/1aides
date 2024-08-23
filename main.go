package main

import (
	"1aides/pkg/components/db"
	"fmt"
)

func main() {
	// 获取全局的MongoDB实例
	mongoDB := db.GetMongoDB()

	// 使用MongoDB实例进行操作
	collection := mongoDB.Collection("testcollection")

	// 示例操作：插入文档
	doc := map[string]interface{}{"name": "Alice", "age": 25}
	insertResult, err := collection.InsertOne(nil, doc)
	if err != nil {
		fmt.Println("插入文档失败:", err)
		return
	}

	fmt.Println("插入的文档ID:", insertResult.InsertedID)
}
