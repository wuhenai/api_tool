package services

import (
	"errors"
	"time"

	"github.com/yoyo/api_bot/internal/models"
	"gorm.io/gorm"
)

// APIKeyService 处理API密钥相关的业务逻辑
type APIKeyService struct {
	db *gorm.DB
}

// NewAPIKeyService 创建新的API密钥服务实例
func NewAPIKeyService(db *gorm.DB) *APIKeyService {
	return &APIKeyService{db: db}
}

// CreateAPIKey 创建新的API密钥
func (s *APIKeyService) CreateAPIKey(name string, userID uint, expiresInDays int) (*models.APIKey, error) {
	if name == "" {
		return nil, errors.New("密钥名称不能为空")
	}

	expiresAt := time.Now().AddDate(0, 0, expiresInDays)

	apiKey := &models.APIKey{
		Name:      name,
		UserID:    userID,
		ExpiresAt: expiresAt,
		Active:    true,
	}

	result := s.db.Create(apiKey)
	if result.Error != nil {
		return nil, result.Error
	}

	return apiKey, nil
}

// GetAPIKeyByKey 通过密钥字符串获取API密钥
func (s *APIKeyService) GetAPIKeyByKey(key string) (*models.APIKey, error) {
	var apiKey models.APIKey
	result := s.db.Where("key = ?", key).First(&apiKey)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("API密钥不存在")
		}
		return nil, result.Error
	}

	return &apiKey, nil
}

// GetAPIKeysByUserID 获取指定用户的所有API密钥
func (s *APIKeyService) GetAPIKeysByUserID(userID uint) ([]models.APIKey, error) {
	var apiKeys []models.APIKey
	result := s.db.Where("user_id = ?", userID).Find(&apiKeys)
	if result.Error != nil {
		return nil, result.Error
	}

	return apiKeys, nil
}

// UpdateAPIKey 更新API密钥信息
func (s *APIKeyService) UpdateAPIKey(id uint, name string, active bool, expiresInDays int) (*models.APIKey, error) {
	var apiKey models.APIKey
	result := s.db.First(&apiKey, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("API密钥不存在")
		}
		return nil, result.Error
	}

	updates := map[string]interface{}{}

	if name != "" {
		updates["name"] = name
	}

	updates["active"] = active

	if expiresInDays > 0 {
		updates["expires_at"] = time.Now().AddDate(0, 0, expiresInDays)
	}

	result = s.db.Model(&apiKey).Updates(updates)
	if result.Error != nil {
		return nil, result.Error
	}

	// 重新获取更新后的记录
	s.db.First(&apiKey, id)

	return &apiKey, nil
}

// DeleteAPIKey 删除API密钥
func (s *APIKeyService) DeleteAPIKey(id uint) error {
	result := s.db.Delete(&models.APIKey{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("API密钥不存在")
	}

	return nil
}

// ValidateAPIKey 验证API密钥是否有效
func (s *APIKeyService) ValidateAPIKey(key string) (*models.APIKey, error) {
	apiKey, err := s.GetAPIKeyByKey(key)
	if err != nil {
		return nil, err
	}

	if !apiKey.IsValid() {
		return nil, errors.New("API密钥已过期或未激活")
	}

	return apiKey, nil
}
