package g

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/toolkits/file"
)

type HttpConfig struct {
	Enable bool   `json:"enable"`
	Listen string `json:"listen"`
}

type AmqpConfig struct {
	Addr    string `json:"addr"`
	Addr1   string `json:"addr1"`
	Addr2   string `json:"addr2"`
	MaxIdle int    `json:"maxIdle"`
}

type RedisConfig struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	MaxIdle  int    `json:"maxIdle"`
	Db       int    `json:"db"`
}

type DBConfig struct {
	Dsn     string `json:"dsn"`
	MaxIdle int    `json:"maxIdle"`
}

type WorkerConfig struct {
	Wechat int `json:"wechat"`
}

type WechatConfig struct {
	WxId         string `json:"WxId"`
	AppSecret    string `json:"AppSecret"`
	AppId        string `json:"AppId"`
	Token        string `json:"Token"`
	Aeskey       string `json:"Aeskey"`
	accessToken  string // 这个是通过接口请求获取到的
	jsapi_ticket string
	AutoAnswer   bool   `json:"AutoAnswer"`
	Welcome      string `json:"Welcome"`
}

type AdminConfig struct {
	Openid   string `json:"openid"`
	Nickname string `json:"nickname"`
}

type OpenConfig struct {
	Addr     string `json:"addr"`
	AppId    string `json:"appid"`
	AppToken string `json:"apptoken"`
}

type RpcConfig struct {
	Addr     string `json:"addr"`
	AppId    string `json:"appid"`
	AppToken string `json:"apptoken"`
}

type LocationConfig struct {
	Addr     string `json:"addr"`
	AppId    string `json:"appid"`
	AppToken string `json:"apptoken"`
}

type WeixinPayConfig struct {
	Addr          string `json:"addr"`
	Key        string `json:"key"`
	Ip         string `json:"ip"` // ip  白名单 
	P          int  `json:"p"`  // 中奖的概率  100  50  到 0 
}

// 号段 设置 
type HaoduanConfig struct {
	ChinaMobile 	string `json:"chinaMobile"` 
	ChinaUnicom 	string `json:"chinaUnicom"` 
	ChinaTelecom 	string `json:"chinaTelecom"` 
	ChinaOther 		string `json:"chinaOther"` 
}

type GlobalConfig struct {
	Debug         bool             `json:"debug"`
	Logs          string           `json:"logs"`
	AdMsg         string           `json:"ad-msg"`
	Http          *HttpConfig      `json:"http"`
	Amqp          *AmqpConfig      `json:"amqp"`
	Redis         *RedisConfig     `json:"redis"`
	RedisLocalNet *RedisConfig     `json:"redis-local-net"`
	Haoduan       *HaoduanConfig   `json:"haoduan"`
	DB            *DBConfig        `json:"db"`
	Worker        *WorkerConfig    `json:"worker"`
	Wechats       []*WechatConfig  `json:"wechats"`
	Admins        []*AdminConfig   `json:"admins"`
	Open          *OpenConfig      `json:"open"`
	Rpc           *RpcConfig       `json:"rpc"`
	Location      *LocationConfig  `json:"location"`
	WeixinPay     *WeixinPayConfig `json:"weixinpay"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

func ParseConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("config file not specified: use -c $filename")
	}

	if !file.IsExist(cfg) {
		log.Fatalln("config file specified not found:", cfg)
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file", cfg, "error:", err.Error())
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("parse config file", cfg, "error:", err.Error())
	}

	// set config
	configLock.Lock()
	defer configLock.Unlock()
	config = &c

	log.Println("g.ParseConfig ok, file", cfg)
}
