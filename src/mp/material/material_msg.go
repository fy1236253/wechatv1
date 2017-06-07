package material

// News 多图文结构
type News struct {
	Article []Article `json:"article"`
}

// Article 图文消息结构体
type Article struct {
	ThumbMediaID     string `json:"thumb_media_id"`               // 必须; 图文消息的封面图片素材id(必须是永久mediaID)
	Title            string `json:"title"`                        // 必须; 标题
	Author           string `json:"author,omitempty"`             // 必须; 作者
	Digest           string `json:"digest,omitempty"`             // 必须; 图文消息的摘要, 仅有单图文消息才有摘要, 多图文此处为空
	Content          string `json:"content"`                      // 必须; 图文消息的具体内容, 支持HTML标签, 必须少于2万字符, 小于1M, 且此处会去除JS
	ContentSourceURL string `json:"content_source_url,omitempty"` // 必须; 图文消息的原文地址, 即点击"阅读原文"后的URL
	ShowCoverPic     int    `json:"show_cover_pic"`               // 必须; 是否显示封面, 0为false, 即不显示, 1为true, 即显示
	URL              string `json:"url,omitempty"`                // !!!创建的时候不需要此参数!!! 图文页的URL, 文章创建成功以后, 会由微信自动生成
}
