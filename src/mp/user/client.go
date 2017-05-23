package user

import (
	"encoding/json"
	"errors"
	"github.com/toolkits/net/httplib"
	"log"
	"mp"
	"time"
	"net/url"
)

// 获取用户信息
//  lang 可能的取值是 zh_CN, zh_TW, en, 如果留空 "" 则默认为 zh_CN.
func GetUserInfo(token, openid, lang string) (info *UserInfo, err error) {
	switch lang {
	case "":
		lang = Language_zh_CN
	case Language_zh_CN, Language_zh_TW, Language_en:
	default:
		lang = Language_zh_CN
	}

	u := "https://api.weixin.qq.com/cgi-bin/user/info?access_token=" + token + "&openid=" + openid + "&lang=" + lang
	req := httplib.Get(u).SetTimeout(3*time.Second, 1*time.Minute)
	resp, err := req.String()
	log.Println(resp)

	if err != nil {
		log.Println("[ERROR]", err)
		return nil, err
	}

	var result struct {
		mp.Error
		UserInfo
	}

	err = json.Unmarshal([]byte(resp), &result)
	if result.ErrCode != 0 {
		log.Println("[ERROR]", result)
		return nil, errors.New(result.ErrMsg)
	}

	return &result.UserInfo, nil
}



// 获取关注者列表.
//  NOTE:
//  1. 每次最多能获取 10000 个用户, 可以多次指定 NextOpenId 来获取以满足需求, 如果 NextOpenId == "" 则表示从头获取
//  2. 目前微信返回的数据并不包括 NextOpenId 本身, 是从 NextOpenId 下一个用户开始的, 和微信文档描述不一样!!!
func GetUserList(token, NextOpenId string) (rslt *UserListResult, err error) {
	

	var u string
	if NextOpenId == "" {
		u = "https://api.weixin.qq.com/cgi-bin/user/get?access_token=" + token 
	} else {
		u = "https://api.weixin.qq.com/cgi-bin/user/get?next_openid=" + url.QueryEscape(NextOpenId) + "&access_token=" + token 
	}

	req := httplib.Get(u).SetTimeout(3*time.Second, 1*time.Minute)
	resp, err := req.String()
	//log.Println(resp)


	var result struct {
		mp.Error
		UserListResult
	}

	err = json.Unmarshal([]byte(resp), &result)
	if result.ErrCode != 0 {
		log.Println("[ERROR]", result)
		return nil, errors.New(result.ErrMsg)
	}

	rslt = &result.UserListResult
	return
}



