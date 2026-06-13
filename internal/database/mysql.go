// Package database 数据库连接初始化。
package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/wxsimon2022/simonStu/config"
)

// DB 全局 GORM 实例。包外部通过 database.DB 直接使用。
var DB *gorm.DB

// InitMySQL 根据配置连接 MySQL。连接失败时仅打日志，不阻塞启动。
func InitMySQL(cfg *config.Config) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=3s&readTimeout=3s&writeTimeout=3s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	var err error
	DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            false,
	})
	if err != nil {
		log.Printf("MySQL 连接失败: %v", err)
		return
	}

	// 配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("MySQL 获取底层连接失败: %v", err)
		return
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	log.Printf("MySQL 连接成功 (%s:%s/%s)", cfg.DBHost, cfg.DBPort, cfg.DBName)
}
