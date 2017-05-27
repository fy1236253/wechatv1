package redis

import (
	"g"
	"log"
	"time"

	"github.com/garyburd/redigo/redis"
)

var ConnPool *redis.Pool
var ConnPoolLocalNet *redis.Pool // 访问 局域网内的 redis

func InitConnPool() {

	ConnPool = &redis.Pool{
		MaxIdle:     g.Config().Redis.MaxIdle,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp",
				g.Config().Redis.Addr,
				redis.DialPassword(g.Config().Redis.Password),
				redis.DialDatabase(g.Config().Redis.Db))
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: PingRedis,
	}

	ConnPoolLocalNet = &redis.Pool{
		MaxIdle:     g.Config().RedisLocalNet.MaxIdle,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp",
				g.Config().RedisLocalNet.Addr,
				redis.DialPassword(g.Config().RedisLocalNet.Password),
				redis.DialDatabase(g.Config().RedisLocalNet.Db))
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: PingRedis,
	}
}

func PingRedis(c redis.Conn, t time.Time) error {
	_, err := c.Do("ping")
	if err != nil {
		log.Println("[ERROR] ping redis fail", err)
	}
	return err
}
