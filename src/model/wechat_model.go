package model

import (
	"mp"
	"net/url"
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
func WechatSignValid(wxcfg *mp.WechatConfig, m url.Values) {
	nonce := m.Get("nonce")
	timestamp := m.Get("timestamp")
	signature := m.Get("signature")
	//log.Println(echostr, nonce, timestamp, signature)
	if mp.Sign(wxcfg.Token, timestamp, nonce) == signature {
		return
	} else {
		panic("signature not match")
	}
}
