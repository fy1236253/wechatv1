package model

import (
	"encoding/xml"
	"log"
	"mp"
	"mp/material"
	"mp/message"
	"mp/message/request"
	"net/http"
	"net/url"
	"util"
)

//WechatQueryParamsValid 检验微信回的消息是否完整
func WechatQueryParamsValid(m url.Values) {
	nonce := m.Get("nonce")
	timestamp := m.Get("timestamp")
	signature := m.Get("signature")
	msgSignature := m.Get("msg_signature")

	if nonce == "" {
		panic("nonce is nil")
	}
	if timestamp == "" {
		panic("timestamp is nil")
	}

	if signature == "" && msgSignature == "" {
		panic("signature and msg_signature is nil")
	}
}

//WechatSignValid 检查微信签名是否正确
func WechatSignValid(wxcfg *mp.WechatConfig, m url.Values) {
	nonce := m.Get("nonce")
	timestamp := m.Get("timestamp")
	signature := m.Get("signature")
	// log.Println(echostr, nonce, timestamp, signature)
	if util.Sign(wxcfg.Token, timestamp, nonce) == signature {
		return
	} else {
		panic("signature not match")
	}
}

// WechatStrValid 检查是否可用
func WechatStrValid(v, w, e string) {
	if v != w {
		panic(e)
	}
}

//WechatSignEncryptValid 加密模式的指纹验证方法
func WechatSignEncryptValid(wxcfg *mp.WechatConfig, m url.Values, body string) {
	nonce := m.Get("nonce")
	timestamp := m.Get("timestamp")
	signature := m.Get("msg_signature")
	//log.Println(echostr, nonce, timestamp, signature)
	if mp.MsgSign(wxcfg.Token, timestamp, nonce, body) == signature {
		return
	} else {
		panic("signature not match")
	}
}

//WechatMessageXMLValid 微信message验证
func WechatMessageXMLValid(req *http.Request, aesBody *message.AesRequestBody) {
	if err := xml.NewDecoder(req.Body).Decode(aesBody); err != nil {
		log.Println("[Warn] xml body", err)
		panic("xml body parse err")
	}
}

//WechatMessageXMLValidNormal 微信message验证
func WechatMessageXMLValidNormal(req *http.Request, normaleBody *message.NormalRequestBody) {
	if err := xml.NewDecoder(req.Body).Decode(normaleBody); err != nil {
		log.Println("[Warn] xml body", err)
		panic("xml body parse err")
	}
}

// ProcessWechatText 微信捕获文字消息
func ProcessWechatText(wxcfg *mp.WechatConfig, mixedMsg *message.MixedMessage) string {
	txt := request.GetText(mixedMsg)
	log.Println(txt)
	txtContent := txt.Content
	log.Println(txtContent)
	if txtContent == "只恐夜深花睡去" {
		SendMessageText(wxcfg.WxID, mixedMsg.FromUserName, "放下屠刀立地成佛！")
		go material.UpLodePIC("gh_f353e8a82fe5")
	} else if txtContent == "佛讲缘我讲钱" {
		SendMessageText(wxcfg.WxID, mixedMsg.FromUserName, "欢迎您管理员！")
	}
	return ""
}

//ProcessWechatEvent 处理微信推送的事件
func ProcessWechatEvent(wxcfg *mp.WechatConfig, mixedMsg *message.MixedMessage) {

	switch mixedMsg.Event {

	// 关注
	case request.EventTypeSubscribe:
		{

		}

	// 取消关注
	case request.EventTypeUnsubscribe:
		{

		}

	// 扫码事件
	case request.EventTypeScan:
		{ // 已经关注后 扫码  老用户 扫码 完成绑定

		}

	case request.EventTypeClick:
		{ // 菜单点击

		}

	// 给用户推送模板消息， 收到后的状态反馈， 需要推送到 open 平台、或业务系统
	case request.EventTypeTempSendOk:
		{ // 模板消息推送 ok

		}
	}

}
