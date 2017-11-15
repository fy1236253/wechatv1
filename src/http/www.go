package http

import (
	"encoding/base64"
	"g"
	"html/template"
	"io"
	"log"
	"math/rand"
	"model"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"util"

	"github.com/toolkits/file"
)

// ConfigWebHTTP 对外http
func ConfigWebHTTP() {
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		//log.Println(openid)
		fullurl := "http://" + r.Host + r.RequestURI
		queryValues, err := url.ParseQuery(r.URL.RawQuery)
		log.Println("ParseQuery", queryValues)
		if err != nil {
			log.Println("[ERROR] URL.RawQuery", err)
			w.WriteHeader(400)
			return
		}
		wxid := "gh_f353e8a82fe5"
		wxcfg := g.GetWechatConfig(wxid)
		code := queryValues.Get("code") //  摇一摇入口 code 有效
		state := queryValues.Get("state")
		if code == "" && state == "" {
			addr := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=" + wxcfg.AppID + "&redirect_uri=" + url.QueryEscape(fullurl) + "&response_type=code&scope=snsapi_userinfo&state=1#wechat_redirect"
			log.Println("http.Redirect", addr)
			http.Redirect(w, r, addr, 302)
			return
		}
		log.Println(code)
		openid, _ := util.GetAccessTokenFromCode(wxcfg.AppID, wxcfg.AppSecret, code)
		log.Println(openid)
		var f string // 模板文件路径
		f = filepath.Join(g.Root, "/public", "index.html")
		if !file.IsExist(f) {
			log.Println("not find", f)
			http.NotFound(w, r)
			return
		}

		data := struct {
		}{}

		t, err := template.ParseFiles(f)
		err = t.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
		return
	})
	http.HandleFunc("/scanner", func(w http.ResponseWriter, r *http.Request) {
		fullurl := "http://" + r.Host + r.RequestURI
		wxid := "gh_f353e8a82fe5"
		appid := g.Config().Wechats[0].AppID
		appsecret := g.Config().Wechats[0].AppSecret
		queryValues, _ := url.ParseQuery(r.URL.RawQuery)
		code := queryValues.Get("code") //  摇一摇入口 code 有效
		if code == "" {
			addr := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=" + appid + "&redirect_uri=" + url.QueryEscape(fullurl) + "&response_type=code&scope=snsapi_userinfo&state=1#wechat_redirect"
			http.Redirect(w, r, addr, 302)
			return
		}
		openid, _ := util.GetAccessTokenFromCode(appid, appsecret, code)
		var f string // 模板文件路径
		f = filepath.Join(g.Root, "/public", "scan.html")
		if !file.IsExist(f) {
			log.Println("not find", f)
			http.NotFound(w, r)
			return
		}
		sess, _ := globalSessions.SessionStart(w, r)
		defer sess.SessionRelease(w)
		if sess.Get("openid") == nil {
			sess.Set("openid", openid)
		}
		// 基本参数设置
		rand.Seed(time.Now().UnixNano())
		nonce := strconv.Itoa(rand.Intn(9999999999))
		ts := time.Now().Unix()
		sign := util.WXConfigSign(g.GetJsAPITicket(wxid), nonce, strconv.FormatInt(ts, 10), fullurl)
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
	http.HandleFunc("/consumer", func(w http.ResponseWriter, r *http.Request) {
		var f string // 模板文件路径
		queryValues, _ := url.ParseQuery(r.URL.RawQuery)
		unionid := queryValues.Get("unionid")
		name := queryValues.Get("name")
		amount := queryValues.Get("amount")
		f = filepath.Join(g.Root, "/public", "scanFinish.html")
		if !file.IsExist(f) {
			log.Println("not find", f)
			http.NotFound(w, r)
			return
		}
		// 基本参数设置
		data := struct {
			Unionid string
			Name    string
			Amount  string
		}{
			Unionid: unionid,
			Name:    name,
			Amount:  amount,
		}
		t, err := template.ParseFiles(f)
		log.Println(err)
		err = t.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
		return
	})
	http.HandleFunc("/credits", func(w http.ResponseWriter, r *http.Request) {
		appid := "wxdfac68fcc7a48fca"
		var f string // 模板文件路径
		f = filepath.Join(g.Root, "/public", "scannerIndex.html")
		if !file.IsExist(f) {
			log.Println("not find", f)
			http.NotFound(w, r)
			return
		}
		// 基本参数设置
		data := struct {
			//Couriers 	string
			AppId string
		}{
			AppId: appid,
		}

		t, err := template.ParseFiles(f)
		err = t.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
		return
	})
	http.HandleFunc("/uploadImg", func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(32 << 20)
		sess, _ := globalSessions.SessionStart(w, r)
		defer sess.SessionRelease(w)
		openid := sess.Get("openid").(string)
		timestamp := time.Now().UnixNano()
		uuid := strconv.FormatInt(timestamp, 10)
		file, _, _ := r.FormFile("img")
		defer file.Close()
		rate := r.FormValue("rate")
		log.Println(rate)
		rateInt, _ := strconv.Atoi(rate)
		var result model.CommonResult
		if rateInt >= 2 {
			//人工处理模块
			log.Println("save handle img:" + uuid)
			f, e := os.Create("upload/" + uuid + ".jpg")
			log.Println(e)
			defer f.Close()
			_, e = io.Copy(f, file)
			log.Println(e)
			// model.CreatNewUploadImg(uuid, openid)
			result.ErrMsg = "1" //表示有错误
			RenderJson(w, result)
			return
		}
		if file == nil || openid == "" {
			log.Println("未检测到文件")
			return
		}
		sourcebuffer := make([]byte, 4*1024*1024) //最大4M
		n, _ := file.Read(sourcebuffer)
		base64Str := base64.StdEncoding.EncodeToString(sourcebuffer[:n])
		res := model.LocalImageRecognition(base64Str)
		result.ErrMsg = "success"
		if res == nil {
			log.Println("fail to upload")
			result.ErrMsg = "1" //表示有错误
			// return
		} else {
			result.DataInfo = res
		}
		log.Println(uuid)
		// model.CreatNewUploadImg(uuid, openid)
		RenderJson(w, result)
		return
	})
	http.HandleFunc("/hand_operation", func(w http.ResponseWriter, r *http.Request) {
		imgItems := model.GetUploadImgInfo()
		var f string // 模板文件路径
		f = filepath.Join(g.Root, "/public", "handOperation.html")
		if !file.IsExist(f) {
			log.Println("not find", f)
			http.NotFound(w, r)
			return
		}
		// 基本参数设置
		data := struct {
			//Couriers 	string
			ImgItems []string
		}{
			ImgItems: imgItems,
		}

		t, err := template.ParseFiles(f)
		err = t.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
		return
	})
	http.HandleFunc("/save_user_info", func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		amount := r.FormValue("amount")
		if name == "" || amount == "" {
			log.Println("[no paramas]")
			return
		}
		log.Println(name, amount)
		return
	})
	http.HandleFunc("/edit_img", func(w http.ResponseWriter, r *http.Request) {
		urlParse, _ := url.ParseQuery(r.URL.RawQuery)
		uuid := urlParse.Get("uuid")
		log.Println(uuid)
		var f string // 模板文件路径
		f = filepath.Join(g.Root, "/public", "edit.html")
		if !file.IsExist(f) {
			log.Println("not find", f)
			http.NotFound(w, r)
			return
		}
		// 基本参数设置
		data := struct {
			UUID string
		}{
			UUID: uuid,
		}

		t, err := template.ParseFiles(f)
		err = t.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
		return
	})
	http.HandleFunc("/handle", func(w http.ResponseWriter, r *http.Request) {
		model.ImportDatbase()
		return
	})
}
