package material

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"

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
func UpLodePIC(url, filepath, filename string) {
	pr, pw := io.Pipe()
	ws := multipart.NewWriter(pw)
	go func() {
		f, err := os.Open("header.jpeg")
		if err != nil {
			return
		}
		defer f.Close()
		fw, err := ws.CreateFormFile("file", "header.jpeg")
		if err != nil {
			log.Println(err)
			return
		}
		io.Copy(fw, f)
		ws.Close()
		pw.Close()
	}()
	log.Println("开始从管道读取数据")
	cli := http.Client{}
	resp, err := cli.Post(url, ws.FormDataContentType(), pr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("POST传输完成")

	body := resp.Body
	defer body.Close()

	if body_bytes, err := ioutil.ReadAll(body); err == nil {
		log.Println("response:", string(body_bytes))
	} else {
		log.Fatalln(err)
	}
}
