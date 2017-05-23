package model

import (
	//"bytes"
	"encoding/json"
	"encoding/xml"
	"g"
	"log"
	"mq"
	//"mp/message/custom"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"mp/menu"
	"mp/message"
	"mp/message/request"
	"mp/message/template"
	mpuser "mp/user"
	"mp/util"
	"net/http"
	"net/url"
	"open"
	redispool "redis"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// -------------------  valid 方法

// 验证 微信配置参数
func WxConfigValid(wxcfg *g.WechatConfig) {
	if wxcfg == nil {
		panic("[Warn] wecat config not find")
	}
}

// 微信回调 参数合法性 检验
func WechatQueryParamsValid(m url.Values) {
	nonce := m.Get("nonce")
	timestamp := m.Get("timestamp")
	signature := m.Get("signature")
	msg_signature := m.Get("msg_signature")

	if nonce == "" {
		panic("nonce is nil")
	}
	if timestamp == "" {
		panic("timestamp is nil")
	}
	if signature == "" && msg_signature == "" {
		panic("signature and msg_signature is nil")
	}
}

// 普通指纹验证
func WechatSignValid(wxcfg *g.WechatConfig, m url.Values) {
	nonce := m.Get("nonce")
	timestamp := m.Get("timestamp")
	signature := m.Get("signature")
	//log.Println(echostr, nonce, timestamp, signature)
	if util.Sign(wxcfg.Token, timestamp, nonce) == signature {
		return
	} else {
		panic("signature not match")
	}
}

func WechatStrValid(v, w, e string) {
	if v != w {
		panic(e)
	}
}

// 加密模式的指纹验证方法
func WechatSignEncryptValid(wxcfg *g.WechatConfig, m url.Values, body string) {
	nonce := m.Get("nonce")
	timestamp := m.Get("timestamp")
	signature := m.Get("msg_signature")
	//log.Println(echostr, nonce, timestamp, signature)
	if util.MsgSign(wxcfg.Token, timestamp, nonce, body) == signature {
		return
	} else {
		panic("signature not match")
	}
}

func WechatMessageXmlValid(req *http.Request, aesBody *message.AesRequestBody) {
	if err := xml.NewDecoder(req.Body).Decode(aesBody); err != nil {
		log.Println("[Warn] xml body", err)
		panic("xml body parse err")
	}
}

func WechatMessageXmlValidNormal(req *http.Request, normaleBody *message.NormalRequestBody) {
	if err := xml.NewDecoder(req.Body).Decode(normaleBody); err != nil {
		log.Println("[Warn] xml body", err)
		panic("xml body parse err")
	}
}

// ------------------

// openid, to, userdata 允许空
// msgfrom 标示 消息的来源   mq 或 开放平台
func WechatSendTemplate(wxid, uuid, openid, to, userdata string, in interface{}, msgfrom string, feed bool) {

	msgid, e := template.Send(in, g.GetWechatAccessToken(wxid))

	if e != nil {
		log.Println("template.Send", e)

		if msgfrom == "open" {
			open.Report(uuid, openid, to, "4", "模板消息推送失败")

			//} else if msgfrom == "rpc" {
			//rpc.Report(uuid, openid, to, "4", "模板消息推送失败")

		} else {
			SendFeedback("WechatFeedback", to, userdata, "0", "微信消息推送失败", uuid, openid) // 失败状态反馈

		}

	} else {
		log.Println("template.Send ok msgid", msgid)

		rc := redispool.ConnPool.Get()
		defer rc.Close()
		if feed { // 需要跟踪状态
			rc.Do("HMSET", msgid, "uuid", uuid, "msgfrom", msgfrom)
			rc.Do("EXPIRE", msgid, 3000) // 5分钟没有收到 消息回执
		} else {
			rc.Do("HMSET", msgid, "uuid", "noreport", "msgfrom", msgfrom)
			rc.Do("EXPIRE", msgid, 3000) // 5分钟没有收到 消息回执
		}

		if msgfrom == "open" {
			open.Report(uuid, openid, to, "3", "模板消息已经发送（是否收到未定）")

			//} else if msgfrom == "rpc" {
			//rpc.Report(uuid, openid, to, "3", "模板消息已经发送（是否收到未定）")

		} else {
			//SendFeedback("WechatFeedback", to, userdata, "180", "微信消息已发送", uuid)
		}

	}

}

// 处理 微信收到的文本消息
func ProcessWechatText(wxcfg *g.WechatConfig, mixedMsg *message.MixedMessage) string {
	txt := request.GetText(mixedMsg)

	user := CreateUser(wxcfg.WxId, mixedMsg.FromUserName)
	log.Println("User", user)

	//未完成绑定的用户，截获用户的输入信息，完成自动绑定
	if true { // 自动绑定模式
		if user.HaveMobile() == false {
			c := txt.Content

			mob := ""
			re := regexp.MustCompile("\\d{11,}")
			arr := re.FindStringSubmatch(c)
			if len(arr) > 0 {
				mob = arr[0]
			}

			rc := redispool.ConnPool.Get()
			defer rc.Close()

			sig := "XXXXXX"
			if n := strings.Index(c, "凭"); n > 0 {
				sig = c[n+3 : n+9]
			} else {
				if mob == "" && len(c) == 6 { // 用户只复制了 验证码
					key := "111_" + c // 111 约定前缀
					mob, _ = redis.String(rc.Do("HGET", key, "m"))
					sig = c
				}
			}

			if mob == "" && sig != "" {
				key := "111_" + sig // 111 约定前缀
				mob, _ = redis.String(rc.Do("HGET", key, "m"))
			}

			log.Println("debug=", c, mob, sig)

			if mob != "" && sig != "XXXXXX" {
				find, _ := redis.Bool(rc.Do("EXISTS", mob+"_"+sig))

				log.Println("debug=", c, mob, sig, find)

				if find == false { // 若 找不到  则只匹配 sig
					key := "111_" + sig // 111 约定前缀
					mob, _ = redis.String(rc.Do("HGET", key, "m"))
					find, _ = redis.Bool(rc.Do("EXISTS", mob+"_"+sig))
					log.Println("debug=", c, mob, sig, find)
				}

				if find { // 自动绑定了
					user.WaitSn(mob, sig)
					user.BindClose(true)
					log.Println("auto bind ok", mob, sig, c)

					oth, _ := redis.String(rc.Do("HGET", mob+"_"+sig, "oth"))
					if oth == "jiangli=true" {
						rc.Do("HMSET", mob+"_"+sig, "oth", "") // 置空 
						if SendBill(wxcfg.WxId, mixedMsg.FromUserName, mob, 1)  {
							SendMessageText(wxcfg.WxId, mixedMsg.FromUserName, "亲，恭喜您中奖话费1元，几分钟后将充值到"+ mob +"。" )
						}
					}
						
					return ""
				}
			}
		}
	}

	if user.IsDoBinding() { //  处理 绑定流程

		if user.IsWaitMobile() { // val 收到手机号码  或  粘贴的通知内容
			c := txt.Content // 亲13980901111有快递(abcdef)即

			mob := ""
			re := regexp.MustCompile("\\d{11,}")
			arr := re.FindStringSubmatch(c)

			if len(arr) > 0 {
				mob = arr[0]
			}
			user.SendSmsSn(mob) // 短信验证

		} else if user.IsWaitSn() {
			if user.ValidSn(txt.Content) { // 绑定成功 提示信息的内容
				//if user.Mobile2 == "" && user.Mobile3 == "" {  // 首次绑定的判断
				//if SendBill(wxcfg.WxId, mixedMsg.FromUserName, user.Mobile1, 1) {
				//	SendMessageText(wxcfg.WxId, mixedMsg.FromUserName, "亲，恭喜您中奖话费1元，几分钟后将充值到"+ user.Mobile1 +"。" )
				//} else {
				//	SendMessageText(wxcfg.WxId, mixedMsg.FromUserName, "亲，您手机号"+ user.Mobile1 )
				//}
				//}
			}
		}
		//obj = nil  //  绑定过程中  数据内部处理 不上报

	} else {
		// 拦截内部 处理命令
		if txt.Content == "芝麻开门" {
			SendMessageText(wxcfg.WxId, mixedMsg.FromUserName, "天王盖地虎，下句是啥？")

		} else if txt.Content == "宝莲灯" {
			SendMessageText(wxcfg.WxId, mixedMsg.FromUserName, "管理员你好，你可以删除指定手机号码的绑定关系，回复 no + 手机号，（只清理绑定手机号码，用户数据不删除）")
			g.SetAdmin(mixedMsg.FromUserName, "nickname")

		} else if txt.Content == "芝麻关门" {
			SendMessageText(wxcfg.WxId, mixedMsg.FromUserName, "已退出管理模式")
			g.ExitAdmin(mixedMsg.FromUserName)

		} else if g.IsAdmin(mixedMsg.FromUserName) {
			if strings.ToUpper(txt.Content)[0:2] == "NO" {
				cont := strings.TrimSpace(txt.Content[2:]) // 自动提取开头的 编号
				re := regexp.MustCompile("\\d+")
				uid := re.FindString(cont)
				if ClearMobile(wxcfg.WxId, uid, mixedMsg.FromUserName) { //  清除手机号码
					SendMessageText(wxcfg.WxId, mixedMsg.FromUserName, "已经清除")
				} else {
					SendMessageText(wxcfg.WxId, mixedMsg.FromUserName, "用户手机号码没有找到")
				}
			}
		} else {
			//obj = txt  // 普通的 会话内容，上报给业务系统
			SendFeedback("WechatFeedback", user.Mobile1, "", "2", txt.Content, "", mixedMsg.FromUserName)

			//if txt.Content == "王东测试" {
			out := fmt.Sprintf("<xml><ToUserName><![CDATA[%s]]></ToUserName><FromUserName><![CDATA[%s]]></FromUserName><CreateTime>%d</CreateTime><MsgType><![CDATA[transfer_customer_service]]></MsgType></xml>",
				mixedMsg.FromUserName, wxcfg.WxId, time.Now().Unix())
			return out
			//}

		}
	}

	return ""
}

// 处理 微信收到的消息
func ProcessWechatEvent(wxcfg *g.WechatConfig, mixedMsg *message.MixedMessage) {

	switch mixedMsg.Event {

	// 关注
	case request.EventTypeSubscribe:
		{

			// 扫码 后关注
			obj := request.GetSubscribeByScanEvent(mixedMsg)
			eventkey := ""
			if obj.EventKey != "" {
				eventkey = obj.EventKey
				// 上报 扫码信息
				log.Println("todo", obj)

			} else { // 普通关注
				obj := request.GetSubscribeEvent(mixedMsg)
				log.Println("todo", obj)
				// 上报关注信息  todo
			}

			go func() { // 获取用户信息
				u, e := mpuser.GetUserInfo(g.GetWechatAccessToken(wxcfg.WxId), mixedMsg.FromUserName, "")
				if e == nil {
					// redis 中获取 已经存在的用户数据
					user := CreateUser(wxcfg.WxId, mixedMsg.FromUserName)
					user.ImgUrl = u.HeadImageURL
					user.NickName = u.Nickname
					u.Mobile = user.Mobile1 // 暂时只上报一个号码
					u.Cmd = "WechatSubscribe"
					u.EventKey = strings.Replace(eventkey, "qrscene_", "", -1) // 记录扫码中的 场景值 qrscene_15882017353
					user.Save()
					bs, _ := json.Marshal(u)
					tmpjson := string(bs)
					// 收到的数据 上报给业务系统
					mq.Publish("wechat.out.exchange", "direct", "out.key", tmpjson, true)
					log.Println("userinfo json push to mq:", tmpjson)
				}
			}()

			// 判断是否是h5 新关注用户
			if true {
				user := CreateUser(wxcfg.WxId, mixedMsg.FromUserName)
				if user.Data == "new" { // 只有新用户 空手机 才发 ， 发完就不发了
					log.Println("send HongBao to", user.Mobile1)

					//if SendBill(wxcfg.WxId, mixedMsg.FromUserName, user.Mobile1, 1)  {
					//	SendMessageText(wxcfg.WxId, mixedMsg.FromUserName, "亲，恭喜您中奖话费1元，几分钟后将充值到"+ user.Mobile1 +"。" )
					//} else {
					//	SendMessageText(wxcfg.WxId, mixedMsg.FromUserName, "亲，您手机号"+ user.Mobile1 +"。" )
					//	// 更多分享 http://wechat2.shenbianvip.com/h5/r
					//}

				} else { // 老用户 重新关注
					SendMessageText(wxcfg.WxId, mixedMsg.FromUserName, wxcfg.Welcome) // 发送普通 文本消息

					if user.Mobile1 == "" { // 直接进入 绑定模式
						go func() {
							time.Sleep(2000 * time.Millisecond)
							user.SendBindWelcom() // 发送 绑定 提示 信息
						}()
					}
				}

				user.Sub = 1
				user.Save()
			}

		}

	// 取消关注
	case request.EventTypeUnsubscribe:
		{
			obj := request.GetUnsubscribeEvent(mixedMsg)
			//log.Println("todo", obj)
			var un struct {
				Cmd    string `json:"cmd"` // 扩展
				OpenId string `json:"openid"`
			}
			un.Cmd = "WechatUnSubscribe"
			un.OpenId = obj.FromUserName

			bs, _ := json.Marshal(un)
			tmpjson := string(bs)
			// 收到的数据 上报给业务系统
			mq.Publish("wechat.out.exchange", "direct", "out.key", tmpjson, true)
			log.Println("push to mq:", tmpjson)

			user := CreateUser(wxcfg.WxId, mixedMsg.FromUserName)
			user.Sub = 0 // 更新用户 为未关注状态
			user.Save()
		}

	// 扫码事件
	case request.EventTypeScan:
		{ // 已经关注后 扫码  老用户 扫码 完成绑定
			obj := request.GetScanEvent(mixedMsg)
			log.Println("todo", obj)
			if obj.EventKey == "123" { // h5活动
				user := CreateUser(wxcfg.WxId, mixedMsg.FromUserName)
				//if SendBill(wxcfg.WxId, mixedMsg.FromUserName, user.Mobile1, 1){
				//	SendMessageText(wxcfg.WxId, mixedMsg.FromUserName, "亲，恭喜您中奖话费1元，几分钟后将充值到"+ user.Mobile1 +"。" )
				//} else {
				//	SendMessageText(wxcfg.WxId, mixedMsg.FromUserName, "亲，您手机号"+ user.Mobile1 )
				//}

				user.Sub = 1
				user.Save() // 扫码用户 一定是关注用户了
			}
		}

	case request.EventLocationSelect:
		{

			if mixedMsg.EventKey == "click_post" {
				//log.Println("x ", mixedMsg.Latitude, " y", mixedMsg.Longitude)
				sendGetCourierInfo(strconv.FormatFloat(mixedMsg.SendLocationInfo.LocationX, 'f', -1, 64),
					strconv.FormatFloat(mixedMsg.SendLocationInfo.LocationY, 'f', -1, 64), wxcfg.WxId, mixedMsg.FromUserName)
			}

		}

	// 位置信息上报
	case request.EventTypeLocation:
		{
			obj := request.GetLocationEvent(mixedMsg)
			log.Println("todo", obj)

		}

	case request.EventTypeClick:
		{ // 菜单点击
			tmp := menu.GetClickEvent(mixedMsg)

			if tmp.IsInnerEvent() { // 只处理内部命令
				if tmp.EventKey == "inner_click_search" {
					user := CreateUser(wxcfg.WxId, mixedMsg.FromUserName)

					if user.Mobile1 == "" {
						SendMessageText(wxcfg.WxId, mixedMsg.FromUserName, "请先关联手机号码，然后才能查询快件！")
					} else {
						//if false {
						SendFeedback("WechatSearchPackage", "", "", "", "包裹查询", "", mixedMsg.FromUserName)
						//} else {
						//	arr := user.GetPackage()
						//	log.Println("all pkg", arr)
						//	for _, p := range arr {
						//		SendMessageText(wxcfg.WxId, mixedMsg.FromUserName, p)
						//	}
						//	if len(arr) == 0 {
						//		SendMessageText(wxcfg.WxId, mixedMsg.FromUserName, "没有找到您的包裹信息")
						//	}
						//}

					}
				}
				if tmp.EventKey == "inner_click_about" {
					SendMessageText(wxcfg.WxId, mixedMsg.FromUserName, ""+
						"云喇叭快递小管家，为您提供快递通知、查询、寄件等服务"+
						"..."+
						"")
				}
				if tmp.EventKey == "inner_click_bind" {
					user := CreateUser(wxcfg.WxId, mixedMsg.FromUserName)
					log.Println("User", user)
					user.SendBindWelcom() // 发送 绑定 提示 信息
				}

				if tmp.EventKey == "inner_click_jifen" {
					//user := CreateUser(wxcfg.WxId, )
					SendFeedback("WechatSearchJifen", "", "", "", "查询积分", "", mixedMsg.FromUserName)
				}

				//obj = nil  // 事件内部处理 不上报了

			} else {
				obj := tmp
				log.Println("todo", obj)

			}
		}

	//case request.EventTypeView:    //  url 导航
	//	obj = request.(mixedMsg)

	// 给用户推送模板消息， 收到后的状态反馈， 需要推送到 open 平台、或业务系统
	case request.EventTypeTempSendOk:
		{ // 模板消息推送 ok

			go func() {
				tmsgevent := template.GetTemplateSendJobFinishEvent(mixedMsg)

				msgid := strconv.FormatInt(tmsgevent.MsgId, 10)

				var uuid string
				var msgfrom string

				uuid = GetUuidByTMsgId(msgid)
				msgfrom = GetMsgFromByTMsgId(msgid) // 消息来源

				if uuid == "noreport" {
					return // 推送事件不需要上报
				}

				// 这uuid 可能为空
				for i := 0; i < 3; i++ {
					if uuid != "" {
						break
					}
					time.Sleep(20 * time.Millisecond)
					uuid = GetUuidByTMsgId(msgid)
					msgfrom = GetMsgFromByTMsgId(msgid) // 消息来源
				}

				log.Println("msgfrom", msgfrom, "uuid", uuid, "msgid", msgid)

				if tmsgevent.Status == template.TemplateSendStatusSuccess {
					if msgfrom == "open" {
						open.Report(uuid, "", "", "5", "模板消息发用户已经收到（是否读未知）")

						//} else if msgfrom == "rpc" {
						//rpc.Report(uuid, "", "", "5", "模板消息发用户已经收到（是否读未知）")

					} else {
						SendFeedback("WechatFeedback", "", "", "1", "微信通知成功", uuid, mixedMsg.FromUserName)
					}

				} else {
					if msgfrom == "open" {
						open.Report(uuid, "", "", "4", "模板消息发送不成功")

						//} else if msgfrom == "rpc" {
						//rpc.Report(uuid, "", "", "4", "模板消息发送不成功")

					} else {
						SendFeedback("WechatFeedback", "", "", "0", "微信通知失败", uuid, mixedMsg.FromUserName)
					}

				}
			}()

		}
	}

}

func GetUuidByTMsgId(msgid string) string {
	rc := redispool.ConnPool.Get()
	defer rc.Close()

	uuid, _ := redis.String(rc.Do("HGET", msgid, "uuid"))
	return uuid
}

func GetMsgFromByTMsgId(msgid string) string {
	rc := redispool.ConnPool.Get()
	defer rc.Close()

	f, _ := redis.String(rc.Do("HGET", msgid, "msgfrom"))
	return f
}

/*
// 消息处理
				var obj interface{}

				if obj != nil {   //  事件需要上报给业务处理
					bs, _ := json.Marshal(obj)
					tmpjson := string(bs)

					// 收到的数据 上报给业务系统
					mq.Publish("wechat.out.exchange", "direct", "out.key", tmpjson, true)
					log.Println("status push to mq:", tmpjson)
				}
*/

// 通过 微信接口，获取所有用户信息，然后上报到业务系统
func UsersListSync(wxid string) {

	nextopenid := ""

	for {
		ulist, err := mpuser.GetUserList(g.GetWechatAccessToken(wxid), nextopenid)

		if err != nil {
			break
		}

		nextopenid = ulist.NextOpenId

		//log.Println(ulist.Data)

		for _, v := range ulist.Data.OpenIdList {

			u, e := mpuser.GetUserInfo(g.GetWechatAccessToken(wxid), v, "")
			if e != nil {
				continue
			}

			// redis 中获取 已经存在的用户数据
			user := CreateUser(wxid, v)
			user.NickName = u.Nickname
			user.Sex = u.Sex
			user.Sub = 1 // 关注状态为1
			user.Save()

			u.Mobile = user.Mobile1 // 暂时只上报一个号码
			u.Cmd = "WechatSubscribe"
			u.EventKey = "syncuser" // 记录扫码中的 场景值

			bs, _ := json.Marshal(u)
			tmpjson := string(bs)
			// 收到的数据 上报给业务系统
			mq.Publish("wechat.out.exchange", "direct", "out.key", tmpjson, true)
			log.Println("userinfo json push to mq:", tmpjson)

		}

		if nextopenid == "" {
			break
		}
	}

}
