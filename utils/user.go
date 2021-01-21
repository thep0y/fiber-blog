/*
* @Author: thepoy
* @Email: thepoy@163.com
* @File Name: user.go (c) 2021
* @Created:  2021-01-18 16:56:27
* @Modified: 2021-01-21 16:01:44
 */

package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// IsAdmin 新注册或当前登录的用户是否为admin
func IsAdmin(username string) bool {
	for _, u := range blogConfig.Admins {
		if username == u {
			return true
		}
	}
	return false
}

// GeneratePassword 生成密码
func GeneratePassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPassword 检查密码
func CheckPassword(hashedPwd string, pwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(pwd))
	if err != nil {
		return false
	}
	return true
}

// Email 网站邮箱
type Email struct {
	Email, Host, Password string
	Port                  int
}

// ValidateCodeIsValid 验证码是否有效
func ValidateCodeIsValid(email, code string) bool {
	conn := redisPool.Get()
	defer conn.Close()

	res, err := redis.String(conn.Do("Get", email))
	if err != nil {
		return false
	}

	if code != res {
		return false
	}

	// 验证成功后，不管是否成功删除，都执行一次删除操作
	go func() {
		conn.Do("DEL", email)
	}()

	return true
}

// SetVerifyCodeOnRedis 将邮箱验证码写入redis里，过期时间为5分钟
func SetVerifyCodeOnRedis(email, code string) error {
	var expire uint = 60 * 5 // 验证码默认过期时间为5分钟

	if err := SetOnRedis(email, code, expire); err != nil {
		return err
	}
	return nil
}

// DatabaseExistError 返回注册时的错误信息
func DatabaseExistError(errStr string) error {
	errInfo := strings.Split(errStr, "'")

	errCol := strings.Split(errInfo[3], ".")
	var res string
	if strings.Contains(errInfo[3], ".") {
		res = fmt.Sprintf("The `%s` you entered alreadyd exists: [ %s ]", errCol[1], errInfo[1])
	} else {
		res = fmt.Sprintf("The `%s` you entered alreadyd exists: [ %s ]", errCol[0], errInfo[1])
	}
	return errors.New(res)
}

// GenerateToken 生成token
func GenerateToken() string {
	u := uuid.NewV4()
	return u.String()
}

// SetToken 向redis里写入token，key为uuid，value为用户id, expire小于等于0时，默认过期时间为24小时
func SetToken(token string, userID uint, expire uint) error {
	// TODO: 每一个id只能对应一个token，每个用户在同一时间只能有一个有效的已登录的终端。用前缀来遍历redis里的值并删除对应的值？
	if expire <= 0 {
		expire = 60 * 60 * 24
	}
	if err := SetOnRedis(token, fmt.Sprint(userID), expire); err != nil {
		return err
	}
	return nil
}
