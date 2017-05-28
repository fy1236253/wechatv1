package template

type TemplateMessage struct {
	ToUser     string `json:"touser"`             // 必须, 接受者OpenID
	TemplateId string `json:"template_id"`        // 必须, 模版ID
	URL        string `json:"url,omitempty"`      // 可选, 用户点击后跳转的URL, 该URL必须处于开发者在公众平台网站中设置的域中
	TopColor   string `json:"topcolor,omitempty"` // 可选, 整个消息的颜色, 可以不设置

	//RawJSONData json.RawMessage `json:"data"` 	  // 必须, JSON 格式的 []byte, 满足特定的模板需求
	Data TemplateData `json:"data"`
}

type TemplateData struct {
	First    KVData `json:"first"`
	Remark   KVData `json:"remark"`
	Keyword1 KVData `json:"keyword1"`
	Keyword2 KVData `json:"keyword2"`
	Keyword3 KVData `json:"keyword3"`
	Keyword4 KVData `json:"keyword4"`
	Keyword5 KVData `json:"keyword5"`
	Keyword6 KVData `json:"keyword6"`
}

type KVData struct {
	Value string `json:"value"`
	Color string `json:"color"`
}
