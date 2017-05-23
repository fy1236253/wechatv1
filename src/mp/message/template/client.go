package template

import (
	//"errors"
	"encoding/json"
	"errors"
	"github.com/toolkits/net/httplib"
	"log"
	"mp"
	"net/url"
	"strconv"
	"time"
)

func Send(msg interface{}, access_token string) (msgid string, err error) {

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=" + url.QueryEscape(access_token)

	req := httplib.Post(incompleteURL).SetTimeout(3*time.Second, 1*time.Minute)
	req.Body(msg)
	resp, err := req.String()

	log.Println(resp) //  {"errcode":43004,"errmsg":"require subscribe hint: [qzd0547age6]"}

	if err != nil {
		log.Println("[ERROR]", err)
		return "", err
	}
	// {"errcode":0,"errmsg":"ok","msgid":401504797}
	var result mp.Error
	err = json.Unmarshal([]byte(resp), &result)
	if result.ErrCode != mp.ErrCodeOK {
		log.Println("[Warn]", result)
		return "", errors.New(strconv.Itoa(result.ErrCode))
	}
	msgid = strconv.FormatInt(result.MsgId, 10)
	return msgid, nil
}
