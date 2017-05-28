package custom

import (
	//"errors"
	"encoding/json"
	"log"
	"mp"
	"net/url"
	"time"

	"github.com/toolkits/net/httplib"
)

func Send(msg interface{}, access_token string) (err error) {

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=" + url.QueryEscape(access_token)

	req := httplib.Post(incompleteURL).SetTimeout(3*time.Second, 1*time.Minute)
	req.Body(msg)
	resp, err := req.String()

	log.Println(msg, resp)

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
