/*
* @Author: thepoy
* @Email: thepoy@163.com
* @File Name: redis.go (c) 2021
* @Created:  2021-01-18 16:14:03
* @Modified: 2021-01-20 19:08:28
 */

package utils

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

var redisPool *redis.Pool

func initRedis() {
	redisPool = &redis.Pool{
		MaxIdle:     20,
		MaxActive:   100,
		IdleTimeout: 100,
		Dial: func() (redis.Conn, error) {
			addr := fmt.Sprintf("%s:%s", dbConfig.Redis.Host, dbConfig.Redis.Port)
			return redis.Dial("tcp", addr)
		},
	}
}

// SetOnRedis redis的set方法，过期时间expire可选
func SetOnRedis(key, value string, expire ...uint) error {
	conn := redisPool.Get()
	defer conn.Close()

	_, err := redis.String(conn.Do("Set", key, value))
	if err != nil {
		return err
	}

	if len(expire) <= 0 {
		return nil
	} else if len(expire) > 1 {
		panic("过期时间只能有一个数字")
	} else {
		_, err = redis.Uint64(conn.Do("EXPIRE", key, expire[0]))
		return err
	}
}

// GetFromRedis 从redis中取数据。
// res可能会是其他类型(int或float)，但本方法只返回string，
// 需要其他类型时需要再转换。
func GetFromRedis(key string) (string, error) {
	conn := redisPool.Get()
	defer conn.Close()

	res, err := redis.String(conn.Do("Get", key))
	if err != nil {
		return "", err
	}

	return res, nil
}

// ExpireOnRedis 设置过期时间
func ExpireOnRedis(conn redis.Conn, key string, expire uint) error {
	_, err := redis.Uint64(conn.Do("EXPIRE", key, expire))
	return err
}
