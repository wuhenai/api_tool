package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yoyo/api_bot/internal/models"
	"github.com/yoyo/api_bot/internal/services"
)

// APIKeyHandler 处理API密钥相关的HTTP请求
type APIKeyHandler struct {
	apiKeyService *services.APIKeyService
}

// NewAPIKeyHandler 创建新的API密钥处理器
func NewAPIKeyHandler(apiKeyService *services.APIKeyService) *APIKeyHandler {
	return &APIKeyHandler{
		apiKeyService: apiKeyService,
	}
}

// CreateAPIKey 创建新的API密钥
func (h *APIKeyHandler) CreateAPIKey(c *gin.Context) {
	var req struct {
		Name          string `json:"name" binding:"required"`
		ExpiresInDays int    `json:"expires_in_days" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	apiKey, err := h.apiKeyService.CreateAPIKey(req.Name, userID.(uint), req.ExpiresInDays)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, apiKey)
}

// GetAPIKeys 获取当前用户的所有API密钥
func (h *APIKeyHandler) GetAPIKeys(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	apiKeys, err := h.apiKeyService.GetAPIKeysByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, apiKeys)
}

// GetAPIKey 获取特定的API密钥
func (h *APIKeyHandler) GetAPIKey(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取API密钥
	apiKeys, err := h.apiKeyService.GetAPIKeysByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 查找特定ID的API密钥
	var targetAPIKey *models.APIKey
	for _, key := range apiKeys {
		if key.ID == uint(id) {
			targetAPIKey = &key
			break
		}
	}

	if targetAPIKey == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "API密钥不存在或不属于当前用户"})
		return
	}

	c.JSON(http.StatusOK, targetAPIKey)
}

// UpdateAPIKey 更新API密钥
func (h *APIKeyHandler) UpdateAPIKey(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var req struct {
		Name          string `json:"name"`
		Active        bool   `json:"active"`
		ExpiresInDays int    `json:"expires_in_days"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取API密钥
	apiKeys, err := h.apiKeyService.GetAPIKeysByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 检查API密钥是否属于当前用户
	var found bool
	for _, key := range apiKeys {
		if key.ID == uint(id) {
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "API密钥不存在或不属于当前用户"})
		return
	}

	// 更新API密钥
	apiKey, err := h.apiKeyService.UpdateAPIKey(uint(id), req.Name, req.Active, req.ExpiresInDays)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, apiKey)
}

// DeleteAPIKey 删除API密钥
func (h *APIKeyHandler) DeleteAPIKey(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取API密钥
	apiKeys, err := h.apiKeyService.GetAPIKeysByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 检查API密钥是否属于当前用户
	var found bool
	for _, key := range apiKeys {
		if key.ID == uint(id) {
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "API密钥不存在或不属于当前用户"})
		return
	}

	// 删除API密钥
	if err := h.apiKeyService.DeleteAPIKey(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "API密钥已成功删除"})
}
