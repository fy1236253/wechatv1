package http

import (
	"fmt"
	"g"
	"io"
	"log"
	"mime/multipart"
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
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		//创建文件
		fW, err := os.Create(g.Root + "/public/img/" + head.Filename)
		if err != nil {
			fmt.Println("文件创建失败")
			return
		}
		defer fW.Close()
		_, err = io.Copy(fW, file)
		if err != nil {
			fmt.Println("文件保存失败")
			return
		}

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
