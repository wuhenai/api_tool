package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// APIKey 表示API密钥记录
type APIKey struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Key       string    `json:"key" gorm:"uniqueIndex;size:64;not null"`
	Name      string    `json:"name" gorm:"size:255;not null"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	ExpiresAt time.Time `json:"expires_at"`
	Active    bool      `json:"active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate 在创建记录前生成UUID作为API密钥
func (a *APIKey) BeforeCreate(tx *gorm.DB) error {
	if a.Key == "" {
		a.Key = uuid.New().String()
	}
	return nil
}

// IsValid 检查API密钥是否有效
func (a *APIKey) IsValid() bool {
	return a.Active && time.Now().Before(a.ExpiresAt)
}
