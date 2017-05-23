package open

import (
//"encoding/xml"
//"encoding/json"
)

// 统一的 短信 语音 通知格式
type Message struct {
	Cmd         string `json:"cmd,omitempty"`
	WxId        string `json:"wxid,omitempty"`
	Openid      string `json:"openid,omitempty"`
	Uuid        string `json:"uuid,omitempty"`
	Sn          string `json:"sn,omitempty"`       // 包裹的自定义编号
	To          string `json:"to,omitempty"`       // 待呼叫号码
	Userdata    string `json:"userdata,omitempty"` // 只允许填写数字字符
	Name        string `json:"name,omitempty"`     // 快递员名字
	CompanyName string `json:"companyname,omitempty"`

	Msg        string `json:"msg,omitempty"` // 编号 + 时间 + 地点。 用户自定义的内容
	Noticetime string `json:"noticetime,omitempty"`
	Url        string `json:"url,omitempty"`
}

type RespMsg struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}


type ReportMsg struct {
	Cmd  string `json:"cmd,omitempty"`
	WxId string `json:"wxid,omitempty"`

	Openid string `json:"openid,omitempty"`
	Uuid   string `json:"uuid,omitempty"`
	To     string `json:"to,omitempty"` // 待呼叫号码

	Resulttype string `json:"resulttype,omitempty"` //
	Resultmsg  string `json:"resultmsg,omitempty"`  //
}

type TokenMsg struct {
	ACCESS_TOKEN string `json:"accessToken"`
	JS_TOKEN     string `json:"jsToken"`
}

type UserMsg struct {
	TEL    string `json:"tel"`
	OPENID string `json:"openid"`
}

//通过code请求用户信息
type UserMsgs struct {
	TEL      string `json:"tel"`
	OPENID   string `json:"openid"`
	SUB	     string `json:"sub"`
	HEADIMG  string `json:"headimg"`
	NICKNAME string `json:"nickname"`
}

type LocationMsg struct {
	Cmd  string `json:"cmd,omitempty"`
	WxId string `json:"wxid,omitempty"`

	Openid    string `json:"openid,omitempty"`
	LocationX string `json:"locationx,omitempty"`
	LocationY string `json:"locationy,omitempty"` //
	Msg       string `json:"resultmsg,omitempty"` //
}

type CouriersMsg struct {
	Cmd      string `json:"cmd,omitempty"`
	WxId     string `json:"wxid,omitempty"`
	Openid   string `json:"openid,omitempty"`
	Couriers []struct {
		Company   string `json:"expressCompanyName,omitempty"`
		Branch    string `json:"expressBranch,omitempty"`
		Name      string `json:"name,omitempty"`
		Phone     string `json:"phone,omitempty"`
		Id        int    `json:"id,omitempty"`
		CompanyId int    `json:"expressCompanyId,omitempty"`
		BranchId  int    `json:"expressBranchId,omitempty"`
		Sex       string `json:"sex,omitempty"`
	} `json:"courierInfo"`
}