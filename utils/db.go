/*
* @Author: thepoy
* @Email:  thepoy@aliyun.com
* @File Name: db.go
* @Created:   2021-01-17 19:03:04
* @Modified:  2021-01-18 16:50:12
 */

package utils

import (
	"blog_api/models"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db  *gorm.DB
	err error
)

// Mysql mysql配置信息
type Mysql struct {
	Username, Password, Host, Port, Database string
}

// Redis Redis配置信息
type Redis struct {
	Host, Port string
}

// DatabaseConfig 数据库模型
type DatabaseConfig struct {
	Mysql Mysql
	Redis Redis
}

func (dbc *DatabaseConfig) mysqlDsn() string {
	dsn := "%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local"
	return fmt.Sprintf(
		dsn,
		dbc.Mysql.Username,
		dbc.Mysql.Password,
		dbc.Mysql.Host,
		dbc.Mysql.Port,
		dbc.Mysql.Database,
	)
}

func initDB() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,  // 慢 SQL 阈值
			LogLevel:      logger.Error, // Log level
			Colorful:      false,        // 禁用彩色打印
		},
	)

	dsn := dbConfig.mysqlDsn()
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                 newLogger,
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic("Connecting database failed:" + err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("Creating `DB` failed:" + err.Error())
	}
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	db.AutoMigrate(&models.User{}, &models.BlogType{}, &models.Blog{}, &models.LikeBlog{}, &models.Favorite{}, &models.Comment{}, &models.LikeComment{})
}

// GetDB 调用gorm生成的数据库连接
func GetDB() *gorm.DB {
	return db
}
