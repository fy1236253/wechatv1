package menu

import (
	//"errors"
	"encoding/json"
	"github.com/toolkits/net/httplib"
	"log"
	"mp"
	"net/url"
	"time"
)

// 创建自定义菜单.
func CreateMenu(obj interface{}, access_token string) (err error) {

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/menu/create?access_token=" + url.QueryEscape(access_token)

	req := httplib.Post(incompleteURL).SetTimeout(3*time.Second, 1*time.Minute)
	req.Body(obj)
	resp, err := req.String()

	log.Println(resp)

	if err != nil {
		log.Println("[ERROR]", err)
		return err
	}

	var result mp.Error
	err = json.Unmarshal([]byte(resp), &result)
	if result.ErrCode != mp.ErrCodeOK {
		log.Println("[ERROR]", result)
		return
	}
	return
}
