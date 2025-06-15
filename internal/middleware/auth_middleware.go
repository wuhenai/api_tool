package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yoyo/api_bot/internal/services"
)

// APIKeyAuth 中间件验证API密钥
func APIKeyAuth(apiKeyService *services.APIKeyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var apiKey string

		// 首先从查询参数中获取API密钥
		apiKey = c.Query("key")

		// 如果查询参数中没有，尝试从请求头获取
		if apiKey == "" {
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" {
				// 检查授权格式
				parts := strings.SplitN(authHeader, " ", 2)
				if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
					apiKey = parts[1]
				}
			}
		}

		// 如果仍然没有API密钥，返回错误
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少API密钥，请通过查询参数'key'或Authorization头提供"})
			c.Abort()
			return
		}

		// 验证API密钥
		key, err := apiKeyService.ValidateAPIKey(apiKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// 将用户ID存储在上下文中，以便后续处理
		c.Set("userID", key.UserID)
		c.Set("apiKey", key)

		c.Next()
	}
}
