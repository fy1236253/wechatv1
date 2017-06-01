package http

import (
	"g"
	"log"
	"mime/multipart"
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
	http.HandleFunc("/api/v1/upload/image", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		log.Println(r.Method)
		if "POST" == r.Method {
		}
		queryValues, err := url.ParseQuery(r.URL.RawQuery)
		log.Println("ParseQuery", queryValues)
		if err != nil {
			log.Println("[ERROR] URL.RawQuery", err)
			w.WriteHeader(400)
			return
		}
		r.ParseMultipartForm(32 << 20)
		form := r.MultipartForm
		logMultipartForm(form)
	})
}

func logMultipartForm(form *multipart.Form) {
	log.Print("Values:", form.Value)
	log.Print("Files:")
	for key := range form.File {
		headers := form.File[key]
		for _, header := range headers {
			log.Printf("Key: %v, Filename: %v, Header: %v", key, header.Filename, header.Header)
		}
	}
}
