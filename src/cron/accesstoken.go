package cron

import (
	"g"
	"log"
	redispool "redis"
	"time"
	"util"

	"github.com/garyburd/redigo/redis"
)

// 定时检查 token ， 如果发现要过期了，重现请求一个
func monitorToken() {
	for {
		rc := redispool.ConnPool.Get()
		log.Println(rc)
		for _, c := range g.Config().Wechats {
			key := "wx_acc_tkn_" + c.WxID // 判断 redis中这个token 是否存
			t, _ := redis.Int64(rc.Do("TTL", key))
			if t < 600 { //  即将过期
				token := util.GetToken(c.AppID, c.AppSecret)
				if token == nil {
					continue
				}
				log.Println("wx access token refresh", "***"+token.Token[12:20]+"***", token.ExpiresIn)
				rc.Do("HMSET", key, "token", token.Token)
				rc.Do("EXPIRE", key, token.ExpiresIn-100)   // 留一个保护间隔
				g.SetWechatAccessToken(c.WxID, token.Token) // 同时保存到 进程内部 提高访问速度

				strTicket := util.GetJsApiTicket(token.Token)
				log.Println("wx jsapiticket ", strTicket)
				rc.Do("HMSET", key, "jsapiticket", strTicket)
				g.SetJsAPITicket(c.WxID, strTicket)

			} else {
				//  每次 同步写入 进程内存中，  这样 多节点，任何一个节点更新后， 都可以实现同步
				token, _ := redis.String(rc.Do("hget", key, "token"))
				g.SetWechatAccessToken(c.WxID, token)
				//log.Println("get access token from redis", c.WxId, "***" + token[12:20] + "***")

				ticket, _ := redis.String(rc.Do("hget", key, "jsapiticket"))
				g.SetJsAPITicket(c.WxID, ticket)

			}
		}
		rc.Close()
		time.Sleep(3 * time.Second) // 定时检查 token
	}
}

// StartToken 开启自动检测token
func StartToken() {
	go monitorToken()
}

//CheckToken 确保所有  access token 都有效
func CheckToken() {
	for _, c := range g.Config().Wechats {
		for {
			if g.GetWechatAccessToken(c.WxID) == "" && c.WxID != "demo" {
				log.Println("[warn] access token not ready, wait", c.WxID)
				time.Sleep(1 * time.Second)
				continue
			}
			break
		}
	}
}
