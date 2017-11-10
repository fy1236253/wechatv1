package http

import (
	"encoding/json"
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
	"regexp"
	"strconv"
	"strings"
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
		fullurl := "http://" + r.Host + r.RequestURI
		wxid := "gh_f353e8a82fe5"
		appid := "wxdfac68fcc7a48fca"
		var f string // 模板文件路径
		f = filepath.Join(g.Root, "/public", "scanFinish.html")
		if !file.IsExist(f) {
			log.Println("not find", f)
			http.NotFound(w, r)
			return
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
		file, _, _ := r.FormFile("img")
		timestamp := time.Now().UnixNano()
		uuid := strconv.FormatInt(timestamp, 10)
		if file == nil {
			log.Println("未检测到文件")
			return
		}
		f, e := os.Create(g.Root + "/public/upload/" + uuid + ".jpg")
		log.Println(e)
		defer f.Close()
		log.Println(f.Name())
		io.Copy(f, file)
		defer file.Close()
		os.Remove(g.Root + "/public/upload/" + uuid + ".jpg")
		sess, _ := globalSessions.SessionStart(w, r)
		defer sess.SessionRelease(w)
		model.CreatNewUploadImg(uuid, sess.Get("openid").(string))
		RenderJson(w, "123")
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
	http.HandleFunc("/edit_img", func(w http.ResponseWriter, r *http.Request) {
		ID := strings.Trim(r.URL.Path, "/edit_img/")
		log.Println(ID)
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
			UUID: ID,
		}

		t, err := template.ParseFiles(f)
		err = t.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
		return
	})
	http.HandleFunc("/handle", func(w http.ResponseWriter, r *http.Request) {
		type Res struct {
			Name   string
			Amount string
		}
		type WResult struct {
			Words string `json:"words"`
		}
		type Result struct {
			WordsResult []WResult `json:"words_result"`
		}

		var res Res
		var result Result
		var amount string
		str, _ := file.ToTrimString("local.json")
		log.Println(str)
		json.Unmarshal([]byte(str), &result)
		regular := "(华联).*(连锁)|(医药).*(连锁)"
		regular1 := "实收"
		reg := regexp.MustCompile(regular)
		reg1 := regexp.MustCompile(regular1)
		for _, v := range result.WordsResult {
			if reg.MatchString(v.Words) {
				res.Name = v.Words
			} else if reg1.MatchString(v.Words) {
				amount = v.Words
			}
		}
		log.Println(amount[2:])
		log.Println(amount[5:])
		res.Amount = amount
		RenderJson(w, res)
		return
	})
}
