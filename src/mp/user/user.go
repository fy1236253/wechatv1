package user

import (
	
)

const (
	Language_zh_CN = "zh_CN" // 简体中文
	Language_zh_TW = "zh_TW" // 繁体中文
	Language_en    = "en"    // 英文
)

const (
	SexUnknown = 0 // 未知
	SexMale    = 1 // 男性
	SexFemale  = 2 // 女性
)

type UserInfo struct {
	Cmd          string `json:"cmd"`       // 扩展
	EventKey     string `json:"eventkey"`  // 扩展
	Mobile       string `json:"mobile"`    // 扩展
	IsSubscriber int    `json:"subscribe"` // 用户是否订阅该公众号标识, 值为0时, 代表此用户没有关注该公众号, 拉取不到其余信息
	OpenId       string `json:"openid"`    // 用户的标识, 对当前公众号唯一
	Nickname     string `json:"nickname"`  // 用户的昵称
	Sex          int    `json:"sex"`       // 用户的性别, 值为1时是男性, 值为2时是女性, 值为0时是未知
	Language     string `json:"language"`  // 用户的语言, zh_CN, zh_TW, en
	City         string `json:"city"`      // 用户所在城市
	Province     string `json:"province"`  // 用户所在省份
	Country      string `json:"country"`   // 用户所在国家

	// 用户头像, 最后一个数值代表正方形头像大小(有0, 46, 64, 96, 132数值可选, 0代表640*640正方形头像),
	// 用户没有头像时该项为空
	HeadImageURL string `json:"headimgurl"`

	// 用户关注时间, 为时间戳. 如果用户曾多次关注, 则取最后关注时间
	SubscribeTime int64 `json:"subscribe_time"`

	// 只有在用户将公众号绑定到微信开放平台帐号后, 才会出现该字段.
	UnionId string `json:"unionid"`

	Remark  string `json:"remark"`  // 公众号运营者对粉丝的备注, 公众号运营者可在微信公众平台用户管理界面对粉丝添加备注
	GroupId int64  `json:"groupid"` // 用户所在的分组ID
}


// 获取关注者列表返回的数据结构
type UserListResult struct {
	TotalCount int `json:"total"` // 关注该公众账号的总用户数
	GotCount   int `json:"count"` // 拉取的 OPENID 个数, 最大值为10000

	Data struct {
		OpenIdList []string `json:"openid,omitempty"`
	} `json:"data"` // 列表数据, OPENID 的列表

	// 拉取列表的后一个用户的OPENID, 如果 next_openid == "" 则表示没有了用户数据
	NextOpenId string `json:"next_openid"`
}


