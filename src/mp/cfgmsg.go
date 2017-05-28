package mp

//GlobalConfig 微信配置结构体
type GlobalConfig struct {
	Debug         bool             `json:"debug"`
	Logs          string           `json:"logs"`
	AdMsg         string           `json:"ad-msg"`
	HTTP          *HTTPConfig      `json:"http"`
	Amqp          *AmqpConfig      `json:"amqp"`
	Redis         *RedisConfig     `json:"redis"`
	RedisLocalNet *RedisConfig     `json:"redis-local-net"`
	DB            *DBConfig        `json:"db"`
	Worker        *WorkerConfig    `json:"worker"`
	Wechats       []*WechatConfig  `json:"wechats"`
	WeixinPay     *WeixinPayConfig `json:"weixinpay"`
	Admins        []*AdminsConfig  `json:"admins"`
}

//AdminsConfig 端口绑定
type AdminsConfig struct {
	Openid   string `json:"openid"`
	Nickname string `json:"nickname"`
}

//HTTPConfig 端口绑定
type HTTPConfig struct {
	Enable bool   `json:"enable"`
	Listen string `json:"listen"`
}

//AmqpConfig rabbitmq地址
type AmqpConfig struct {
	Addr    string `json:"addr"`
	Addr1   string `json:"addr1"`
	Addr2   string `json:"addr2"`
	MaxIdle int    `json:"maxIdle"`
}

//RedisConfig redis配置
type RedisConfig struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	MaxIdle  int    `json:"maxIdle"`
	Db       int    `json:"db"`
}

//DBConfig mysql配置
type DBConfig struct {
	Dsn     string `json:"dsn"`
	MaxIdle int    `json:"maxIdle"`
}

//WorkerConfig worker数量配置
type WorkerConfig struct {
	Wechat int `json:"wechat"`
}

//WechatConfig 微信公众号配置，支持多个
type WechatConfig struct {
	WxID        string `json:"WxId"`
	AppSecret   string `json:"AppSecret"`
	AppID       string `json:"AppId"`
	Token       string `json:"Token"`
	Aeskey      string `json:"Aeskey"`
	AccessToken string // 这个是通过接口请求获取到的
	JsapiTicket string
	AutoAnswer  bool   `json:"AutoAnswer"`
	Welcome     string `json:"Welcome"`
}

//WeixinPayConfig 微信支付相关
type WeixinPayConfig struct {
	Addr string `json:"addr"`
	Key  string `json:"key"`
	IP   string `json:"ip"` // ip  白名单
	P    int    `json:"p"`  // 中奖的概率  100  50  到 0
}
