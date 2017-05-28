package model

import (
	"bytes"
	"encoding/json"
	"g"
	"log"
	"mp/message/custom"
)

//JSONHead 命令码
type JSONHead struct {
	Cmd  string `json:"cmd,omitempty"`
	WxID string `json:"wxid,omitempty"`
	UUID string `json:"uuid,omitempty"` // 如果是异步消息，需要 uuid 匹配发送和接收
}

//  -------------   客服消息 微信接口 ------------

//SendMessageText 给用户发送 普通文本消息  客服消息接口
func SendMessageText(wxid, openid, content string) {

	obj := custom.NewText(wxid, openid, content, "")

	buf := bytes.NewBuffer(make([]byte, 0, 16<<10))
	buf.Reset()
	json.NewEncoder(buf).Encode(obj)
	tmpjson := buf.String()

	go custom.Send(tmpjson, g.GetWechatAccessToken(wxid))
}

//SendMessageNews 发送客服消息 图文消息
func SendMessageNews(wxid, openid, title, desc, url, pic string) {

	art := custom.Article{
		Title:       title,
		Description: desc,
		URL:         url,
		PicURL:      pic,
	}

	articles := []custom.Article{art}

	obj := custom.NewNews(wxid, openid, articles, "")

	buf := bytes.NewBuffer(make([]byte, 0, 16<<10))
	buf.Reset()
	json.NewEncoder(buf).Encode(obj)
	tmpjson := buf.String()
	log.Println(tmpjson)
	go custom.Send(tmpjson, g.GetWechatAccessToken(wxid))
}

//SendMessagePic 发送客服消息 图文消息
func SendMessagePic(wxid, openid, mediaid, pic string) {
	obj := custom.NewImage(wxid, openid, mediaid, "")

	buf := bytes.NewBuffer(make([]byte, 0, 16<<10))
	buf.Reset()
	json.NewEncoder(buf).Encode(obj)
	tmpjson := buf.String()

	go custom.Send(tmpjson, g.GetWechatAccessToken(wxid))
}

//SendMessageVedio 发送客服消息 图文消息
func SendMessageVedio(wxid, openid, content string) {

	obj := custom.NewText(wxid, openid, content, "")

	buf := bytes.NewBuffer(make([]byte, 0, 16<<10))
	buf.Reset()
	json.NewEncoder(buf).Encode(obj)
	tmpjson := buf.String()

	go custom.Send(tmpjson, g.GetWechatAccessToken(wxid))
}

// ------------ 微信模板消息 ---------------

//SendSmsNotify 不同的消息入口，状态返回方式是不同的， 所以用一个 参数来标示， 消息的来源  msgfrom  , 默认来源于mq
func SendSmsNotify(wxid, uuid string, in, strurl, msgfrom string) {

}
