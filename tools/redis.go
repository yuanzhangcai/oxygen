package tools

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

// Redis redis组件
type Redis struct {
	*redis.Client
	prefix string
}

var client *Redis

// InitRedis 初始化redis
func InitRedis(server, password, prefix string) error {
	if client != nil {
		return nil
	}

	if server == "" {
		return errors.New("redis服务器地址为空。")
	}

	client = &Redis{
		Client: redis.NewClient(&redis.Options{
			Addr:     server,
			Password: password,
		}),
		prefix: prefix,
	}

	pong, err := client.Ping().Result()
	if err != nil {
		logrus.Error(pong)
		logrus.Error(err)
		client = nil
		return err
	}
	return nil
}

// GetRedis 获取redis实例
func GetRedis() *Redis {
	return client
}

// SetObject 设置redis对象
func (c *Redis) SetObject(key string, value interface{}, expire time.Duration) error {
	buf, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.Set(key, string(buf), expire).Err()
}

// GetObject 获取redis对象
func (c *Redis) GetObject(key string, value interface{}) error {
	ret := c.Get(key)
	if ret.Err() != nil {
		return ret.Err()
	}

	err := json.Unmarshal([]byte(ret.Val()), value)
	if err == nil {
		return err
	}

	return nil
}
