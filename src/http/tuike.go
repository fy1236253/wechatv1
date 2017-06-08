package http

import (
	"log"
	"net/http"
	"net/url"
)

// Config3rdWechatRoutes 微信页面路由
func Config3rdWechatRoutes() {

	http.HandleFunc("/3rd/", func(w http.ResponseWriter, req *http.Request) {

		// 捕获异常
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Runtime error caught: %v", r)
				w.WriteHeader(400)
				w.Write([]byte(""))
				return
			}
		}()
		queryValues, _ := url.ParseQuery(req.URL.RawQuery)
		log.Println(queryValues)
	})
}
