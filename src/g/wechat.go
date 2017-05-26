package g

import (
	"log"
	"mp"
	"sync"
)

var (
	Wxcfg     map[string]*mp.WechatConfig
	wxcfgLock = new(sync.RWMutex)
)

//InitWxConfig 初始化WeChat
func InitWxConfig() {
	Wxcfg = make(map[string]*mp.WechatConfig)
	log.Println("g.InitWxConfig ok")
	for _, c := range Config().Wechats {
		Wxcfg[c.WxID] = c
	}

}

func GetWechatConfig(wxid string) *mp.WechatConfig {
	if wxid == "" {
		wxid = "gh_f353e8a82fe5"
	}
	wxcfgLock.RLock()
	defer wxcfgLock.RUnlock()
	return Wxcfg[wxid]
}
