package http

import (
	"g"
	"log"
	"mp/menu"
	"net/http"
	"net/url"
)

// ConfigAPIRoutes api相关接口
func ConfigAPIRoutes() {
	http.HandleFunc("/api/v1/createmenu", func(w http.ResponseWriter, r *http.Request) {
		queryValues, err := url.ParseQuery(r.URL.RawQuery)
		log.Println("ParseQuery", queryValues)
		if err != nil {
			log.Println("[ERROR] URL.RawQuery", err)
			w.WriteHeader(400)
			return
		}
		cfg := queryValues.Get("cfg")
		wxid := queryValues.Get("wxid")
		menu.CreateMenu(cfg, g.GetWechatAccessToken(wxid))
	})
}
