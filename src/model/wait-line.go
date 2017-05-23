// 排队

package model

import (
	//"bytes"
	//"encoding/json"
	"log"
	//"mq"
	//"g"
	//"mp/message/custom"
	//"time"
	//"math/rand"
	//"strconv"

	"github.com/garyburd/redigo/redis"
	redispool "redis"

	//"database/sql"
)

type WaitLine struct {
	DevId string
	//Users []*User
}

// 根据设备返回  排队队列
func CreateLine(wxid, devid string) *WaitLine {

	line := &WaitLine{
		DevId: devid,
	}

	return line
}

func (self *WaitLine) AddUser(u *User) {
	rc := redispool.ConnPool.Get()
	defer rc.Close()

	key := "dev_" + self.DevId

	// 要考虑去掉 重复的

	a, _ := redis.Strings(rc.Do("LRANGE", key, 0, 100))
	for _, k := range a {
		if k == u.OpenId {
			return // 已经在队列中了
		}
	}

	rc.Do("LPUSH", key, u.OpenId)
	rc.Do("EXPIRE", key, 7200)

	log.Println("add user in line", key, u.OpenId)

	/*
		find, _ := redis.Bool(rc.Do("EXISTS",  key))

		if find {

		} else {

		}
	*/
	//smap, _ := redis.StringMap(rc.Do("HGETALL", openid))
}

func (self *WaitLine) GetUsers() (s []*User) {
	rc := redispool.ConnPool.Get()
	defer rc.Close()

	key := "dev_" + self.DevId

	a, _ := redis.Strings(rc.Do("LRANGE", key, 0, 100))
	for _, k := range a {
		u := CreateUser("", k)
		s = append(s, u)
	}
	return
}
