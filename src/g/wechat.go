package g

import (
	"log"
	"mp"
	"sync"
)

var (
	// Wxcfg 全局的微信配置
	Wxcfg       map[string]*mp.WechatConfig
	wxcfgLock   = new(sync.RWMutex)
	wxTokenLock = new(sync.RWMutex)
)

//InitWxConfig 初始化WeChat
func InitWxConfig() {
	Wxcfg = make(map[string]*mp.WechatConfig)
	log.Println("g.InitWxConfig ok")
	for _, c := range Config().Wechats {
		Wxcfg[c.WxID] = c
	}
}

// GetWechatConfig 通过wxid获取配置信息
func GetWechatConfig(wxid string) *mp.WechatConfig {
	wxcfgLock.RLock()
	defer wxcfgLock.RUnlock()
	return Wxcfg[wxid]
}

// GetWechatAccessToken 通过wxid获取accesstoken（进程中）
func GetWechatAccessToken(wxid string) string {
	wxTokenLock.RLock()
	defer wxTokenLock.RUnlock()
	c := GetWechatConfig(wxid)
	if c == nil {
		return ""
	} else {
		return c.AccessToken
	}
}

// SetWechatAccessToken 设置accesstoken
func SetWechatAccessToken(wxid, token string) {
	wxTokenLock.Lock()
	defer wxTokenLock.Unlock()
	c := GetWechatConfig(wxid)
	c.AccessToken = token
}

// GetJsAPITicket 获取apiticket
func GetJsAPITicket(wxid string) string {
	wxTokenLock.RLock()
	defer wxTokenLock.RUnlock()

	c := GetWechatConfig(wxid)
	if c == nil {
		return ""
	} else {
		return c.JsapiTicket
	}
}

//SetJsAPITicket 设置apiticket
func SetJsAPITicket(wxid, ticket string) {
	wxTokenLock.Lock()
	defer wxTokenLock.Unlock()
	c := GetWechatConfig(wxid)
	c.JsapiTicket = ticket
}
