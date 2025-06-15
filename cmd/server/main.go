package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yoyo/api_bot/config"
	"github.com/yoyo/api_bot/internal/handlers"
	"github.com/yoyo/api_bot/internal/middleware"
	"github.com/yoyo/api_bot/internal/models"
	"github.com/yoyo/api_bot/internal/services"
)

func main() {
	// 初始化数据库
	db := config.SetupDatabase()

	// 初始化服务
	apiKeyService := services.NewAPIKeyService(db)

	// 初始化处理器
	apiKeyHandler := handlers.NewAPIKeyHandler(apiKeyService)

	// 设置路由
	r := gin.Default()

	// 添加初始API密钥创建路由（仅限本地使用）
	r.POST("/api/init-key", func(c *gin.Context) {
		// 检查是否已经存在API密钥
		var count int64
		db.Model(&models.APIKey{}).Count(&count)
		if count > 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "API密钥已经存在，无法创建初始密钥"})
			return
		}

		var req struct {
			Name          string `json:"name" binding:"required"`
			UserID        uint   `json:"user_id" binding:"required"`
			ExpiresInDays int    `json:"expires_in_days" binding:"required,min=1"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		apiKey, err := apiKeyService.CreateAPIKey(req.Name, req.UserID, req.ExpiresInDays)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, apiKey)
	})

	// API密钥管理路由组
	apiKeys := r.Group("/api/keys")
	{
		// 需要API密钥认证的路由
		apiKeys.Use(middleware.APIKeyAuth(apiKeyService))

		// 创建API密钥
		apiKeys.POST("/", apiKeyHandler.CreateAPIKey)

		// 获取所有API密钥
		apiKeys.GET("/", apiKeyHandler.GetAPIKeys)

		// 获取特定API密钥
		apiKeys.GET("/:id", apiKeyHandler.GetAPIKey)

		// 更新API密钥
		apiKeys.PUT("/:id", apiKeyHandler.UpdateAPIKey)

		// 删除API密钥
		apiKeys.DELETE("/:id", apiKeyHandler.DeleteAPIKey)
	}

	// 添加一个简单的健康检查端点
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("服务器启动在 :%s 端口", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
