/*
* @Author: thepoy
* @Email: thepoy@163.com
* @File Name: common.go
* @Created:  2021-01-17 19:22:50
* @Modified:  2021-01-20 18:21:13
 */

package utils

import (
	"blog_api/models"
	"io/ioutil"
	"os"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/yaml.v2"
)

// Config 通用配置
type Config struct {
	DBs    map[string]*DatabaseConfig `yaml:"databases"`
	Admins []string                   `yaml:"admins"`
	Email  *Email                     `yaml:"email"`
}

var (
	blogConfig *Config
	dbConfig   *DatabaseConfig
	email      *Email
	basicDir   string
)

func init() {
	blogConfig = &Config{}

	file, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(file, &blogConfig)
	if err != nil {
		panic(err)
	}

	blogAPIEnv := os.Getenv("BLOG_API_ENV")
	if blogAPIEnv == "" {
		blogAPIEnv = "development"
	}
	dbConfig = blogConfig.DBs[blogAPIEnv]
	email = blogConfig.Email

	basicDir, err = os.Getwd()
	if err != nil {
		panic(err)
	}

	initDB()
	initRedis()
}

// ErrorJSON 错误响应
func ErrorJSON(c *fiber.Ctx, statusCode int, err error) error {
	return c.Status(statusCode).JSON(&models.ErrorResponse{
		Error: err.Error(),
	})
}
