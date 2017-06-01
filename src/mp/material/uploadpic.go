package material

import (
	"g"

	"log"
	"os"

	"github.com/toolkits/net/httplib"
)

// UpLodePIC 上传图片素材
func UpLodePIC(wxid string) {
	url := "https://api.weixin.qq.com/cgi-bin/media/upload?access_token=" + g.GetWechatAccessToken(wxid) + "&type=image"
	file, err := os.Open("/usr/local/src/wechatv1/public/img/u1604.png")
	log.Println(err)
	defer file.Close()
	req := httplib.Post(url)
	req.PostFile("media", g.Root+"/public/img/u1355.png")
	// reqs, errs := req.Post(url, req.FileUpload{
	// 	File:      file,
	// 	FieldName: "media",     // FieldName 是表单字段名
	// 	FileName:  "u1604.png", // Filename 是要上传的文件的名称，我们使用它来猜测mimetype，并将其上传到服务器上
	// })
	// log.Println(errs)
	resp, err := req.String()
	log.Println(resp, err)
}
