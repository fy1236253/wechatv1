package shakearound

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/toolkits/net/httplib"
	"log"
	"mp"
	"net/url"
	"time"
)

// 摇一摇后， 获取 ticket ，然后通过 ticket 换取设备信息与用户信息
func GetShakeInfo(ticket, access_token string) (info *Shakeinfo, err error) {

	obj := struct {
		Ticket  string `json:"ticket"`
		NeedPoi int    `json:"need_poi,omitempty"`
	}{
		Ticket: ticket,
	}

	buf := bytes.NewBuffer(make([]byte, 0, 16<<10))
	buf.Reset()
	json.NewEncoder(buf).Encode(obj)
	tmpjson := buf.String()

	incompleteURL := "https://api.weixin.qq.com/shakearound/user/getshakeinfo?access_token=" + url.QueryEscape(access_token)
	req := httplib.Post(incompleteURL).SetTimeout(3*time.Second, 1*time.Minute)
	req.Body(tmpjson)
	resp, err := req.String()
	log.Println(resp)

	if err != nil {
		log.Println("[ERROR]", err)
		return nil, err
	}

	var result struct {
		mp.Error
		Shakeinfo `json:"data"`
	}

	err = json.Unmarshal([]byte(resp), &result)
	if result.ErrCode != 0 {
		log.Println("[ERROR]", result)
		return nil, errors.New(result.ErrMsg)
	}

	return &result.Shakeinfo, nil
}
