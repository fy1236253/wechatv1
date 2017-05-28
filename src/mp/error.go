package mp

const (
	ErrCodeOK                 = 0	
	ErrMsgOk                  = "success"
	ErrCodeInvalidCredential  = 40001 // access_token 过期(无效)返回这个错误
	ErrCodeAccessTokenExpired = 42001 // access_token 过期(无效)返回这个错误(maybe!!!)
)

type Error struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	MsgId   int64  `json:"msgid"`
}

type YzError struct {
	Code int `json:"code"`
	Msg  string `json:"msg"`
}