package http

import (
	"g"
	"github.com/toolkits/file"
	"html/template"
	"log"
	"math/rand"
	"model"
	//"mp/account"
	"mp/shakearound"
	"mp/user/oauth2"
	"mp/util"
	"net/http"
	"net/url"
	"path/filepath"
	"proc"
	"strconv"
	"strings"
	"time"
)

func configW3Routes() {

	//  -----  充流量的实现     
	http.HandleFunc("/pay/wx/recharge", func(w http.ResponseWriter, r *http.Request) {


		if r.Method == "HEAD" {  // 如果仅仅查询头部信息 则直接返回 
			w.WriteHeader(200)
			return 
		}

		if r.Method == "GET" {
			var f string // 模板文件路径

			f = filepath.Join(g.Root, "/public", "recharge.html")
			if !file.IsExist(f) {
				log.Println("not find", f)
				http.NotFound(w, r)
				return
			}

			queryValues, err := url.ParseQuery(r.URL.RawQuery)
			log.Println("ParseQuery", queryValues)
			if err != nil {
				log.Println("[ERROR] URL.RawQuery", err)
				w.WriteHeader(400)
				return
			}

			// 基本参数设置
			fullurl := "http://" + r.Host + r.RequestURI
		    wxid := "gh_8ac8a8821eb9"
		    wxcfg := g.GetWechatConfig(wxid)
		    appid := wxcfg.AppId
			rand.Seed(time.Now().UnixNano())
			nonce := strconv.Itoa(rand.Intn(9999999999))
			ts := time.Now().Unix()

			log.Println(r.Method, fullurl)

			
			code := queryValues.Get("code") //  摇一摇入口 code 有效
			state := queryValues.Get("state")
			openid := ""
			phonebind := "" // 默认的绑定手机号码 

			// code 是空的 需要 重定向 https://open.weixin.qq.com/connect/oauth2/authorize?appid=APPID&redirect_uri=REDIRECT_URI&response_type=code&scope=SCOPE&state=STATE#wechat_redirect
			if code == "" && state == "" {
				addr := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=" + appid + "&redirect_uri=" + url.QueryEscape(fullurl) + "&response_type=code&scope=snsapi_base&state=1#wechat_redirect"
				//log.Println(addr)
				http.Redirect(w, r, addr, http.StatusFound)
				return
			} else {
				// code 一定存在
				c := g.GetWechatConfig(wxid)
				openid, _ = util.GetAccessTokenFromCode(c.AppId, c.AppSecret, code)
				log.Println("get openid", openid)
				if openid == "" {
					return 
				}

				u := model.CreateUser(wxid, openid)
				if u != nil {
					phonebind = u.Mobile1
				}
			}

			data := struct {
				AppId  string
				OpenId string 
				Ts     int64
				Nonce  string
				Sign   string
				PhoneBind string 
			}{
				AppId:  appid,
				OpenId: openid,
				Ts:     ts,
				Nonce:  nonce,
				Sign:   util.WXConfigSign(g.GetJsApiTicket(wxid), nonce, strconv.FormatInt(ts, 10), fullurl),
				PhoneBind: phonebind,
			}

			t, err := template.ParseFiles(f)
			err = t.Execute(w, data)
			if err != nil {
				log.Println(err)
			}
		}

		return
	})



	//  -----  快递物流信息查询    
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		fullurl := "http://" + r.Host + r.RequestURI
	    wxid := "gh_8ac8a8821eb9"
	    wxcfg := g.GetWechatConfig(wxid)
	    appid := wxcfg.AppId
	    AppSecret := wxcfg.AppSecret
		queryValues, err := url.ParseQuery(r.URL.RawQuery)
		code := queryValues.Get("code") 
		if code == "" {
			addr := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=" + appid + "&redirect_uri=" + url.QueryEscape(fullurl) + "&response_type=code&scope=snsapi_base&state=1#wechat_redirect"
			//log.Println("http.Redirect", addr)
			http.Redirect(w, r, addr, 302)
			return
		}
		//log.Println(code)
		openid, _:= util.GetAccessTokenFromCode(appid, AppSecret, code)
		//log.Println(openid)
		var f string // 模板文件路径

		f = filepath.Join(g.Root, "/public", "search.html")
		if !file.IsExist(f) {
			log.Println("not find", f)
			http.NotFound(w, r)
			return
		}

		// 基本参数设置

		rand.Seed(time.Now().UnixNano())
		nonce := strconv.Itoa(rand.Intn(9999999999))
		ts := time.Now().Unix()

		data := struct {
			OpenId string
			AppId  string
			Ts     int64
			Nonce  string
			Sign   string
		}{
			OpenId: openid,
			AppId:  appid,
			Ts:     ts,
			Nonce:  nonce,
			Sign:   util.WXConfigSign(g.GetJsApiTicket(wxid), nonce, strconv.FormatInt(ts, 10), fullurl),
		}

		t, err := template.ParseFiles(f)
		err = t.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	
		return
	})


	// --------------------------------  以下 代码 都要废弃的  --------------------------------



	// 摇一摇后的首页， ** 您好， 您最近2天内有 N个快递， 有您代取的快递有M个， 当前有 Q个人在排队，  我要排队
	http.HandleFunc("/wait-in-line", func(w http.ResponseWriter, r *http.Request) {

		var f string // 模板文件路径

		f = filepath.Join(g.Root, "/public", "wait-in-line.html")
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

		// 基本参数设置
		fullurl := "http://" + r.Host + r.RequestURI
		wxid := "gh_8ac8a8821eb9"
		appid := "wxacc105428fe41835" // 云喇叭服务号
		rand.Seed(time.Now().UnixNano())
		nonce := strconv.Itoa(rand.Intn(9999999999))
		ts := time.Now().Unix()

		var u *model.User

		data := struct {
			AppId  string
			Ts     int64
			Nonce  string
			Sign   string
			OpenId string
			DebugT string

			NickName     string
			Sex          string
			At           string
			DistanceInfo string
			Count        int
			Other        string
			Items        []*model.User
		}{
			AppId:  appid,
			Ts:     ts,
			Nonce:  nonce,
			Sign:   util.WXConfigSign(g.GetJsApiTicket(wxid), nonce, strconv.FormatInt(ts, 10), fullurl),
			Items:  []*model.User{},
			DebugT: "false",
		}

		//oauth 跳转 ， 页面授权获取用户基本信息
		code := queryValues.Get("code") //  摇一摇入口 code 有效
		state := queryValues.Get("state")

		// code 是空的 需要 重定向 https://open.weixin.qq.com/connect/oauth2/authorize?appid=APPID&redirect_uri=REDIRECT_URI&response_type=code&scope=SCOPE&state=STATE#wechat_redirect
		if code == "" && state == "" {
			addr := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=" + appid + "&redirect_uri=" + url.QueryEscape(fullurl) + "&response_type=code&scope=snsapi_userinfo&state=1#wechat_redirect"
			//log.Println(addr)
			http.Redirect(w, r, addr, 302)
			return

		} else {
			// code 一定存在
			c := g.GetWechatConfig(wxid)

			var token, openid string

			openid, token = util.GetAccessTokenFromCode(c.AppId, c.AppSecret, code)

			u = model.CreateUser(wxid, openid)
			if u.NickName == "" { // 用户数据本地没有

				userinfo, err := oauth2.GetUserInfo(token, openid, "zh_CN")
				if err != nil {
					//u = nil
				} else if userinfo != nil {
					log.Println("get userinfo from oauht2 api", userinfo)
					// 保存 userinfo 信息到 redis
					u.NickName = userinfo.Nickname
					u.Sex = userinfo.Sex
					u.ImgUrl = userinfo.HeadImageURL
					u.Save()
				}

			}

			if u.Sex == oauth2.SexMale {
				data.Sex = "先生"
			} else if u.Sex == oauth2.SexFemale {
				data.Sex = "女士"
			}
			data.NickName = u.NickName
		}

		if u != nil {
			pkgs := model.GetPackages(u, "")
			data.Count = len(pkgs)
		}

		ticket := queryValues.Get("ticket")
		devid := 0

		if ticket != "" {
			info, err := shakearound.GetShakeInfo(ticket, g.GetWechatAccessToken(wxid))
			//log.Println(info , err)
			if err != nil {
				// 获取设备信息  失败
			} else if info != nil && info.BeaconInfo.Distance > 0 {
				d := strconv.FormatInt(int64(info.BeaconInfo.Distance), 10)
				data.DistanceInfo = "距离取件地点还有" + d + "米"
			}

			// 根据设备信息  确定 门店信息
			devid = info.BeaconInfo.Minor
			data.At = strconv.Itoa(info.BeaconInfo.Minor) + "测试信息"
		}

		//  排队的人员列表  获取设备下的所有排队人数

		if u != nil {
			line := model.CreateLine(wxid, strconv.Itoa(devid))
			line.AddUser(u) // 排队啦
			data.Items = line.GetUsers()
		}

		log.Println(data)

		t, err := template.ParseFiles(f)
		err = t.Execute(w, data)
		if err != nil {
			log.Println(err)
		}

		return

	})



	// 首页的处理  只处理点击 过来的请求
	http.HandleFunc("/del-this", func(w http.ResponseWriter, r *http.Request) {

		// 简单的路由机制
		if strings.HasSuffix(r.URL.Path, "/") {

			log.Println("Path==>", r.URL.Path)
			//log.Println("http://", r.Host, r.RequestURI)

			/*
				path 的几种情况
				/		默认加载  index.html
				/aa/ 	默认加载  aa/index.html  如果不存在 加载 aa.html
				/aa/bb/ 默认加载  aa/bb/index.html 如果不存在加载 aa/bb.html
			*/
			path := strings.TrimRight(r.URL.Path, "/") //
			var f string                               // 模板文件路径

			//  模板文件检查

			f = filepath.Join(g.Root, "/public", path, "index.html")
			if !file.IsExist(f) {
				f = filepath.Join(g.Root, "/public", path+".html")
				if !file.IsExist(f) {
					log.Println("not find", f)
					http.NotFound(w, r)
					return
				}
			}

			// 参数检查
			queryValues, err := url.ParseQuery(r.URL.RawQuery)
			log.Println("ParseQuery", queryValues)
			if err != nil {
				log.Println("[ERROR] URL.RawQuery", err)
				w.WriteHeader(400)
				return
			}

			// 基本参数设置
			fullurl := "http://" + r.Host + r.RequestURI
			wxid := "gh_8ac8a8821eb9"
			appid := "wxacc105428fe41835" // 云喇叭服务号
			rand.Seed(time.Now().UnixNano())
			nonce := strconv.Itoa(rand.Intn(9999999999))
			ts := time.Now().Unix()

			var u *model.User
			u = nil
			//

			//oauth 跳转 ， 页面授权获取用户基本信息
			code := queryValues.Get("code") //  摇一摇入口 code 有效
			state := queryValues.Get("state")

			openid := queryValues.Get("openid") // H5入口 可能携带有 openid 信息
			uuid := queryValues.Get("uuid")

			if openid != "" { // 先判断本地用户信息是否存在
				u = model.CreateUser(wxid, openid)
				if u != nil && u.NickName == "" { // 用户数据本地没有
					u = nil
				}
			}

			// 用户信息不完善的处理  需要 oauth2
			if u == nil {
				// code 是空的 需要 重定向 https://open.weixin.qq.com/connect/oauth2/authorize?appid=APPID&redirect_uri=REDIRECT_URI&response_type=code&scope=SCOPE&state=STATE#wechat_redirect
				if code == "" && state == "" {
					addr := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=" + appid + "&redirect_uri=" + url.QueryEscape(fullurl) + "&response_type=code&scope=snsapi_userinfo&state=1#wechat_redirect"
					//log.Println(addr)
					http.Redirect(w, r, addr, 302)
					return

				} else {
					// code 一定存在
					c := g.GetWechatConfig(wxid)

					var token string

					openid, token = util.GetAccessTokenFromCode(c.AppId, c.AppSecret, code)

					u = model.CreateUser(wxid, openid)
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
				}
			}

			// 用户信息是ok的

			//  下面开始  render 页面了

			data := struct {
				AppId  string
				Ts     int64
				Nonce  string
				Sign   string
				OpenId string

				NickName string
				Sex      string

				Sn          string
				Uuid        string
				CompanyName string
				Name        string
				Userdata    string
				Note        string
				NotifyTime  string
				Other       string
				QRCodeUrl   string
				//Items []string
				//DistanceInfo string
				DebugT bool
			}{
				AppId:  appid,
				Ts:     ts,
				Nonce:  nonce,
				Sign:   util.WXConfigSign(g.GetJsApiTicket(wxid), nonce, strconv.FormatInt(ts, 10), fullurl),
				OpenId: openid,

				NickName: "",
				Sex:      "",

				Uuid: uuid,
				//Items: []string{},
				DebugT: false,
			}

			if u != nil {
				if u.Sex == oauth2.SexMale {
					data.Sex = "先生"
				} else if u.Sex == oauth2.SexFemale {
					data.Sex = "女士"
				}
				data.NickName = u.NickName

				/*
				if u.QRticket == "" {
					qr, e := account.CreateTemporaryQRCode(123, 604800, g.GetWechatAccessToken(wxid))
					if e == nil {
						u.QRticket = qr.Ticket
						u.Save()
					}
				}

				if u.QRticket != "" {
					data.QRCodeUrl = "https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=" + url.QueryEscape(u.QRticket)
				}

				if g.IsAdmin(openid) {
					log.Println("show weidian debug url -------")
					data.DebugT = true
				}
				*/
			}

			// 获取 设备信息
			/*

			 */

			// 设备关联 快递员 然后获取当前用户的 可取的包裹信息
			if u != nil {
				pkgs := model.GetPackages(u, uuid) //  uuid 可能为空
				//log.Println(pkgs)
				if len(pkgs) >= 1 {
					p := pkgs[0]

					data.Sn = p.Sn
					data.CompanyName = p.CompanyName
					data.Name = p.Name
					data.Userdata = p.Userdata
					data.Note = "快递" + p.Note
					data.NotifyTime = p.NotifyTime.Format("2006-01-02 15:04:05")

					//tmp := []string{ "通知时间：" + p.NotifyTime.Format("2006-01-02 15:04:05") ,
					//		"包裹编号：" + p.Sn ,
					// 		//"手机号码：" + p.To ,
					// 		"快递公司：" + p.CompanyName + " " + p.Name + " "  + p.Userdata ,
					// 		"通知内容：快递" + p.Note ,
					//		"    "}

					// data.Items = append(data.Items, tmp...)

					// p.Uuid    这个通知已经被点击啦  上报状态
					model.SendFeedback("WechatFeedback", p.To, p.Userdata, "1", "接收成功", p.Uuid, openid)
				} else {
					data.Other = "没有找到您的包裹信息"
				}

				//for _, p := range pkgs {
				//	//log.Println(i, p)
				//}
			}

			log.Println(data)

			t, err := template.ParseFiles(f)
			err = t.Execute(w, data)
			if err != nil {
				log.Println(err)
			}

			proc.PVCnt.Incr() //  广告投放的计数
			log.Println("PV", proc.PVCnt.Cnt)

			return
		}

		
	})

}
