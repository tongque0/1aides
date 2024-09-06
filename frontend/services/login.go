package services

import (
	"1aides/pkg/components/db"
	"1aides/pkg/log/zlog"
	"context"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

// RegisterLoginRoutes 注册登录相关的路由
func RegisterLoginRoutes(router *gin.Engine) {
	// 登录页面
	router.GET("/login", LoginHandler)
	// 登录表单提交
	router.POST("/login/login", Login)
}

// LoginHandler 处理登陆请求
func LoginHandler(c *gin.Context) {
	c.HTML(200, "login.tmpl", gin.H{
		"ActivePage": "login", // 设置活动页面
	})
}

// Login 处理登录请求
func Login(c *gin.Context) {
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// 从请求体中解析 JSON 数据
	if err := c.BindJSON(&loginData); err != nil {
		zlog.Error("无法解析登录请求的JSON数据", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	username := loginData.Username
	password := loginData.Password
	collection := db.GetMongoDB().Collection("webadmin")
	zlog.Info("登录请求", zap.String("username", username))

	// 查询数据库中的用户信息
	var result bson.M
	err := collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&result)
	if err != nil {
		zlog.Error("无法获取用户信息", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// 验证密码是否匹配
	if password != result["password"] {
		zlog.Warn("密码验证失败", zap.String("username", username))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// 登录成功，将用户信息存储到Session中
	session := sessions.Default(c)
	session.Set("username", username) // 存储用户名
	session.Save()

	// 登录成功后返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}
