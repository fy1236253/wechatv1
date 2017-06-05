package material

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"log"
)

// UploadNews 上传图文消息
func UploadNews(wxid string) {
	// url := "https://api.weixin.qq.com/cgi-bin/material/add_news?access_token=" + g.GetWechatAccessToken(wxid)
	url := "https://api.weixin.qq.com/cgi-bin/material/add_news?access_token=-ZRxRRs2mVhrBKzr1MTTrkYsPPli4fH1qrq6dbEklBrKHGhRyGj2T-BqILlDtsS0ZkZs4cTV_cdlt-1lVlLA6pcNLoyR1S0fULoVsDhCbXR3nGC-38_dLWtMK1ZA4Q4oHLRhADARVB"
	var news Article
	news.Author = ""
	news.Content = "hello world"
	news.ContentSourceURL = "http://www.baidu.com"
	news.Digest = "测试地址"
	news.ShowCoverPic = 0
	news.ThumbMediaID = "putcMP5_kwvnCEdux1dO0wZd6nPY1RBzRwTJ0TTg0U3kMg7mvh0O7zgH8gKGklEG"
	news.Title = "葫芦娃"
	buf := bytes.NewBuffer(make([]byte, 0, 16<<10))
	buf.Reset()
	json.NewEncoder(buf).Encode(news)
	body := buf.String()
	log.Println(body)
	// req := httplib.Post(url).SetTimeout(3*time.Second, 30*time.Second)
	// req.Body(body)
	// resp, err := req.String()
	// if err != nil {
	// 	log.Println(err)
	// }
	// log.Println(resp)

	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		log.Println(err)
		return
	}
	client := &http.Client{}
	res, err := client.Do(req)
	bodys, _ := ioutil.ReadAll(res.Body)
	log.Println(string(bodys))
	// url := "http://39.108.14.29/api/v1/upload/image"
	// url := "https://api.weixin.qq.com/cgi-bin/media/upload?type=image&access_token=3DqsZGRXlwPO_lF3RxKvplRZoifH8WrAbgnI6tVbr78v0GbrwVusgXmp0UENjygxfTqByQHRkykGscO6fiR3pN8Nw_XpoopFtd3o_Rl8tLr5VMSIbLWP_-Q-k7zbnXfiEFNbAHAKLT"
	// material.UploadLocalPic(url, g.Root+"header.jpeg", "header.jpeg")
	// material.UpLodePIC()
	// material.UploadNews("")
}
