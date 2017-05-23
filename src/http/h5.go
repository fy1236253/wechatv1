package http

import (
	"g"
	"log"
	"net/http"
	//"errors"
	//"crypto/md5"
	//"encoding/hex"
	//"github.com/garyburd/redigo/redis"
	//"encoding/json"
	//redispool "redis"
	"github.com/toolkits/file"
	"html/template"
	"math/rand"
	"model"
	"mp/user/oauth2"
	"mp/util"
	"net/url"
	"path/filepath"
	"strconv"
	"time"
	//"github.com/astaxie/beego/session"

	"github.com/dchest/captcha"
)

func getuser(w http.ResponseWriter, r *http.Request) *model.User {
	fullurl := "http://" + r.Host + r.RequestURI
    wxid := "gh_8ac8a8821eb9"
    wxcfg := g.GetWechatConfig(wxid)
    appid := wxcfg.AppId
	//wxid := "gh_8ac8a8821eb9"
	//appid := "wxacc105428fe41835" // 云喇叭服务号

	// 参数检查
	queryValues, err := url.ParseQuery(r.URL.RawQuery)
	log.Println("ParseQuery", queryValues)
	if err != nil {
		log.Println("[ERROR] URL.RawQuery", err)
		w.WriteHeader(400)
		return nil
	}

	// 从 session 中获取用户的 openid
	sess, _ := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)
	if sess.Get("openid") == nil {
		sess.Set("openid", "")
	}
	openid := sess.Get("openid").(string)

	// session 不存在
	if openid == "" {
		//oauth 跳转 ， 页面授权获取用户基本信息
		code := queryValues.Get("code") //  摇一摇入口 code 有效
		state := queryValues.Get("state")
		if code == "" && state == "" {
			addr := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=" + appid + "&redirect_uri=" + url.QueryEscape(fullurl) + "&response_type=code&scope=snsapi_userinfo&state=1#wechat_redirect"
			log.Println("http.Redirect", addr)
			http.Redirect(w, r, addr, 302)
			return nil
		}

		// 获取用户信息
		c := g.GetWechatConfig(wxid)
		var token string
		openid, token = util.GetAccessTokenFromCode(c.AppId, c.AppSecret, code)

		if openid == "" {
			return nil
		}

		u := model.CreateUser(wxid, openid)
		if u != nil && u.NickName == "" { // 用户数据本地没有
			userinfo, err := oauth2.GetUserInfo(token, openid, "zh_CN")
			if err != nil {
				u = nil
			} else if userinfo != nil {
				log.Println("get userinfo from oauht2 api", userinfo)
				// 保存 userinfo 信息到 redis
				u.NickName = userinfo.Nickname
				u.Sex = userinfo.Sex
				u.ImgUrl = userinfo.HeadImageURL
				u.Save()
			}
		}

		sess.Set("openid", openid)
	}

	u := model.CreateUser(wxid, openid)
	log.Println(u, wxid, openid)
	log.Println("get openid from session", sess.Get("openid").(string), u.Mobile1)

	return u
}

func configH5Routes() {
	//发快递页面
	http.HandleFunc("/h5/send", func(w http.ResponseWriter, r *http.Request) {
		fullurl := "http://" + r.Host + r.RequestURI
	    wxid := "gh_8ac8a8821eb9"
	    wxcfg := g.GetWechatConfig(wxid)
	    appid := wxcfg.AppId
		sess, _ := globalSessions.SessionStart(w, r)
		defer sess.SessionRelease(w)

		var f string // 模板文件路径
		f = filepath.Join(g.Root, "/public", "h5-send.html")
		if !file.IsExist(f) {
			log.Println("not find", f)
			http.NotFound(w, r)
			return
		}

		// 基本参数设置
		rand.Seed(time.Now().UnixNano())
		nonce := strconv.Itoa(rand.Intn(9999999999))
		ts := time.Now().Unix()
		sign := util.WXConfigSign(g.GetJsApiTicket(wxid), nonce, strconv.FormatInt(ts, 10), fullurl)
		data := struct {
			//Couriers 	string
			AppId string
			Ts    int64
			Nonce string
			Sign  string
		}{
			AppId: appid,
			Ts:    ts,
			Nonce: nonce,
			Sign:  sign,
		}

		t, err := template.ParseFiles(f)
		err = t.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
		return

	})
//发快递处理模块
	http.HandleFunc("/h5/search", func(w http.ResponseWriter, r *http.Request) {
		var f string // 模板文件路径
		f = filepath.Join(g.Root, "/public", "h5-search.html")
		if !file.IsExist(f) {
			log.Println("not find", f)
			http.NotFound(w, r)
			return
		}
		queryValues, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			log.Println("[ERROR] URL.RawQuery", err)
		}
		r.ParseForm()
		num := queryValues.Get("num")
		openid := queryValues.Get("openid")		
		data := struct {
			NUM     string
			COMPANY string
			AMOUT	int
			Item  []*model.ExpressData
		}{}
		data.Item = model.GetCouryInfo(num,openid)
		data.AMOUT = len(data.Item)
		data.COMPANY = model.GetCompany(num)
		data.NUM = num
		//log.Println(data)
		t, err := template.ParseFiles(f)
		err = t.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
		return
	})

	//发快递处理模块
	http.HandleFunc("/h5/couriers", func(w http.ResponseWriter, r *http.Request) {
		var f string // 模板文件路径
		f = filepath.Join(g.Root, "/public", "h5-couriers.html")
		if !file.IsExist(f) {
			log.Println("not find", f)
			http.NotFound(w, r)
			return
		}

		fullurl := "http://" + r.Host + r.RequestURI
		log.Println(fullurl)

		qvalues, err := url.ParseRequestURI(fullurl)
		if err != nil {
			log.Println("[ERROR] URL.RawQuery", err)
			w.WriteHeader(400)
			return
		}
		queryValues, err := url.ParseQuery(qvalues.RawQuery)
		if err != nil {
			log.Println("[ERROR] URL.RawQuery", err)
			w.WriteHeader(400)
			return
		}
		r.ParseForm()
		log.Println("ParseQuery", queryValues)
		var longitude, latitude string
		longitude = queryValues.Get("lng")
		latitude = queryValues.Get("lat")
		//默认是通过微信接口location
		if longitude == "" {
			longitude = queryValues.Get("wxlng")
			latitude = queryValues.Get("wxlat")
			log.Println(latitude + " " + longitude)
		}
		// tel := u.Mobile1
		// log.Println(tel)
		log.Println(latitude + " " + longitude)
		lon, _ := strconv.ParseFloat(longitude, 64)
		lat, _ := strconv.ParseFloat(latitude, 64)
		geo := model.Encode(lat, lon)
		log.Println(geo)
		data := struct {
			Count int
			Item  []*model.Courier
		}{}
		data.Item = model.GetCouries(geo)
		count := len(data.Item)
		//log.Println(count)
		data.Count = count
		//log.Println(data)
		t, err := template.ParseFiles(f)
		err = t.Execute(w, data)
		if err != nil {
			log.Println(err)
		}

		return
	})

	// 最简单的 H5 ，手机号码绑定页面
	http.HandleFunc("/h5/bind", func(w http.ResponseWriter, r *http.Request) {

		// 获取用户信息
		log.Println("bind.getuser before <<<<<<<<")
		u := getuser(w, r) // oauth2 获取用户信息
		log.Println("bind.getuser after  >>>>>>>>", u)

		if u == nil {
			//w.WriteHeader(400)
			return
		}

		//
		sess, _ := globalSessions.SessionStart(w, r)
		defer sess.SessionRelease(w)

		// 手机号码已经有了， 直接调回活动 页面
		if u.Mobile1 != "" {
			id := "2" // 返回 活动2 默认
			if sess.Get("action-id") != nil {
				id = sess.Get("action-id").(string)
			}
			addr := "http://wechat2.shenbianvip.com/h5/action?id=" + id
			http.Redirect(w, r, addr, 302)
			return
		}

		//
		queryValues, err := url.ParseQuery(r.URL.RawQuery)
		log.Println("ParseQuery", queryValues)
		if err != nil {
			log.Println("[ERROR] URL.RawQuery", err)
			w.WriteHeader(400)
			return
		}

		x := queryValues.Get("snkey")
		m := queryValues.Get("userPhoneNumber")

		//  post 提交 ， 判断验证码 是否ok
		if x != "" &&
			sess.Get("mobile") != nil && sess.Get("snkey") != nil &&
			sess.Get("snkey").(string) == x &&
			sess.Get("mobile").(string) == m {

			log.Println("验证码ok", sess.Get("snkey").(string), sess.Get("mobile").(string))
			sess.Set("snkey", "null-null") // 验证码只一次有效

			u.Mobile1 = m
			u.Save() // 保存 手机号码

			// 重定向到 活动页面
			addr := "http://wechat2.shenbianvip.com/h5/action?id=" + sess.Get("action-id").(string)
			http.Redirect(w, r, addr, 302)
			return

		} else {
			log.Println("show bind form")
		}

		//  show  form
		var f string // 模板文件路径

		f = filepath.Join(g.Root, "/public", "h5-bind.html")
		if !file.IsExist(f) {
			log.Println("not find", f)
			http.NotFound(w, r)
			return
		}

		captchaId := captcha.NewLen(4)
		sess.Set("captchaId", captchaId)

		data := struct {
			HaveBind  string
			CaptchaId string
		}{
			HaveBind:  "", // 空代表 未注册
			CaptchaId: captchaId,
		}

		t, err := template.ParseFiles(f)
		err = t.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	})

	// 所有活动的 代理入口 , 统一判断，参加活动前的条件， 是否需要绑定手机号码
	// id=0  关注送话费 http://wechat2.shenbianvip.com/h5/r
	// id=1  发快递   http://wechat2.shenbianvip.com/InnerService/api/oauthAction.action
	// id=2  一元抢   http://wechat2.shenbianvip.com/InnerService/api/oauthYiyuanAction.action
	// id=3  20元韩装 http://ylb.im/A
	http.HandleFunc("/h5/action", func(w http.ResponseWriter, r *http.Request) {

		// 参数检查
		queryValues, err := url.ParseQuery(r.URL.RawQuery)
		//log.Println("ParseQuery", queryValues)
		if err != nil {
			log.Println("[ERROR] URL.RawQuery", err)
			w.WriteHeader(400)
			return
		}

		// 获取用户信息
		log.Println("action.getuser before <<<<<<<<")
		u := getuser(w, r) // oauth2 获取用户信息
		log.Println("action.getuser after  >>>>>>>>", u)

		if u == nil {
			//w.WriteHeader(400)
			return
		}

		sess, _ := globalSessions.SessionStart(w, r)
		defer sess.SessionRelease(w)

		id := queryValues.Get("id")

		if id == "1" { //发快递
			sess.Set("action-id", id)
			addr := "http://wechat2.shenbianvip.com/InnerService/api/oauthAction.action"
			http.Redirect(w, r, addr, 302)
			return
		} else if id == "2" { //一元夺宝
			sess.Set("action-id", id)
			if u.Mobile1 != "" {
				addr := "http://wechat2.shenbianvip.com/InnerService/api/oauthYiyuanAction.action"
				http.Redirect(w, r, addr, 302)
				return
			} else {
				addr := "http://wechat2.shenbianvip.com/h5/bind"
				http.Redirect(w, r, addr, 302)
				return
			}
		} else if id == "3" { // 20元韩装
			sess.Set("action-id", id)
			addr := "http://ylb.im/A"
			http.Redirect(w, r, addr, 302)
			return
		} else if id == "0" { // 绑定手机号码赢话费  活动
			sess.Set("action-id", id)
			addr := "http://wechat2.shenbianvip.com/h5/r"
			http.Redirect(w, r, addr, 302)
			return
		}
		return
	})

	// 活动奖励页面  关注 送话费
	http.HandleFunc("/h5/r", func(w http.ResponseWriter, r *http.Request) {
		// 基本参数设置
		fullurl := "http://" + r.Host + r.RequestURI
		wxid := "gh_8ac8a8821eb9"
		appid := "wxacc105428fe41835" // 云喇叭服务号

		var f string // 模板文件路径

		f = filepath.Join(g.Root, "/public", "h5-r.html")
		if !file.IsExist(f) {
			log.Println("not find", f)
			http.NotFound(w, r)
			return
		}

		// 参数检查
		queryValues, err := url.ParseQuery(r.URL.RawQuery)
		log.Println("ParseQuery", queryValues)
		if err != nil {
			log.Println("[ERROR] URL.RawQuery", err)
			w.WriteHeader(400)
			return
		}

		// 从 session 中获取用户的 openid
		sess, _ := globalSessions.SessionStart(w, r)
		defer sess.SessionRelease(w)
		if sess.Get("openid") == nil {
			sess.Set("openid", "")
		}
		openid := sess.Get("openid").(string)

		// session 不存在
		if openid == "" {
			//oauth 跳转 ， 页面授权获取用户基本信息
			code := queryValues.Get("code") //  摇一摇入口 code 有效
			state := queryValues.Get("state")
			if code == "" && state == "" {
				addr := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=" + appid + "&redirect_uri=" + url.QueryEscape(fullurl) + "&response_type=code&scope=snsapi_userinfo&state=1#wechat_redirect"
				log.Println("http.Redirect", addr)
				http.Redirect(w, r, addr, 302)
				return
			}

			// 获取用户信息
			c := g.GetWechatConfig(wxid)
			var token string
			openid, token = util.GetAccessTokenFromCode(c.AppId, c.AppSecret, code)

			if openid == "" {
				return
			}

			u := model.CreateUser(wxid, openid)
			if u != nil && u.NickName == "" { // 用户数据本地没有
				userinfo, err := oauth2.GetUserInfo(token, openid, "zh_CN")
				if err != nil {
					u = nil
				} else if userinfo != nil {
					log.Println("get userinfo from oauht2 api", userinfo)
					// 保存 userinfo 信息到 redis
					u.NickName = userinfo.Nickname
					u.Sex = userinfo.Sex
					u.ImgUrl = userinfo.HeadImageURL
					u.Save()
				}
			}

			sess.Set("openid", openid)
		}

		u := model.CreateUser(wxid, openid)
		log.Println(u, wxid, openid)
		log.Println("get openid from session", sess.Get("openid").(string), u.Mobile1)

		rand.Seed(time.Now().UnixNano())
		nonce := strconv.Itoa(rand.Intn(99999999))
		ts := time.Now().Unix()
		captchaId := captcha.NewLen(4)

		sess.Set("captchaId", captchaId)

		data := struct {
			AppId string
			Ts    int64
			Nonce string
			Sign  string

			OpenId    string
			HaveReg   string
			CaptchaId string
		}{
			AppId: appid,
			Ts:    ts,
			Nonce: nonce,
			Sign:  util.WXConfigSign(g.GetJsApiTicket(wxid), nonce, strconv.FormatInt(ts, 10), fullurl),

			OpenId:    openid,
			HaveReg:   u.Mobile1, // 空代表 未注册
			CaptchaId: captchaId,
		}

		log.Println(data)

		t, err := template.ParseFiles(f)
		err = t.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
		//RenderDataJson(w, nil)
	})

}
