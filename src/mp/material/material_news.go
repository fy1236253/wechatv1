package material

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/toolkits/net/httplib"

	"log"
)

// UploadNews 上传图文消息
func UploadNews(wxid string) {
	// url := "https://api.weixin.qq.com/cgi-bin/material/add_news?access_token=" + g.GetWechatAccessToken(wxid)
	url := "https://api.weixin.qq.com/cgi-bin/material/add_news?access_token=vKrlEbWgsjVf4trOtxkhqzwckKw23ym6_rb-oay8EZmfr-ReqAmyTJzFJ6Lfi78MGKBAH2IIlyl3miJcKL3uiSlnHGyBpexw6GVfya8GUEINZAfACAQCD"
	var news Article
	news.Author = "body"
	news.Content = "hello world"
	news.ContentSourceURL = "http://www.baidu.com"
	news.Digest = "测试地址"
	news.ShowCoverPic = 0
	news.ThumbMediaID = "putcMP5_kwvnCEdux1dO0wZd6nPY1RBzRwTJ0TTg0U3kMg7mvh0O7zgH8gKGklEG"
	news.Title = "葫芦娃"
	var newsItem News
	newsItem.Article = append(newsItem.Article, news)
	buf := bytes.NewBuffer(make([]byte, 0, 16<<10))
	buf.Reset()

	json.NewEncoder(buf).Encode(newsItem)
	body := buf.String()
	log.Println(body)

	req := httplib.Post(url).SetTimeout(3*time.Second, 30*time.Second)
	req.Body(body)
	resp, err := req.String()
	if err != nil {
		log.Println(err)
	}
	log.Println(resp)

}
