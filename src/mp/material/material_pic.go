package material

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/toolkits/net/httplib"

	"log"
)

// UploadLocalPic 上传图片素材（微信需要一次性上传）
func UploadLocalPic(url, filepath, filename string) {
	var b bytes.Buffer

	// pr, pw := io.Pipe()
	w := multipart.NewWriter(&b)
	// w := multipart.NewWriter(pw)

	go func() {
		f, err := os.Open(filepath)
		if err != nil {
			return
		}
		defer f.Close()
		fw, err := w.CreateFormFile("file", filename)
		if err != nil {
			log.Println(err)
			return
		}
		io.Copy(fw, f)
		w.Close()
		// pw.Close()
	}()

	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		log.Println(err)
		return
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	client := &http.Client{}
	res, err := client.Do(req)
	log.Println(req.PostForm)
	body, _ := ioutil.ReadAll(res.Body)
	log.Println(string(body))
}

// UpLodePIC 上传图片素材
func UpLodePIC() {
	// url := "https://api.weixin.qq.com/cgi-bin/media/uploadimg?access_token=JMJCtJ4L8T2Nbigv9fYgQ92oczZAljvgDfP3gp5OOSK_d1eVJ6OWsW0fcjJdDI-QlQTBsSH8i0rt9Yti8oX7BLpu77HOXdqEK3mAkTshFM4IhRAzjubE_J8uLPLj60NAKJEhAEADCX"
	url := "http://39.108.14.29/api/v1/upload/image"
	file, _ := os.Open("/Users/fengya/go/localtest/header.jpeg")
	defer file.Close()
	req := httplib.Post(url).SetTimeout(3*time.Second, 10*time.Second)
	req.PostFile("file", "header.jpeg")
	resp, err := req.String()
	log.Println(resp, err)

	// url := "https://api.weixin.qq.com/cgi-bin/media/uploadimg?access_token=JMJCtJ4L8T2Nbigv9fYgQ92oczZAljvgDfP3gp5OOSK_d1eVJ6OWsW0fcjJdDI-QlQTBsSH8i0rt9Yti8oX7BLpu77HOXdqEK3mAkTshFM4IhRAzjubE_J8uLPLj60NAKJEhAEADCX"
	// file, err := os.Open("/Users/fengya/go/wechatv1/sea.jpeg")
	// log.Println(err)
	// defer file.Close()
	// req.SetTimeout(30 * time.Second)

	// reqs, errs := req.Post(url, req.FileUpload{
	// 	File:      file,
	// 	FieldName: "file",     // FieldName 是表单字段名
	// 	FileName:  "sea.jpeg", // Filename 是要上传的文件的名称，我们使用它来猜测mimetype，并将其上传到服务器上
	// })
	// log.Println(errs)
	// resp := reqs.String()
	// log.Println(resp, err)

}
