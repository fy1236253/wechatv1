package http

import (
	"g"
	"log"
	"net/http"
	"errors"
	"net/url"
	//"crypto/md5"
	//"encoding/hex"
	"github.com/garyburd/redigo/redis"
	
	"bytes"
	"encoding/json"
	redispool "redis"
	"mp/account"
	"strconv"
	"model"
	"sort"
	"crypto/md5"
	"crypto/sha1"
	"strings"
	"encoding/hex"
	"math/rand"
	"time"
	"math"

	"github.com/dchest/captcha" 
	"mq"
)


// api 接口 只对内网开放访问 
func configApiRoutes() {


	http.HandleFunc("/api/debug", func(w http.ResponseWriter, req *http.Request) {
		
		/*
		queryValues, _ := url.ParseQuery(req.URL.RawQuery)
		a := queryValues.Get("a")
		b := queryValues.Get("b")

		if captcha.VerifyString(a, b) {
			log.Println(a, b, "true")
		} else {
			log.Println(a, b, "false")
		}
		*/
		StdRender(w, g.VERSION, nil)
	})

	// 图像验证码接口 
	http.HandleFunc("/api/imgcode/", func(w http.ResponseWriter, req *http.Request) {
		//log.Println("imgcode", captcha.New())
		captcha.Server(captcha.StdWidth, captcha.StdHeight).ServeHTTP(w, req)
	})

	//快递员 主叫号码设置 接口 ，
	//get
	http.HandleFunc("/api/v1/logs/", func(w http.ResponseWriter, req *http.Request) {
		//model.SendCashBill("oyZW-w9b1cqaS-UncIENWIf_mdEo","13551243019",1)
		RenderDataJson(w, nil)
	})

	// ping
	http.HandleFunc("/api/version", func(w http.ResponseWriter, req *http.Request) {
		StdRender(w, g.VERSION, nil)
		return
	})

	http.HandleFunc("/api/ping", func(w http.ResponseWriter, req *http.Request) {
		StdRender(w, "pong", nil)
		return
	})

	// 发送 验证码 到指定手机号码, 需要图像验证码  
	http.HandleFunc("/api/v1/mobile/sn/send", func(w http.ResponseWriter, r *http.Request) {
		log.Println("---> /api/v1/mobile/sn/send ")
		//StdRender(w, "pong", nil)
		queryValues, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			log.Println("[ERROR] URL.RawQuery", err)
			w.WriteHeader(400)
			return
		}

		mob := queryValues.Get("m")
		x 	:= queryValues.Get("x") // 图形验证码 mm

		if len(mob) != 11 || mob[0:1] != "1" { 
			RenderMsgJson(w, "手机号码不正确！")
			return 
		}


		// sess.Set("captchaId", captchaId)
		sess, _ := globalSessions.SessionStart(w, r)
    	defer sess.SessionRelease(w)

    	if sess.Get("captchaId") == nil {
    		sess.Set("captchaId", "")
    	}
    	captchaId := sess.Get("captchaId").(string)

		if captcha.VerifyString(captchaId, x) {
			log.Println("imgcode is ok", x )
		} else {
			RenderMsgJson(w, "图形验证码错误！")
			return 
		}

 		var snkey string 

		rand.Seed(time.Now().UnixNano())
		snkey = strconv.Itoa(rand.Intn(99999)) // 随机4位扩展码

		sess.Set("snkey", snkey)
		sess.Set("mobile", mob)


		model.SendValidSn(mob, snkey)

		RenderDataJson(w, nil) // 

		return
	})

	// 验证手机号码 是否有效 
	http.HandleFunc("/api/v1/mobile/sn/check", func(w http.ResponseWriter, r *http.Request) {
		log.Println("---> /api/v1/mobile/sn/check ")
		//StdRender(w, "pong", nil)
		queryValues, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			log.Println("[ERROR] URL.RawQuery", err)
			w.WriteHeader(400)
			return
		}

		m 	:= queryValues.Get("m")
		x 	:= queryValues.Get("x")

		sess, _ := globalSessions.SessionStart(w, r)
    	defer sess.SessionRelease(w)

    	if sess.Get("mobile") != nil && sess.Get("snkey") != nil && sess.Get("snkey").(string) == x  && sess.Get("mobile").(string) == m  {
    		log.Println("验证码ok", sess.Get("snkey").(string), sess.Get("mobile").(string))
    		sess.Set("snkey", "null-null") // 验证码只一次有效 
    	} else {
    		log.Println("验证码不正确 ", m, x )
    		AutoRender(w, nil, errors.New("验证码不正确")) // 错误信息提示 
    		return 
    	}


    	// 用户数据 写入 redis 中 
		openid := "" 
		if sess.Get("openid") != nil {
			openid = sess.Get("openid").(string)
		}
		wxid := "gh_8ac8a8821eb9"
		u := model.CreateUser(wxid, openid)
		u.Mobile1 = m 
		u.Data = "new"  // 新用户
		u.Save()


		valid := 1234  // 手机号码临时关联 

		qr, e := account.CreateTemporaryQRCode(uint32(valid), 300, g.GetWechatAccessToken(wxid))
		if e == nil {
			data := map[string]string{}
			data["qrurl"] = "https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=" + url.QueryEscape(qr.Ticket)
			log.Println(data)
			RenderDataJson(w, data)
			return 
		}
		AutoRender(w, nil, e) // 错误信息提示 
		 

		return
	})


	// 只允许内网访问 ， 不对外 
	// 传递手机号码 返回一个 二维码路径 
	// 内网 访问地址 http://10.117.106.16:8022/   或  http://SMSD:8022 
	// t  t=0  临时二维码 ； t=1 永久二维码  ；t=2 永久二维码 
	// sence 场景id值；  临时二维码有效范围 1 - 4294967296；  永久二维码有效范围 1--100000； sence string类型 1到64长度
	// ttl   t=0 时 有效  默认 300秒 5分钟
	// 返回 json 数据 qrurl 为二维码路径
	// {"msg":"success","ts":"20160127173412","data":{"sence":"123","t":"0","qrurl":"https://%3D%3D"}}
	// http://SMSD:8022/api/v1/qrcode?t=2&sence=13551243019
	http.HandleFunc("/api/v1/qrcode", func(w http.ResponseWriter, r *http.Request) {
		log.Println("---> /api/v1/qrcode")

		queryValues, err := url.ParseQuery(r.URL.RawQuery)
		log.Println("ParseQuery", queryValues)
		if err != nil {
			log.Println("[ERROR] URL.RawQuery", err)
			w.WriteHeader(400)
			return
		}

		wxid := "gh_8ac8a8821eb9" 

		t 	:= queryValues.Get("t")
		ttl, _ := strconv.Atoi(queryValues.Get("ttl"))
		sence := queryValues.Get("sence")
		
		valid, _ := strconv.ParseInt(sence, 10, 64)

		if t == "0" {
			if valid < 0 || valid > 4294967296 {
				RenderMsgJson(w, "sence无效")
			 	return 
			}
		} else if t == "1" {
			if valid < 0 || valid > 100000 {
				RenderMsgJson(w, "sence无效")
				return 
			}
		}

		
		if ttl == 0 {
			ttl = 300 // ttl 一周有效 604800
		}

		data := map[string]string{}
		data["t"] = queryValues.Get("t")
		data["ttl"] = queryValues.Get("ttl")
		data["sence"] = queryValues.Get("sence")

		if t == "0" { // 临时 
			qr, e := account.CreateTemporaryQRCode(uint32(valid), ttl, g.GetWechatAccessToken(wxid))
			if e == nil {
				data["qrurl"] = "https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=" + url.QueryEscape(qr.Ticket)
				log.Println(data)
				RenderDataJson(w, data)
				return 
			}
			AutoRender(w, nil, e) // 错误信息提示 
			return 

		} else if t == "1" {
			qr, e := account.CreateQRCode(account.QR_LIMIT_SCENE, valid, ttl, g.GetWechatAccessToken(wxid))
			if e == nil {
				data["qrurl"] = "https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=" + url.QueryEscape(qr.Ticket)
				log.Println(data)
				RenderDataJson(w, data)
				return 
			}
			AutoRender(w, nil, e) // 错误信息提示 
			return 

		} else if t == "2" {
			qr, e := account.CreateQRCode(account.QR_LIMIT_STR_SCENE, valid, ttl, g.GetWechatAccessToken(wxid))
			if e == nil {
				data["qrurl"] = "https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=" + url.QueryEscape(qr.Ticket)
				log.Println(data)
				RenderDataJson(w, data)
				return 
			}
			AutoRender(w, nil, e) // 错误信息提示 
			return 

		}
		
	})


	// /api/v1/recharge/order-status?pid=wx20160519183512bef5b2d20f0117738565 & uid=
	http.HandleFunc("/api/v1/recharge/order-status", func(w http.ResponseWriter, r *http.Request) {
		log.Println("---> /api/v1/recharge/order-status ")

		queryValues, err := url.ParseQuery(r.URL.RawQuery)
		log.Println("ParseQuery", queryValues, err)  

		pid := queryValues.Get("pid")   
		uid := queryValues.Get("uid")  

		conn := redispool.ConnPool.Get()
		defer conn.Close()

		/*
		                                  预订单            支付ok     预充流量         流量到账
		conn.Do("HMSET", uid, "prepayid", result.PrepayId, "pay", 0, "precharge", 0, "charge", 0, "val", val)
		conn.Do("EXPIRE",uid, 7*24*3600)  // 保存一周
		*/
		find, _ := redis.Bool(conn.Do("EXISTS", uid))
		if find { // 订单有效 
			smap, _ := redis.StringMap(conn.Do("HGETALL", uid))
			log.Println("smap", smap)
			if smap["prepayid"] == pid && smap["pay"] == "0" && smap["precharge"] == "0" {

				payok := model.WeixinOrderStatus("gh_8ac8a8821eb9", uid)

				if payok {
					conn.Do("HMSET", uid, "pay", 1, "precharge", 1) // 支付成功  mq 命令ok 
			
					// 充值订单 入 mq 
					var o model.RechargeOrder 
					o.Cmd = "recharge"
					o.Uuid = uid 
					o.Phone = smap["to"]
					o.Value = smap["oldv"]
					o.Bits = strings.Replace(smap["tname"],"M","",-1)
					o.Type = "1"

					buf := bytes.NewBuffer(make([]byte, 0, 16<<10))
					buf.Reset()
					json.NewEncoder(buf).Encode(o)
					tmpjson := buf.String()

					mq.Publish("ka007.exchange", "direct", "rechange.key", tmpjson, true)
					log.Println("mq publish", tmpjson)

					RenderMsgJson(w,"success")
					return
				} else {
					log.Println("order status is false", pid, uid)
					RenderMsgJson(w,"该订单支付不成功")
					return
				}
					

			} else {
				log.Println("recharge ignore", pid, uid)
				RenderMsgJson(w,"该订单充值请求已经处理")
				return
			}
		}
		RenderMsgJson(w,"该订单数据不存在")
		return	
	})


	// 充流量 订单 确认
	http.HandleFunc("/api/v1/recharge/order", func(w http.ResponseWriter, r *http.Request) {
		log.Println("---> /api/v1/recharge/order ")

		queryValues, err := url.ParseQuery(r.URL.RawQuery)
		log.Println("ParseQuery", queryValues, err)  

		//  m:[13551243019] x:[0] i:[oyZW-w6uuFwN3l-XGVLq9XZHndDs] t:[10M] v:[2.7]

		m := queryValues.Get("m")  //  手机号码 
		i := queryValues.Get("i")  //  用户 openid 

		x := queryValues.Get("x")  //  套餐序号  
		t := queryValues.Get("t")  //  套餐名称 
		v := queryValues.Get("v")  //  实际价格 单位元 
		y := queryValues.Get("y")  //  原价格 

		// 套餐有效性判断 
		if validValBis(m, x, t, v, y) == false {
			log.Println("套餐选择异常",m,x,t,v, y)
			RenderMsgJson(w, "套餐选择异常" )
			return 
		}

		

		wxid := "gh_8ac8a8821eb9" 
		appid := "wxacc105428fe41835"
		rand.Seed(time.Now().UnixNano())
		nonce := strconv.Itoa(rand.Intn(9999999999))
		ts := strconv.FormatInt(time.Now().Unix(), 10)

		uid := time.Now().Format("20060102150405_") + m + "_" + t  //  最大 32 字符 
		// r.RemoteAddr  // 127.0.0.1:39825 
		// r.Header.Get("X-Forwarded-For")   118.112.58.109, 182.140.184.8

		// WeixinOrder(wxid, uid, openid, orderdesc, val, ip
		vf, _:= strconv.ParseFloat(v, 64)  //  订单价格 元 单位 
		succ, pid := model.WeixinOrder(wxid, uid, i, m, t, y, vf, GetClientIp(r) )

		if succ {
			pkg := "prepay_id="+pid

			// 签名的参数有appId, timeStamp, nonceStr, package, signType=SHA1
			data := map[string]string{}
			data["appId"] 		= appid
			data["timeStamp"] 	= ts
			data["nonceStr"] 	= nonce
			//data["package"] 	= pkg
			data["pkg"] 		= pkg // 别名 
			data["signType"]    = "MD5"
			data["paySign"] 	= paySign(data, g.Config().WeixinPay.Key)
			data["pid"]			= pid
			data["uid"]			= uid 
			
			log.Println(data) // 
			//appId:wxacc105428fe41835 timeStamp:1463625468 nonceStr:4231892229 package:prepay_id=wx201605191037484a1de7ca240974655753 
			// signType:SHA1 paySign:4E929D338E747CC8983543280860673A130BB282 pkg:prepay_id=wx201605191037484a1de7ca240974655753
			RenderDataJson(w, data) // 
		} else {
			log.Println("订单接口失败", pid)
			RenderMsgJson(w, pid)
		}
		
		return
	})

}



// 临时放在这，需要整理
//最后参与签名的参数有 appId, timeStamp, nonceStr, package, signType
func paySign(m map[string]string, key string) string {

	strs := sort.StringSlice{ "appId=" + m["appId"], 
		"timeStamp=" + m["timeStamp"],   //  这注意  timestamp  s 小写  
		"nonceStr=" + m["nonceStr"], 
		"package=" + m["pkg"], 
		"signType=" + m["signType"]}

	strs.Sort()

	strA := strings.Join(strs[:], "&")

	log.Println(strA)

	strB := strA + "&key=" + key 

	log.Println(strB)

	// appId=wxacc105428fe41835&nonceStr=1113819655    &package=prepay_id=wx20160519135525ae89e8e8720393892068&signType=SHA1&timeStamp=1463637324       &key=9scjWXbTc5yS3VYBQMKxvv8FIaBNFg99
	// appId="+appId        + "&nonceStr="+noncestr + "&package=prepay_id=wx2015041419450958e073ca4a0071648005&signType=MD5 &timeStamp=" + timestamp + "&key="+key

	if m["signType"] == "MD5" {
		md5Sum := md5.Sum([]byte(strB))
    	sig  := strings.ToUpper(hex.EncodeToString(md5Sum[:]))
    	return sig 
	} else {
		hashsum := sha1.Sum([]byte(strB))
		sig  := strings.ToUpper(hex.EncodeToString(hashsum[:]))
    	return sig 
	}
}



var (
	T1 = map[string]int{"10M": 3, "30M": 5, "100M": 10, "300M": 20, "500M": 30}
	T2 = map[string]int{"20M": 3, "50M": 6, "100M": 10, "200M": 15, "500M": 30}
	T3 = map[string]int{"5M": 1, "10M": 2, "30M": 5, "50M": 7, "100M": 10, "200M": 15, "500M": 30}
)

// 获取
// {10M 3} {30M 5} {100M 10} {300M 20}  {500M 30}  移动 
// {20M 3} {50M 6} {100M 10} {200M 15}  {500M 30}  联通 
// {5M  1} {10M 2} {30M 5} {50M 7} {100M 10} {200M 15} {500M 30}  电信 
func validValBis(to string, idx, name, realval, oldval string) bool {

	y, _ := strconv.Atoi(oldval)

	if g.MobileBelongTo(to) == g.ChinaMobile {
		if y != T1[name] {
			return false 
		}

		p   := float64(T1[name]) * 0.9  // 全部9折 
		r, _:= strconv.ParseFloat(realval, 64)
		if math.Abs(p-r) <= 0.01 {
			return true
		}
		log.Println(" yidong rechage type not find", to, name, realval)
	}

	if g.MobileBelongTo(to) == g.ChinaUnicom {
		if y != T2[name] {
			return false 
		}

		p := float64(T2[name]) * 0.9  // 全部9折 
		r,_ := strconv.ParseFloat(realval, 64)
		if math.Abs(p-r) <= 0.01 {
			return true
		}
		log.Println("liantong rechage type not find", to, name, realval)
	}

	if g.MobileBelongTo(to) == g.ChinaTelecom {
		if y != T3[name] {
			return false 
		}

		p := float64(T3[name]) * 0.9  // 全部9折 
		r,_ := strconv.ParseFloat(realval, 64)
		if math.Abs(p-r) <= 0.01 {
			return true
		}
		log.Println("dianxin rechage type not find", to, name, realval)
	}

	return false
}

