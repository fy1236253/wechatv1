package material

import (
	"g"

	"log"

	"github.com/toolkits/net/httplib"
)

// UpLodePIC 上传图片素材
func UpLodePIC(wxid string) {
	url := "https://api.weixin.qq.com/cgi-bin/media/upload?access_token=" + g.GetWechatAccessToken(wxid) + "&type=image"
	req := httplib.Post(url)
	req.PostFile("image", g.Root+"/public/img/u1355.png")
	resp, err := req.String()
	if err == nil {
		log.Println("[上传失败]", err)
	}
	log.Println(resp)

}
