
package model

import (
	//"mq"
	"strconv"
	"time"
	"math/rand"

	"bytes"
	"encoding/xml"
	"g"
	"sort"
	"log"
	"crypto/md5"
	"strings"
	"encoding/hex"
   	"github.com/toolkits/net/httplib"
   	"hash"
	redispool "redis"
	//"github.com/garyburd/redigo/redis"	
)

// 微信 订单模型 
// 统一下单， 付款结果查询 

type WeixinUniOrder struct {
	XMLName 	struct{} `xml:"xml"`
	Appid    	string `xml:"appid"` 
	MchId    	string `xml:"mch_id"` 
	NonceStr    string `xml:"nonce_str"` 
	Sign    	string `xml:"sign"`    	// 签名 
	Body 		string `xml:"body"`   	// 商品描述
	Uid    		string `xml:"out_trade_no"`  // 内部订单号   32 字符
	TotalFee    int `xml:"total_fee"` // 订单金额  注意是 分 
	ClientIp 	string `xml:"spbill_create_ip"` // 用户端ip 
	NotifyUrl   string `xml:"notify_url"`   // 异步通知地址 
	TradeType   string `xml:"trade_type"`   // JSAPI
	OpenId      string `xml:"openid"`   // JSAPI 订单，必须要传递用户的openid 信息 
}

type WeixinUniOrderResult struct {
	XMLName 	struct{} `xml:"xml"`
	ReturnCode 	string `xml:"return_code"` 
	ResultCode 	string `xml:"result_code"` 
	PrepayId   	string `xml:"prepay_id"` 
	ErrCode   	string `xml:"err_code_des"`
}


// 订单查询接口 
type WeixinOrderQuery struct {
	XMLName 	struct{} `xml:"xml"`
	Appid    	string `xml:"appid"` 
	MchId    	string `xml:"mch_id"` 
	Uid    		string `xml:"out_trade_no"`  // 内部订单号   32 字符
	NonceStr    string `xml:"nonce_str"` 
	Sign    	string `xml:"sign"`    	// 签名 
}

type WeixinOrderQueryResult struct {
	XMLName 	struct{} `xml:"xml"`
	ReturnCode 	string `xml:"return_code"` 
	ResultCode 	string `xml:"result_code"` 
	TradeState 	string `xml:"trade_state"`   // SUCCESS—支付成功 
}



// 订单金额  val 是元为单位   string 类型 
// wxid, uid, i, m, t, y, vf, GetClientIp(r)
func WeixinOrder(wxid, uid, openid, to, tname, oldv string, val float64, ip string) (bool, string) {

	rand.Seed(time.Now().UnixNano())
	nonce := strconv.Itoa(rand.Intn(999999999)) 

	o := &WeixinUniOrder{
		Appid: "wxacc105428fe41835", 
		MchId: "1283952401",
		NonceStr: nonce,  
		Sign: "", 
		Body: "云喇叭流量 " + tname + " " + to,
		Uid: uid,
		TotalFee: int(val * 100),  //  需要转为分 单位  
		ClientIp: ip, 
		NotifyUrl: "http://weichat2.shenbianvip.com/wx/wxacc105428fe41835",
		TradeType: "JSAPI",
		OpenId: openid, 
	}

	o.Sign = ordersign(o, g.Config().WeixinPay.Key)


	buf := bytes.NewBuffer(make([]byte, 0, 16<<10))
	buf.Reset()
	xml.NewEncoder(buf).Encode(o)
	body := buf.String()

	log.Println(body)


	r := httplib.Post("https://api.mch.weixin.qq.com/pay/unifiedorder").SetTimeout(3*time.Second, 1*time.Minute)
	r.Header("Content-Type", "application/xml;charset=UTF-8")
	r.Body(body)

	resp, err := r.String()

	if err != nil {
		log.Println("[ERROR] WeixinOrder", err)
		return false, ""
	}

/*
<xml>
<return_code><![CDATA[SUCCESS]]></return_code>
<return_msg><![CDATA[OK]]></return_msg>
<appid><![CDATA[wxacc105428fe41835]]></appid>
<mch_id><![CDATA[1283952401]]></mch_id>
<nonce_str><![CDATA[PFZZRU4vxUalS93f]]></nonce_str>
<sign><![CDATA[AD937178AE6C7759939399734B41A987]]></sign>
<result_code><![CDATA[SUCCESS]]></result_code>
<prepay_id><![CDATA[wx201605181201222c5c43adbe0059637895]]></prepay_id>
<trade_type><![CDATA[JSAPI]]></trade_type>
</xml>
*/
	log.Println("weixin WeixinOrder result", resp, uid )

	var result WeixinUniOrderResult

	if err = xml.Unmarshal([]byte(resp), &result); err != nil {
		log.Println("[ERROR] xml ", err, resp)
		return false, ""
	}

	if result.ReturnCode == "SUCCESS" && result.ResultCode == "SUCCESS" {
		// 订单数据 保存在本地 redis中 
		conn := redispool.ConnPool.Get()
		defer conn.Close()

		//                                预订单            支付ok     预充流量         流量到账      真正支付     手机号码    套餐名称        原价
		conn.Do("HMSET", uid, "prepayid", result.PrepayId, "pay", 0, "precharge", 0, "charge", 0, "val", val, "to", to, "tname", tname, "oldv", oldv)
		conn.Do("EXPIRE",uid, 7*24*3600)  // 保存一周

		return true, result.PrepayId
	}

	return false, result.ErrCode
}


// 订单 支付状态查询
func WeixinOrderStatus(wxid, uid string) (bool) {

	rand.Seed(time.Now().UnixNano())

	data := map[string]string{ "appid": "wxacc105428fe41835",
		"mch_id": "1283952401",
		"out_trade_no": uid, 
		"nonce_str": strconv.Itoa(rand.Intn(999999999)) ,
	}

	sign := WxPaySign(data, g.Config().WeixinPay.Key, nil )

	o := &WeixinOrderQuery{
		Appid: 		data["appid"], 
		MchId: 		data["mch_id"], 
		NonceStr: 	data["nonce_str"], 
		Uid: 		data["out_trade_no"], 
		Sign: 		sign, 
	}

	
	buf := bytes.NewBuffer(make([]byte, 0, 16<<10))
	buf.Reset()
	xml.NewEncoder(buf).Encode(o)
	body := buf.String()

	log.Println(body)


	r := httplib.Post("https://api.mch.weixin.qq.com/pay/orderquery").SetTimeout(3*time.Second, 1*time.Minute)
	r.Header("Content-Type", "application/xml;charset=UTF-8")
	r.Body(body)

	resp, err := r.String()

	if err != nil {
		log.Println("[ERROR] WeixinOrder", err)
		return false
	}

	log.Println("weixin WeixinOrderQuery result", resp, uid )

	var result WeixinOrderQueryResult

	if err = xml.Unmarshal([]byte(resp), &result); err != nil {
		log.Println("[ERROR] xml ", err, resp)
		return false
	}

	if result.ReturnCode == "SUCCESS" && result.ResultCode == "SUCCESS" && result.TradeState == "SUCCESS" {
		return true  // 支付成功啦  
	}

	return false 
}


func WxPaySign(parameters map[string]string, apiKey string, fn func() hash.Hash) string {
	ks := make([]string, 0, len(parameters))
	for k := range parameters {
		if k == "sign" {
			continue
		}
		ks = append(ks, k)
	}
	sort.Strings(ks)

	if fn == nil {
		fn = md5.New
	}
	h := fn()

	buf := make([]byte, 256)
	for _, k := range ks {
		v := parameters[k]
		if v == "" {
			continue
		}

		buf = buf[:0]
		buf = append(buf, k...)
		buf = append(buf, '=')
		buf = append(buf, v...)
		buf = append(buf, '&')
		h.Write(buf)
	}
	buf = buf[:0]
	buf = append(buf, "key="...)
	buf = append(buf, apiKey...)
	h.Write(buf)

	signature := make([]byte, h.Size()*2)
	hex.Encode(signature, h.Sum(nil))
	return string(bytes.ToUpper(signature))
}



// 需要整理 
func ordersign(o *WeixinUniOrder, key string) string {
	strs := sort.StringSlice{ "appid=" + o.Appid, 
		"openid=" + o.OpenId, 
		"mch_id=" + o.MchId, 
		"nonce_str=" + o.NonceStr, 
		"body=" + o.Body, 
		"out_trade_no=" + o.Uid, 
		"total_fee=" + strconv.Itoa(o.TotalFee), 
		"spbill_create_ip=" + o.ClientIp, 
		"notify_url=" + o.NotifyUrl, 
		"trade_type=" + o.TradeType }
	strs.Sort()

	strA := strings.Join(strs[:], "&")

	//log.Println(strA)

	strB := strA + "&key=" + key 

	//log.Println(strB)

	md5Sum := md5.Sum([]byte(strB))
    sig  := strings.ToUpper(hex.EncodeToString(md5Sum[:]))

    //log.Println(sig)

	return sig
}

