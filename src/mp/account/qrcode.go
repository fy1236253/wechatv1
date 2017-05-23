// 二维码接口

package account

import ()

const (
	TemporaryQRCodeExpireSecondsLimit = 604800 // 临时二维码 expire_seconds 限制
	PermanentQRCodeSceneIdLimit       = 100000 // 永久二维码 scene_id 限制

	QR_SCENE 			= "QR_SCENE"
	QR_LIMIT_SCENE 		= "QR_LIMIT_SCENE"
	QR_LIMIT_STR_SCENE 	= "QR_LIMIT_STR_SCENE"
)

// 永久二维码
type PermanentQRCode struct {
	// 下面两个字段同时只有一个有效, 非zero值表示有效.
	SceneId     uint32 `json:"scene_id,omitempty"`  // 场景值ID, 临时二维码时为32位非0整型, 永久二维码时最大值为100000(目前参数只支持1--100000)
	SceneString string `json:"scene_str,omitempty"` // 场景值ID(字符串形式的ID), 字符串类型, 长度限制为1到64, 仅永久二维码支持此字段

	Ticket string `json:"ticket"` // 获取的二维码ticket, 凭借此ticket可以在有效时间内换取二维码.
	URL    string `json:"url"`    // 二维码图片解析后的地址, 开发者可根据该地址自行生成需要的二维码图片
}

// 临时二维码
type TemporaryQRCode struct {
	ExpireSeconds int `json:"expire_seconds,omitempty"` // 二维码的有效时间, 以秒为单位. 最大不超过 604800.
	PermanentQRCode
}
