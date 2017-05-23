package http

import (
	"g"
	"log"
	"net/http"
	//"errors"
	"net/url"
	"strings"
	"mp/util"
	"mp/user/oauth2"
	//"crypto/md5"
	//"encoding/hex"
	"github.com/garyburd/redigo/redis"
	//"encoding/xml"
	"encoding/json"
	//"encoding/base64"
	redispool "redis"
	"bytes"
	"io/ioutil"
	"math"
	"model"
	"open"
	"strconv"
	"time"
)

// 接收 YLB 开放平台回调
func configYLBOpenRoutes() {

	//接收 open 平台过来的  通知发送请求
	http.HandleFunc("/ylb-open/", func(w http.ResponseWriter, req *http.Request) {

		wxid := strings.Trim(req.URL.Path, "/ylb-open/")
		log.Println("Path -->", req.URL.Path, "wxid", wxid) //

		wxcfg := g.GetWechatConfig(wxid) // 通过微信id 获取 对接的配置信息
		if wxcfg == nil {
			log.Println("[Warn] wecat config not find", wxid)
			w.WriteHeader(400)
			return
		}

		queryValues, err := url.ParseQuery(req.URL.RawQuery)
		//log.Println("ParseQuery", queryValues)
		if err != nil {
			log.Println("[ERROR] URL.RawQuery", err)
			w.WriteHeader(400)
			return
		}

		nonce := queryValues.Get("nonce")
		timestamp := queryValues.Get("timestamp")
		signature := strings.ToLower(queryValues.Get("signature"))
		appid := queryValues.Get("appid")

		if appid != g.Config().Open.AppId {
			log.Println("[WARN] err appid", appid, ", want to", g.Config().Open.AppId)
			w.WriteHeader(400)
			return
		}

		ts := time.Now().Unix()
		tsf, _ := strconv.ParseFloat(timestamp, 64)
		//log.Println(ts, timestamp)

		if math.Abs(float64(ts)-tsf) > 600 {
			log.Println("[WARN] timestamp ts out of range 600 seconds", timestamp)
			w.WriteHeader(400)
			return
		}

		if signature != open.Sign(g.Config().Open.AppToken, timestamp, nonce, appid) {
			log.Println("[WARN] signature wrong", signature, open.Sign(g.Config().Open.AppToken, timestamp, nonce, appid))
			//w.WriteHeader(400)

			RenderJson(w, open.RespMsg{ErrCode: 101, ErrMsg: "签名不正确"})
			return
		}

		// 参数验证都ok 了 获取 json 数据

		//
		var msgjson open.Message
		if err := json.NewDecoder(req.Body).Decode(&msgjson); err != nil {
			log.Println("[Warn] json body", err)

			RenderJson(w, open.RespMsg{ErrCode: 102, ErrMsg: "json格式不正确"})
			return
		}

		// json 数据不正确
		if msgjson.Uuid == "" || msgjson.To == "" || msgjson.Msg == "" {
			log.Println("[Warn] json ", msgjson)

			RenderJson(w, open.RespMsg{ErrCode: 102, ErrMsg: "json中的数据不正确，uuid to msg不能为空"})
			return
		}

		log.Println("receive json data from ylb-open", msgjson)
		openid := model.GetOpenidFromMobile(wxcfg.AppId, msgjson.To) // 手机号码变换为 微信下的 openid

		if openid == "" {
			// 用户不存在
			go open.Report(msgjson.Uuid, openid, msgjson.To, "2", "用户不存在")
			RenderJson(w, open.RespMsg{ErrCode: 0, ErrMsg: "接口请求已经处理"})
			return
		}

		// 用户存在
		var d model.SysImJson
		d.Cmd = "SendSmsNotify"
		d.WxId = wxid
		d.Openid = openid
		d.Uuid = msgjson.Uuid
		d.Sn = msgjson.Sn
		d.To = msgjson.To
		d.Userdata = msgjson.Userdata
		d.CompanyName = msgjson.CompanyName
		d.Note = msgjson.Msg

		// 
		d.Note = strings.Replace(d.Note, "【云喇叭快递小管家】", "", -1)

		buf := bytes.NewBuffer(make([]byte, 0, 16<<10))
		buf.Reset()
		json.NewEncoder(buf).Encode(d)
		tmpjson := buf.String()
		log.Println(tmpjson)

		go model.SendSmsNotify(wxid, msgjson.Uuid, tmpjson, msgjson.Url, "open") // 模板通知 消息从 开放平台过来的

		RenderJson(w, open.RespMsg{ErrCode: 0, ErrMsg: "模板消息已经发送"})
		return
	})

	// 接收 open 平台 接口调用， 然后将快递员数据，推送给 真实的用户
	http.HandleFunc("/rpc/couriers", func(w http.ResponseWriter, req *http.Request) {

		queryValues, err := url.ParseQuery(req.URL.RawQuery)
		//log.Println("ParseQuery", queryValues)
		if err != nil {
			log.Println("[ERROR] URL.RawQuery", err)
			w.WriteHeader(400)
			return
		}

		nonce := queryValues.Get("nonce")
		timestamp := queryValues.Get("timestamp")
		signature := strings.ToLower(queryValues.Get("signature"))
		appid := queryValues.Get("appid")

		log.Println(nonce, " ", appid)

		if appid != g.Config().Location.AppId {
			log.Println("[WARN] err appid", appid, ", want to", g.Config().Location.AppId)
			w.WriteHeader(400)
			return
		}

		ts := time.Now().Unix()
		tsf, _ := strconv.ParseFloat(timestamp, 64)
		//log.Println(ts, timestamp)

		if math.Abs(float64(ts)-tsf) > 600 {
			log.Println("[WARN] timestamp ts out of range 600 seconds", timestamp)
			w.WriteHeader(400)
			return
		}

		if signature != open.Sign(g.Config().Location.AppToken, timestamp, nonce, appid) {
			log.Println("[WARN] signature wrong", signature, open.Sign(g.Config().Location.AppToken, timestamp, nonce, appid))
			w.WriteHeader(400)

			//RenderJson(w, rpc.RespMsg{ErrCode: 101, ErrMsg: "签名不正确"})
			return
		}

		// 参数验证都ok 了 获取 json 数据
		/*
			msg, _ := ioutil.ReadAll(req.Body)
			body := string(msg)
			log.Println(body) */

		var info open.CouriersMsg
		if err := json.NewDecoder(req.Body).Decode(&info); err != nil {
			log.Println("[Warn] json body", err)
			//RenderJson(w, rpc.RespMsg{ErrCode: 102, ErrMsg: "json格式不正确"})
			return
		}

		couriers := info.Couriers
		length := len(couriers)

		message := "您身边总共有" + strconv.Itoa(length) + "位快递小哥，分别为： \n"
		for _, courier := range couriers {
			message = message + courier.Company + "的" + courier.Name + " 手机号码：" + courier.Phone + "  \n"
		}

		log.Println(message)

		if length == 0 {
			message = "您身边暂时没有快递员"
		}

		//发送文本消息
		go model.SendMessageText(info.WxId, info.Openid, message)

	})

	http.HandleFunc("/rpc/getToken", func(w http.ResponseWriter, req *http.Request) {

		queryValues, err := url.ParseQuery(req.URL.RawQuery)
		//log.Println("ParseQuery", queryValues)
		if err != nil {
			log.Println("[ERROR] URL.RawQuery", err)
			w.WriteHeader(400)
			return
		}

		nonce := queryValues.Get("nonce")
		timestamp := queryValues.Get("timestamp")
		signature := strings.ToLower(queryValues.Get("signature"))
		appid := queryValues.Get("appid")

		//log.Println(nonce, " ", appid)

		if appid != g.Config().Location.AppId {
			log.Println("[WARN] err appid", appid, ", want to", g.Config().Location.AppId)
			w.WriteHeader(400)
			return
		}

		ts := time.Now().Unix()
		tsf, _ := strconv.ParseFloat(timestamp, 64)
		//log.Println(ts, timestamp)

		if math.Abs(float64(ts)-tsf) > 600 {
			log.Println("[WARN] timestamp ts out of range 600 seconds", timestamp)
			w.WriteHeader(400)
			return
		}

		if signature != open.Sign(g.Config().Location.AppToken, timestamp, nonce, appid) {
			log.Println("[WARN] signature wrong", signature, open.Sign(g.Config().Location.AppToken, timestamp, nonce, appid))
			w.WriteHeader(400)

			//RenderJson(w, open.RespMsg{ErrCode: 101, ErrMsg: "签名不正确"})
			return
		}

		//获取到微信ID
		//log.Println("get token info from wx server!")
		wxid := queryValues.Get("wxid")

		RenderJson(w, open.TokenMsg{ACCESS_TOKEN: g.GetWechatAccessToken(wxid), JS_TOKEN: g.GetJsApiTicket(wxid)})

	})

	http.HandleFunc("/rpc/getTelByOpenid", func(w http.ResponseWriter, req *http.Request) {

		queryValues, err := url.ParseQuery(req.URL.RawQuery)
		//log.Println("ParseQuery", queryValues)
		if err != nil {
			log.Println("[ERROR] URL.RawQuery", err)
			w.WriteHeader(400)
			return
		}

		nonce := queryValues.Get("nonce")
		timestamp := queryValues.Get("timestamp")
		signature := strings.ToLower(queryValues.Get("signature"))
		appid := queryValues.Get("appid")

		//log.Println(nonce, " ", appid)

		if appid != g.Config().Location.AppId {
			log.Println("[WARN] err appid", appid, ", want to", g.Config().Location.AppId)
			w.WriteHeader(400)
			return
		}

		ts := time.Now().Unix()
		tsf, _ := strconv.ParseFloat(timestamp, 64)
		//log.Println(ts, timestamp)

		if math.Abs(float64(ts)-tsf) > 600 {
			log.Println("[WARN] timestamp ts out of range 600 seconds", timestamp)
			w.WriteHeader(400)
			return
		}

		if signature != open.Sign(g.Config().Location.AppToken, timestamp, nonce, appid) {
			log.Println("[WARN] signature wrong", signature, open.Sign(g.Config().Location.AppToken, timestamp, nonce, appid))
			w.WriteHeader(400)

			//RenderJson(w, open.RespMsg{ErrCode: 101, ErrMsg: "签名不正确"})
			return
		}

		//获取到用户的openid，用于从后端获取到手机号码
		//log.Println("get token info from wx server!")
		openid := queryValues.Get("openid")
		wxid := queryValues.Get("wxid")
		log.Println(wxid + "..." + openid)
		u := model.CreateUser(wxid, openid)

		tel := u.Mobile1 + "," + u.Mobile2 + "," + u.Mobile3

		RenderJson(w, open.UserMsg{TEL: tel, OPENID: openid})

	})

	http.HandleFunc("/rpc/sendCouriersInfo", func(w http.ResponseWriter, req *http.Request) {
		log.Println("--->/rpc/sendCouriersInfo")
		queryValues, err := url.ParseQuery(req.URL.RawQuery)
		//log.Println("ParseQuery", queryValues)
		if err != nil {
			log.Println("[ERROR] URL.RawQuery", err)
			w.WriteHeader(400)
			return
		}

		nonce := queryValues.Get("nonce")
		timestamp := queryValues.Get("timestamp")
		signature := strings.ToLower(queryValues.Get("signature"))
		appid := queryValues.Get("appid")

		log.Println(nonce, " ", appid)

		if appid != g.Config().Location.AppId {
			log.Println("[WARN] err appid", appid, ", want to", g.Config().Location.AppId)
			w.WriteHeader(400)
			return
		}

		ts := time.Now().Unix()
		tsf, _ := strconv.ParseFloat(timestamp, 64)
		//log.Println(ts, timestamp)

		if math.Abs(float64(ts)-tsf) > 600 {
			log.Println("[WARN] timestamp ts out of range 600 seconds", timestamp)
			w.WriteHeader(400)
			return
		}

		if signature != open.Sign(g.Config().Location.AppToken, timestamp, nonce, appid) {
			log.Println("[WARN] signature wrong", signature, open.Sign(g.Config().Location.AppToken, timestamp, nonce, appid))
			w.WriteHeader(400)

			//RenderJson(w, rpc.RespMsg{ErrCode: 101, ErrMsg: "签名不正确"})
			return
		}

		//获取到用户的openid，用于从后端获取到手机号码
		//log.Println("get token info from wx server!")
		openid := queryValues.Get("openid")
		wxid := queryValues.Get("wxid")
		info, _ := ioutil.ReadAll(req.Body)

		log.Println("openid :" + openid + " infos:" + string(info) + "  wxid:" + wxid)
		go model.SendMessageText(wxid, openid, string(info))

		RenderJson(w, open.RespMsg{ErrCode: 0, ErrMsg: "模板消息已经发送"})
		return

	})

	//接受其他页面的code，返回用户信息
	http.HandleFunc("/api/info", func(w http.ResponseWriter, r *http.Request) {
		fullurl := "http://" + r.Host + r.RequestURI
	    wxid := "gh_8ac8a8821eb9"
	    wxcfg := g.GetWechatConfig(wxid)
	    appid := wxcfg.AppId
	    AppSecret := wxcfg.AppSecret
		qvalues, _ := url.ParseRequestURI(fullurl) 
		queryValues, _ := url.ParseQuery(qvalues.RawQuery)
		r.ParseForm()
		code := queryValues.Get("code")
		nonce := queryValues.Get("nonce")
		timestamp := queryValues.Get("timestamp")
		signature := strings.ToLower(queryValues.Get("signature"))
		ts := time.Now().Unix()
		tsf, _ := strconv.ParseFloat(timestamp, 64)
		if math.Abs(float64(ts)-tsf) > 600 {
			log.Println("[WARN] timestamp ts out of range 600 seconds", timestamp)
			w.WriteHeader(400)
			return
		}
		if signature != open.Sign(g.Config().Location.AppToken, timestamp, nonce, g.Config().Location.AppId) {
			log.Println("[WARN] signature wrong", signature, open.Sign(g.Config().Location.AppToken, timestamp, nonce, appid))
			w.WriteHeader(400)

			RenderMsgJson(w,"签名错误！！！")
			return
		}
		openid, token:= util.GetAccessTokenFromCode(appid, AppSecret, code)
		rc := redispool.ConnPool.Get()
		defer rc.Close()
		find, _ := redis.Bool(rc.Do("EXISTS", openid))
		var data open.UserMsgs
		if find == true {
			log.Println("已经是老用户了")
			smap, _ := redis.StringMap(rc.Do("HGETALL", openid))
			nm := smap["nickname"]
			if nm == ""{
				userinfo, err := oauth2.GetUserInfo(token, openid, "zh_CN")
				if err != nil {
					log.Println("[ERROR] userinfo", err)
				}
				data.SUB = smap["sub"]
				data.TEL = smap["mobile1"]
				data.NICKNAME = userinfo.Nickname
				data.HEADIMG = userinfo.HeadImageURL
				data.OPENID = openid
				log.Println("getusr from wechat"+openid)
				RenderJson(w, data)
				return
			}else{
				log.Println(smap)
				data.SUB = smap["sub"]
				data.TEL = smap["mobile1"]
				data.NICKNAME = smap["nickname"]
				data.HEADIMG = smap["imgurl"]
				data.OPENID = openid
				RenderJson(w, data)
				return
			}

		}else if openid == "" {
			log.Println("code invalid")
			RenderMsgJson(w,"code invalid")
			return
		}else{
			log.Println(openid)
			log.Println("new user")
			data.SUB = "0"
			data.OPENID = openid
			RenderJson(w, data)
			return
		}
			

	})
}
