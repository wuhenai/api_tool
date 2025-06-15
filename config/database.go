package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/yoyo/api_bot/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupDatabase 初始化数据库连接并进行迁移
func SetupDatabase() *gorm.DB {
	// 确保data目录存在
	if err := os.MkdirAll("data", os.ModePerm); err != nil {
		log.Fatalf("创建数据目录失败: %v", err)
	}

	dbPath := filepath.Join("data", "apikeys.db")

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	// 自动迁移数据库结构
	if err := db.AutoMigrate(&models.APIKey{}); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	log.Println("数据库初始化成功")
	return db
}
