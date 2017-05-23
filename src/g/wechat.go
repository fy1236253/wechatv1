package g

import (
	//"bytes"
	//"encoding/json"
	//"mq"
	"log"
	"sync"
)

var (
	wxcfg       map[string]*WechatConfig
	wxcfgLock   = new(sync.RWMutex)
	wxTokenLock = new(sync.RWMutex)
)

//
//type WechatConfig struct {
//	WxId  	string
//	AppId   string
//	Token 	string
//	Aeskey 	string
//}

func InitWxConfig() {
	wxcfg = make(map[string]*WechatConfig)
	log.Println("g.InitWxConfig ok")

	//cfg := &WechatConfig{
	//	WxId: "",
	//	AppId: "",
	//	Token: "",
	//	Aeskey: "",
	//}
	//wxcfg[""] = cfg

	for _, c := range Config().Wechats {
		wxcfg[c.WxId] = c
	}

}

func GetWechatConfig(wxid string) *WechatConfig {
	if wxid == "" {
		wxid = "gh_8ac8a8821eb9"
	}

	wxcfgLock.RLock()
	defer wxcfgLock.RUnlock()

	return wxcfg[wxid]
}

func SetWechatAccessToken(wxid, token string) {
	wxTokenLock.Lock()
	defer wxTokenLock.Unlock()
	c := GetWechatConfig(wxid)
	c.accessToken = token
}

func GetWechatAccessToken(wxid string) string {
	wxTokenLock.RLock()
	defer wxTokenLock.RUnlock()

	c := GetWechatConfig(wxid)
	if c == nil {
		return ""
	} else {
		return c.accessToken
	}

}

//  jsapi_ticket 管理
func SetJsApiTicket(wxid, ticket string) {
	wxTokenLock.Lock()
	defer wxTokenLock.Unlock()

	c := GetWechatConfig(wxid)
	c.jsapi_ticket = ticket
}

func GetJsApiTicket(wxid string) string {
	wxTokenLock.RLock()
	defer wxTokenLock.RUnlock()

	c := GetWechatConfig(wxid)
	if c == nil {
		return ""
	} else {
		return c.jsapi_ticket
	}

}
