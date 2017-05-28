package http

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/xml"
	"g"
	"io/ioutil"
	"log"
	"model"
	"mp"
	"mp/message"
	"mp/message/request"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"util"
)

// ConfigWechatRoutes 微信页面路由
func ConfigWechatRoutes() {

	//快递员 主叫号码设置 接口 ，
	//get
	http.HandleFunc("/wx/", func(w http.ResponseWriter, req *http.Request) {

		// 捕获异常
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Runtime error caught: %v", r)
				w.WriteHeader(400)
				w.Write([]byte(""))
				return
			}
		}()
		var wxcfg *mp.WechatConfig
		var queryValues url.Values

		wxid := strings.Trim(req.URL.Path, "/wx/")
		log.Println("wxid -->", wxid)   //    /gh_8ac8a8821eb9
		wxcfg = g.GetWechatConfig(wxid) // 通过微信id 获取 对接的配置信息
		if wxcfg == nil {
			panic("[Warn] wecat config not find")
		}
		queryValues, _ = url.ParseQuery(req.URL.RawQuery)
		model.WechatQueryParamsValid(queryValues)
		switch req.Method {
		case "GET":
			{
				model.WechatSignValid(wxcfg, queryValues)
				RenderText200(w, queryValues.Get("echostr"))
				return
			}

		case "POST":
			{

				if queryValues.Get("encrypt_type") == "aes" {
					var aesBody message.AesRequestBody
					var aeskey [32]byte // 秘钥
					var mixedMsg message.MixedMessage
					// 非加密码模式 不接入
					model.WechatStrValid(queryValues.Get("encrypt_type"), "aes", "[ERROR] encryptType not support")
					model.WechatMessageXMLValid(req, &aesBody)                                                  // xml 解析验证
					model.WechatStrValid(aesBody.ToUserName, wxcfg.WxID, "[Warn] wechat id mismatch, from err") // 来源验证
					model.WechatSignEncryptValid(wxcfg, queryValues, aesBody.EncryptedMsg)                      // 指纹验证

					k, _ := util.AESKeyDecode(wxcfg.Aeskey)
					copy(aeskey[:], k)

					// 解密
					encryptedMsgBytes, _ := base64.StdEncoding.DecodeString(aesBody.EncryptedMsg)
					_, rawMsgXML, appid, _ := util.AESDecryptMsg(encryptedMsgBytes, aeskey)
					model.WechatStrValid(string(appid), wxcfg.AppID, "[Warn] AppId mismatch")
					// 解密ok

					log.Println(string(rawMsgXML))

					if err := xml.Unmarshal(rawMsgXML, &mixedMsg); err != nil {
						log.Println("[Warn] rawMsgXML Unmarshal", err)
						w.WriteHeader(400)
						return
					}

					model.WechatStrValid(mixedMsg.ToUserName, wxcfg.WxID, "[Warn] mixedMsg.ToUserName mismatch, from err") // 来源验证

					textXML := ""

					switch mixedMsg.MsgType {

					// text
					case request.MsgTypeText:
						{
							textXML = model.ProcessWechatText(wxcfg, &mixedMsg) // 文本消息的处理逻辑
						}
					// event
					case request.MsgTypeEvent:
						{
							model.ProcessWechatEvent(wxcfg, &mixedMsg)
						}

					}
					// 做同步响应
					nonce := queryValues.Get("nonce")
					timestamp := queryValues.Get("timestamp")

					random := make([]byte, 16)
					tmp := md5.Sum([]byte(nonce + timestamp))
					copy(random[:16], tmp[:16]) // 设置随机数 一个简单的处理方法

					// 注意这不能返回 201
					ts, _ := strconv.ParseInt(timestamp, 10, 64)
					responseHTTPBody := message.AesResponseBody{
						EncryptedMsg: base64.StdEncoding.EncodeToString(util.AESEncryptMsg(random, []byte(textXML), wxcfg.AppID, aeskey)),
						Timestamp:    ts,
						Nonce:        nonce,
					}
					responseHTTPBody.MsgSignature = util.MsgSign(wxcfg.Token, timestamp, responseHTTPBody.Nonce, responseHTTPBody.EncryptedMsg)
					w.WriteHeader(200)
					RenderXML(w, responseHTTPBody) // 所有流程都采用异步处理， 所以不需要同步返回xml 数据
					return
				} else {
					var commonBody message.MixedMessage
					msg, _ := ioutil.ReadAll(req.Body)
					if err := xml.Unmarshal(msg, &commonBody); err != nil {
						log.Println("[Warn] body Unmarshal", err)
						w.WriteHeader(400)
						return
					}

					//log.Println(commonBody.Content)
					model.WechatStrValid(commonBody.ToUserName, wxcfg.WxID, "[Warn] commonBody.ToUserName mismatch, from err") // 来源验证
					switch commonBody.MsgType {
					// text
					case request.MsgTypeText:
						{
							model.ProcessWechatText(wxcfg, &commonBody) // 文本消息的处理逻辑
						}
					// event
					case request.MsgTypeEvent:
						{
							model.ProcessWechatEvent(wxcfg, &commonBody)
						}
					}
					w.WriteHeader(200)
					RenderText(w, "")
				}
			}
		}

	})
}
