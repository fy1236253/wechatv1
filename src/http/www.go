package http

import (
	"g"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
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
		var f string // 模板文件路径
		f = filepath.Join(g.Root, "/public", "scan.html")
		if !file.IsExist(f) {
			log.Println("not find", f)
			http.NotFound(w, r)
			return
		}
		// 基本参数设置
		data := struct {
			//Couriers 	string
			Wxid string
		}{
			Wxid: "",
		}

		t, err := template.ParseFiles(f)
		err = t.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
		return
	})

}
