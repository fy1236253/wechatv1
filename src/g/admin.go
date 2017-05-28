package g

import (
	"log"
	"mp"
	"sync"
)

var (
	adminLock = new(sync.RWMutex)
)

//IsAdmin 是否是管理员
func IsAdmin(openid string) bool {
	adminLock.RLock()
	defer adminLock.RUnlock()

	for _, c := range Config().Admins {
		if c.Openid == openid {
			return true
		}
	}
	return false
}

//SetAdmin 设置管理员
func SetAdmin(openid, nickname string) {
	if IsAdmin(openid) {
		return
	}

	adminLock.Lock()
	defer adminLock.Unlock()

	a := &mp.AdminsConfig{
		Openid:   openid,
		Nickname: nickname,
	}
	Config().Admins = append(Config().Admins, a)
	log.Println("add user to admins", openid)
}

// ExitAdmin 退出管理员模式
func ExitAdmin(openid string) {
	adminLock.Lock()
	defer adminLock.Unlock()

	var as []*mp.AdminsConfig
	for _, c := range Config().Admins {
		if c.Openid == openid {
			// 忽略掉
		} else {
			as = append(as, c)
		}
	}
	Config().Admins = as
	return
}
