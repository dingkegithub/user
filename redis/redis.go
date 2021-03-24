package redis

import (
	"fmt"
	"time"

	"github.com/go-redsync/redsync"
	"github.com/gomodule/redigo/redis"
)

var pool *redis.Pool
var redisLock *redsync.Redsync

func InitRedis(host string, port string, password string) error {
	pool := &redis.Pool{
		MaxIdle:     20,
		IdleTimeout: 240 * time.Second,
		MaxActive:   50,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	redisLock = redsync.New([]redsync.Pool{pool})
	return nil
}

func GetRedisConn() (redis.Conn, error) {
	conn := pool.Get()
	return conn, conn.Err()
}

func GetLock(key string, expire time.Duration) *redsync.Mutex {
	return redisLock.NewMutex(key, redsync.SetExpiry(expire))
}
