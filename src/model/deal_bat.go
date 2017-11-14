package model

import (
	"encoding/json"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/toolkits/net/httplib"
)

// BATList 百度数据返回
type BATList struct {
	Words string `json:"words"` //目前只正则匹配汉字
}

// BATResult 识别的列表
type BATResult struct {
	WordsResult []BATList `json:"words_result"`
}

// RecognizeResult 自处理后的结果
type RecognizeResult struct {
	ShopName    string  `json:"shop_name"`
	TotalAmount float64 `json:"total_amount"`
	Uninoid     string  `json:"unionid"`
}

// CommonResult api接口返回数据
type CommonResult struct {
	ErrData  string      `json:"errMsg"`
	ErrMsg   string      `json:"msgInfo"`
	DataInfo interface{} `json:"data"`
}

// BatImageRecognition 百度的图像识别接口
func BatImageRecognition(base64Str string) string {
	url := "https://aip.baidubce.com/rest/2.0/ocr/v1/accurate?access_token=24.1f248484d5b7faf54537dfae92fed52c.2592000.1512598910.282335-10330945"
	req := httplib.Post(url).SetTimeout(3*time.Second, 1*time.Minute)
	req.Header("Content-Type", "application/x-www-form-urlencoded")
	// req.Body("{\"img\":" + base64Str + "}")
	req.Param("image", base64Str)
	resp, err := req.String()
	if err != nil {
		log.Println(err)
		return ""
	}
	return resp
}

// LocalImageRecognition 自由图片处理 提取数据
func LocalImageRecognition(base64 string) *RecognizeResult {
	resp := BatImageRecognition(base64)
	if resp == "" {
		log.Println("request BAT fail")
		return nil
	}
	log.Println(resp)
	var res BATResult
	var amountFloat, amount float64
	var unionid string
	result := new(RecognizeResult)
	json.Unmarshal([]byte(resp), &res)
	for _, v := range res.WordsResult { //轮训关键字
		log.Println(v)
		name := recongnitionName(v.Words)
		if name != "" {
			result.ShopName = name
		}
		amountFloat = recongnitionAmount(v.Words)
		if amountFloat >= amount {
			amount = amountFloat
		}
		id := recongnitionOrderNum(v.Words)
		if id != "" {
			unionid = id
		}
	}
	result.TotalAmount = amount
	result.Uninoid = unionid
	return result
}

// RecongnitionOrderNum 处理订单中的编号
func recongnitionOrderNum(str string) string { //加上单据号搜索
	regular := `^(单据号|单据).\d[0-9]+|\d{15}`
	match, name := commonMatch(regular, str)
	reg := regexp.MustCompile("[\u4E00-\u9FA5].")
	name = reg.ReplaceAllLiteralString(name, "")
	if match {
		return name
	}
	return ""
}

// recongnitionAmount 识别订单中的金额
func recongnitionAmount(str string) float64 {
	regular := `\d+\..*\d+`
	//	regular := `(-?\d*)\.?\d+`
	match, amount := commonMatch(regular, str)
	if match {
		amountFloat, _ := strconv.ParseFloat(amount, 64)
		return amountFloat
	}
	return 0
}

// recongnitionName 匹配订单中的药店名称
func recongnitionName(str string) string {
	regular := `.*.(大药房)|.*.(连锁店)`
	match, name := commonMatch(regular, str)
	name = strings.Replace(name, "落款单位:", "", -1)
	if match {
		return name
	}
	return ""
}

// commonMatch 通用正则匹配
func commonMatch(regular, str string) (bool, string) {
	reg := regexp.MustCompile(regular)
	name := reg.FindAllString(str, -1)
	match := reg.MatchString(str)
	if match {
		//		log.Println(name)
		return true, name[0]
	}
	return false, ""
}
