package open

import (
	"g"
	"math/rand"
	"strconv"

	"time"

	//"net/url"
	"encoding/json"
	"github.com/toolkits/net/httplib"
	"log"
	//"errors"
	"bytes"
)


// 服务号 消息状态，上报到 open 平台 
func Report(uuid, openid, to, typeid, typemsg string) (err error) {

	appid := g.Config().Open.AppId
	token := g.Config().Open.AppToken
	nonce := strconv.Itoa(rand.Intn(999999999))
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	signature := Sign(token, timestamp, nonce, appid)

	var r ReportMsg

	r.Openid = openid
	r.Uuid = uuid
	r.To = to
	r.Resulttype = typeid
	r.Resultmsg = typemsg

	buf := bytes.NewBuffer(make([]byte, 0, 16<<10))
	buf.Reset()
	json.NewEncoder(buf).Encode(r)
	tmpjson := buf.String()
	log.Println(tmpjson)

	incompleteURL := g.Config().Open.Addr + "?appid=" + appid + "&nonce=" + nonce + "&timestamp=" + timestamp + "&signature=" + signature
	log.Println(incompleteURL)
	req := httplib.Post(incompleteURL).SetTimeout(3*time.Second, 1*time.Minute)
	req.Body(tmpjson)
	resp, err := req.String()

	log.Println(resp)

	if err != nil {
		log.Println("[ERROR]", err)
		return err
	}

	/*
		var result mp.Error
		err = json.Unmarshal([]byte(resp), &result)
		if result.ErrCode != mp.ErrCodeOK {
			log.Println("[ERROR]", result)
			return
		}
	*/
	return nil

}




//向服务器请求服务，获取道周边快递员的信息
func ReportLocation(locationx, locationy, wxid, openid string) (err error) {
	appid := g.Config().Location.AppId
	token := g.Config().Location.AppToken
	nonce := strconv.Itoa(rand.Intn(999999999))
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	signature := Sign(token, timestamp, nonce, appid)

	var l LocationMsg

	l.Openid = openid
	l.LocationX = locationx
	l.LocationY = locationy
	l.WxId = wxid

	buf := bytes.NewBuffer(make([]byte, 0, 16<<10))
	buf.Reset()
	json.NewEncoder(buf).Encode(l)
	tmpjson := buf.String()
	log.Println(tmpjson)

	incompleteURL := g.Config().Location.Addr + "?appid=" + appid + "&nonce=" + nonce + "&timestamp=" + timestamp + "&signature=" + signature
	log.Println(incompleteURL)
	req := httplib.Post(incompleteURL).SetTimeout(1*time.Second, 1*time.Minute)
	req.Body(tmpjson)
	resp, err := req.String()

	log.Println(resp)

	if err != nil {
		log.Println("[ERROR]", err)
		return err
	}

	return nil

}
