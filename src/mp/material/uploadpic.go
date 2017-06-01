package material

import (
	"os"

	"github.com/imroc/req"

	"log"
)

// UpLodePIC 上传图片素材
func UpLodePIC(wxid string) {
	url := "https://api.weixin.qq.com/cgi-bin/media/upload?access_token=V_IXTnNdFxalRPh3wlfHRWHg0Kw-7Tjp4f7tpe2o60Qi6CXZ9NE8xv9rEqTHbAJAuqR8SvRM6K6lVE-THkQskF0BK8dbx3mvcLIOzlXctphCczeHPDSzBuI797CIK0m7ZNYdADAHVG&type=image"
	// urls := "/api/v1/upload/image"
	file, err := os.Open("/usr/local/src/wechatv1/public/img/duo1.jpeg")
	log.Println(err)
	defer file.Close()
	// req := httplib.Post(urls).SetTimeout(3*time.Second, 1*time.Minute)
	// req.PostFile("file", g.Root+"/public/img/duo1.jpeg")
	reqs, errs := req.Post(url, req.FileUpload{
		File:      file,
		FieldName: "media",     // FieldName 是表单字段名
		FileName:  "duo1.jpeg", // Filename 是要上传的文件的名称，我们使用它来猜测mimetype，并将其上传到服务器上
	})
	log.Println(errs)
	resp := reqs.String()
	log.Println(resp, err)
}
