// worker 监控 mq 队列，获取json 数据
package cron

import (
	"encoding/json"
	"g"
	"log"
	"mq"
	//"strings"

	//redispool "redis"
	"model"
	//"proc"
	"mp/message/custom"
	//"mp/message/template"
	"mp/menu"
)

func handler(in string) (string, error) {

	var head model.JsonHead

	e := json.Unmarshal([]byte(in), &head)
	if e != nil {
		log.Println("json 解析失败: %s", e)
		return "", nil // 吃掉错误
	}

	wxid := head.WxId
	cmd := head.Cmd

	if g.GetWechatConfig(wxid) == nil {
		log.Println("[warn] wxid not find", wxid)
		return "", nil
	}

	switch cmd {

	case "SendSmsNotify": // 最终发送 模板消息
		{
			log.Println("[mq json] head", head, "json", in)
			go model.SendSmsNotify(wxid, head.Uuid, in, "", "") // 发送内部格式的 模板消息
		}

	//case "admin_template": // 短信模板审核 发送客服类消息
	//	{
	//		log.Println("[mq json] head", head, "json", in)
	//		go model.SendCheckDataToAdmin(wxid, in)
	//	}

	case "message":
		{
			log.Println("[mq json] head", head, "json", in)
			go custom.Send(in, g.GetWechatAccessToken(wxid)) // 主动推消息接口 客服消息
		}

	//case "template":
	//	{
	//		log.Println("[mq json] head", head, "json", in)
	//		go model.WechatSendTemplate(wxid, head.Uuid, "", "", "", in, "") //  发送 微信格式的模板消息
	//	}

	case "menu": //  创建菜单接口
		{
			log.Println("[mq json] head", head, "json", in)
			go menu.CreateMenu(in, g.GetWechatAccessToken(wxid))
		}
	case "syncuser": // 同步用户 
		{
			log.Println("[mq json] head", head, "json", in)
			go model.UsersListSync(wxid)
		}
	}

	return "", nil
}

func WechatConsumer() *mq.Consumer {
	consumer, err := mq.NewConsumer(
		"wechat.in.exchange", // exchange
		"direct",
		"wechat.in.queue", // queue name
		"in.key",          // route key
		"wechat-consumer", // ctag   //simple-consumer
		handler)           // call back Fun
	g.FailOnError(err, "[ERROR] Consumer 启动失败")

	return consumer
}

func DebugConsumer() *mq.Consumer {
	con, _ := mq.NewConsumer(
		"wechat.out.exchange", // exchange
		"direct",
		"wechat.out.queue",      // queue name
		"out.key",               // route key
		"wechat-consumer-debug", // ctag   //simple-consumer
		handler)
	return con
}
