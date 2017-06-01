package http

import (
	"fmt"
	"g"
	"io"
	"log"
	"mp/menu"
	"net/http"
	"net/url"
	"os"
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
	http.HandleFunc("/api/v1/upload/image", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		fmt.Println(r.FormFile("file"))
		queryValues, err := url.ParseQuery(r.URL.RawQuery)
		log.Println("ParseQuery", queryValues)
		if err != nil {
			log.Println("[ERROR] URL.RawQuery", err)
			w.WriteHeader(400)
			return
		}
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		log.Println(g.Root)
		f, err := os.OpenFile(g.Root+"/public/img/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	})
}
