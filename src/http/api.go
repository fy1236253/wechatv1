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

	"github.com/imroc/req"
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
		// form := r.MultipartForm
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

		url := "https://api.weixin.qq.com/cgi-bin/media/upload?access_token=" + g.GetWechatAccessToken("gh_f353e8a82fe5") + "&type=image"
		req.Post(url, req.FileUpload{
			File:      file,
			FieldName: "file",       // FieldName 是表单字段名
			FileName:  "avatar.png", // Filename 是要上传的文件的名称，我们使用它来猜测mimetype，并将其上传到服务器上
		})
		_, err = io.Copy(fW, file)
		if err != nil {
			fmt.Println("文件保存失败")
			return
		}
		log.Println("保存成功" + head.Filename)
	})
}
