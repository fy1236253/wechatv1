package model

import (
	"bytes"
	"encoding/json"
	"g"
	"log"
	"mp/message/custom"
	"mp/message/template"
	"github.com/toolkits/net/httplib"
	"mq"
	"open"
	"strings"
	"time"
)

// json中的 命令码
type JsonHead struct {
	Cmd  string `json:"cmd,omitempty"`
	WxId string `json:"wxid,omitempty"`
	Uuid string `json:"uuid,omitempty"` // 如果是异步消息，需要 uuid 匹配发送和接收
}

// 充值 充流量 
type RechargeOrder struct {
	Cmd  string `json:"cmd,omitempty"`
	Uuid  string `json:"uuid,omitempty"`  // 唯一标示
	Phone string `json:"phone,omitempty"` // 电话号码
	Openid string `json:"openid,omitempty"` // 微信 openid
	Value string `json:"value,omitempty"` // 金额
	Bits  string `json:"bits,omitempty"`  // 流量包 大小 M 
	Type  string `json:"type,omitempty"`  // 话费 或 流量
	Timestamp int64 `json:"timestamp,omitempty"` // 下单时候的时间戳 整数 
}

// 统一的 短信 语音 通知格式
type SysImJson struct {
	Cmd         string `json:"cmd,omitempty"`
	WxId        string `json:"wxid,omitempty"`
	Openid      string `json:"openid,omitempty"`
	Uuid        string `json:"uuid,omitempty"`
	Sn          string `json:"sn,omitempty"`       // 包裹的自定义编号
	To          string `json:"to,omitempty"`       // 待呼叫号码
	Userdata    string `json:"userdata,omitempty"` // 只允许填写数字字符
	Name        string `json:"name,omitempty"`     // 快递员名字
	CompanyName string `json:"companyname,omitempty"`
	Note        string `json:"note,omitempty"`         // 编号 + 时间 + 地点。 用户自定义的内容
	Sid         string `json:"sid,omitempty"`          // 短链接使用的 id
	History     string `json:"is_history,omitempty"`   //  = "1" 历史消息
	HistoryAt   string `json:"receive_time,omitempty"` //通知时间
}

type SysFeedback struct {
	Cmd      string `json:"cmd,omitempty"`
	Uuid     string `json:"uuid,omitempty"`
	Openid   string `json:"openid,omitempty"`
	To       string `json:"to,omitempty"`       // 待呼叫号码
	Userdata string `json:"userdata,omitempty"` // 只允许填写数字字符
	State    string `json:"state,omitempty"`    // 状态， 数字按键
	Message  string `json:"message,omitempty"`
	Sn       string `json:"sn,omitempty"` // 包裹编号
	Channel  string `json:"channel,omitempty"`
}

// 入mq 的方法  ------------------------

// 短信接口调用
func SendValidSn(mob, sn string) {

	tmpjson := "{\"cmd\":\"SendValidCode\", \"to\":\"" + mob + "\", \"code\":\"" + sn + "\" }"

	mq.Publish("sms.pusher.exchange", "direct", "sms.key", tmpjson, true)
	log.Println("status push to mq:", tmpjson)
}

//  反馈 管理员的操作指令信息
// cmd = WxTemplateFeedback   模板审核拒绝
// cmd = ...
//func FeedBackAdminOp(cmd, uuid, st, msg string) {
//	var s SysFeedback
//	s.Cmd = cmd
//	s.State = st
//	s.Message = msg
//	s.Uuid = uuid
//
//	buf := bytes.NewBuffer(make([]byte, 0, 16<<10))
//	buf.Reset()
//	json.NewEncoder(buf).Encode(s)
//	tmpjson := buf.String()
//
//	// 通用反馈队列
//	mq.Publish("state.pusher.exchange", "direct", "state.key", tmpjson, true)
//	log.Println("status push to mq:", tmpjson)
//}

//  发送 反馈消息
func SendFeedback(cmd, to, userdata, stateid, msg, uuid, openid string) {

	var s SysFeedback
	s.Cmd = cmd
	s.To = to
	s.Userdata = userdata
	s.State = stateid
	s.Message = msg
	s.Uuid = uuid
	s.Openid = openid
	s.Channel = "7" // 容联1，cmpp2 3， 红树4， 筑望5， 华为6，微信7

	buf := bytes.NewBuffer(make([]byte, 0, 16<<10))
	buf.Reset()
	json.NewEncoder(buf).Encode(s)
	tmpjson := buf.String()

	// 入单独的队列

	mq.Publish("wechat.out.exchange", "direct", "out.key", tmpjson, true)

	if stateid != "404" { // 太多了 不输出
		log.Println("status push to mq:", tmpjson)
	}

	return
}

//  -------------   客服消息 微信接口 ------------

// 给用户发送 普通文本消息  客服消息接口
func SendMessageText(wxid, openid, content string) {

	obj := custom.NewText(wxid, openid, content, "")

	buf := bytes.NewBuffer(make([]byte, 0, 16<<10))
	buf.Reset()
	json.NewEncoder(buf).Encode(obj)
	tmpjson := buf.String()

	go custom.Send(tmpjson, g.GetWechatAccessToken(wxid))
}

//发送客服消息 图文消息
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

	go custom.Send(tmpjson, g.GetWechatAccessToken(wxid))
}

func SendMessagePic(wxid, openid, title, desc, url, pic string) {
}

/*
// 给所有管理员发消息
func SendMessageTextToAdmins(wxid, content string) {
	//for _, c = range g.Config()
	for _, c := range g.Config().Admins {
		openid := c.Openid

		SendMessageText(wxid, openid, content)
	}
}

// 模板审核消息，推送给管理员
func SendCheckDataToAdmin(wxid string, in string) {

	var t SysImJson
	e := json.Unmarshal([]byte(in), &t)
	if e != nil {
		log.Println("json 解析失败: %s", e)
		return
	}

	SendMessageTextToAdmins(wxid, "内容: "+t.Note+"\n"+
		"编号: "+t.Uuid+"\n"+
		t.CompanyName+" "+
		t.Name+" "+
		t.Userdata+"\n\n"+
		"拒绝请回复 no "+t.Uuid+" 拒绝原因（可选）")
	return
}
*/

//发送微信账户的地理位置，获取周边的快递员信息
func sendGetCourierInfo(locationx, locationy, wxid, openid string) {
	open.ReportLocation(locationx, locationy, wxid, openid)
}

// ------------ 微信模板消息 ---------------

// 发送 取件消息 给 最终用户
// 发送消息 有2个入口， 一个是 mq ， 一个 开放平台的接口调用
// 不同的消息入口，状态返回方式是不同的， 所以用一个 参数来标示， 消息的来源  msgfrom  , 默认来源于mq
func SendSmsNotify(wxid, uuid string, in, strurl, msgfrom string) {
	var t SysImJson
	e := json.Unmarshal([]byte(in), &t)
	if e != nil {
		log.Println("json 解析失败: %s", e)
		return
	}

	var tmsg template.TemplateMessage

	if t.Openid == "" {
		tmsg.ToUser = GetOpenidFromMobile(t.WxId, t.To)
	} else {
		tmsg.ToUser = t.Openid
	}

	if tmsg.ToUser == "" {
		SendFeedback("WechatFeedback", t.To, t.Userdata, "404", "微信用户不存在", uuid, "")
		return // 该用户无法推送消息
	}

	tmsg.TemplateId = "IaH3G5ayQ1jU6mVrJLL1X9813vCeKnS2lVrNTJZhJjs"
	if strurl == "" && t.Sid == "" {
		tmsg.URL = "http://wechat2.shenbianvip.com/?openid=" + tmsg.ToUser + "&uuid=" + uuid //  模板消息 地址
	} else if t.Sid != "" {
		tmsg.URL = "http://ylb.im/" + t.Sid
	} else {
		tmsg.URL = strurl
	}

	feed := true
	if t.History == "1" {
		tmsg.Data.First.Value = "【云喇叭】查询结果！"
		tmsg.Data.Keyword1.Value = t.HistoryAt
		feed = false
	} else {
		tmsg.Data.First.Value = "【云喇叭】通知您快递到了！"
		tmsg.Data.Keyword1.Value = time.Now().Format("2006-01-02 15:04")
	}

	tmsg.Data.Remark.Value = "快递员电话：" + t.Name + " " + t.Userdata
	if strings.Index(tmsg.Data.Remark.Value, "详情") == -1 {
		tmsg.Data.Remark.Value += g.Config().AdMsg
	}

	tmsg.Data.Keyword2.Value = t.Sn
	tmsg.Data.Keyword3.Value = t.CompanyName
	tmsg.Data.Keyword4.Value = t.Note
	tmsg.Data.Keyword5.Value = ""

	tmsg.Data.Keyword1.Color = "#173177"
	tmsg.Data.Keyword2.Color = "#173177"
	tmsg.Data.Keyword3.Color = "#173177"
	tmsg.Data.Keyword4.Color = "#173177"
	tmsg.Data.Keyword5.Color = "#173177"

	buf := bytes.NewBuffer(make([]byte, 0, 16<<10))
	buf.Reset()
	json.NewEncoder(buf).Encode(tmsg)
	tmpjson := buf.String()

	// json 格式变化
	WechatSendTemplate(wxid, uuid, tmsg.ToUser, t.To, t.Userdata, tmpjson, msgfrom, feed) // feed 是否需要上报状态
}
//发送模板消息（自动检测物流信息，有信息时发送）
func SendSmsNotifies(openid string, access_token string, in, strurl string) {
	var tmsg template.TemplateMessage
	var t PMessage
	e := json.Unmarshal([]byte(in), &t)
	if e != nil {
		log.Println("json 解析失败: %s", e)
		return
	}
	tmsg.ToUser = openid
	tmsg.TemplateId = "GPbyjw2jcgPclh9s1AwLb0NUzzzNQKQIuN-0zlebYkc"
	tmsg.URL = strurl
	tmsg.Data.First.Value = "【云喇叭】通知您快递有新消息"
	tmsg.Data.Keyword1.Color = "#173177"
	tmsg.Data.Keyword2.Color = "#173177"
	tmsg.Data.Keyword3.Color = "#173177"
	tmsg.Data.Keyword1.Value = t.NUM
	tmsg.Data.Keyword2.Value = t.Msg
	tmsg.Data.Keyword3.Value = t.TIME
	token := access_token
	incompleteURL := "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=" + token
	req := httplib.Post(incompleteURL)

	buf := bytes.NewBuffer(make([]byte, 0, 16<<10))
	buf.Reset()
	json.NewEncoder(buf).Encode(tmsg)
	tmpjson := buf.String()
	log.Println(tmpjson)
	req.Body(tmpjson)
	resp, err := req.String()
	if err != nil {
		log.Println("[ERROR]", err)
	}
	log.Println(resp)
}
