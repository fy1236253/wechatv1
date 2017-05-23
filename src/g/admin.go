//  管理员的 openid  管理
package g

import (
	//"bytes"
	//"encoding/json"
	//"mq"
	"log"
	"sync"
)

var (
	adminLock = new(sync.RWMutex)
)

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

func SetAdmin(openid, nickname string) {

	if IsAdmin(openid) {
		return
	}

	adminLock.Lock()
	defer adminLock.Unlock()

	a := &AdminConfig{
		Openid:   openid,
		Nickname: nickname,
	}

	Config().Admins = append(Config().Admins, a)
	log.Println("add user to admins", openid)
}

func ExitAdmin(openid string) {
	adminLock.Lock()
	defer adminLock.Unlock()

	var as []*AdminConfig
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
