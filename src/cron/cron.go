package cron

import (
	"g"
	"log"
	"mq"
)

var (
	wechatWorkers []*mq.Consumer
)

func Start() {

	StartPackage()//调用每六小时查询快递

	StartToken() // 这个方法确保了 access token 可用

	// 这样检查一下  accesstoken 可用
	CheckToken()

	wechatWorkers = make([]*mq.Consumer, g.Config().Worker.Wechat)

	for i := range wechatWorkers {
		wechatWorkers[i] = WechatConsumer()
		go wechatWorkers[i].StartUp()
	}

	log.Println("cron.Start ok")

}

func Stop() {
	for i := range wechatWorkers {
		wechatWorkers[i].Shutdown()
	}

	log.Println("cron.Stop ok")
}
