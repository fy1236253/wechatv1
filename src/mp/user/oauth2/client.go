package oauth2

import (
	"encoding/json"
	"errors"
	"github.com/toolkits/net/httplib"
	"log"
	"mp"
	"time"
)

// 获取用户信息(需scope为 snsapi_userinfo).
//  lang 可能的取值是 zh_CN, zh_TW, en, 如果留空 "" 则默认为 zh_CN.
func GetUserInfo(token, openid, lang string) (info *UserInfo, err error) {
	switch lang {
	case "":
		lang = Language_zh_CN
	case Language_zh_CN, Language_zh_TW, Language_en:
	default:
		lang = Language_zh_CN
	}

	u := "https://api.weixin.qq.com/sns/userinfo?access_token=" + token + "&openid=" + openid + "&lang=" + lang
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
