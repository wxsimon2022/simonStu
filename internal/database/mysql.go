// Package database 数据库连接初始化。
package database

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/wxsimon8888/simonStu/config"
)

// DB 全局 GORM 实例。包外部通过 database.DB 直接使用。
var DB *gorm.DB

// InitMySQL 根据配置连接 MySQL。连接失败时仅打日志，不阻塞启动。
func InitMySQL(cfg *config.Config) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("MySQL 连接失败: %v", err)
		return
	}

	log.Printf("MySQL 连接成功 (%s:%s/%s)", cfg.DBHost, cfg.DBPort, cfg.DBName)
}
