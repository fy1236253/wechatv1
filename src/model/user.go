package model

import (
	//"bytes"
	"encoding/json"
	"log"
	"mq"
	"g"
	//"mp/message/custom"
	"math/rand"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
	redispool "redis"

	//"database/sql"
)

const (
	UserNormal         = 0 // 普通状态
	UserBindWaitMobile = 1 // 开始绑定手机
	UserBindWaitSn     = 2
)

type User struct {
	WxId     string
	OpenId   string
	Status   int
	Sub      int    // 关注=1  未关注=0
	Value1   string
	Value2   string
	Mobile1  string
	Mobile2  string
	Mobile3  string
	NickName string // 昵称
	Sex      int    // 1=男   2=女
	ImgUrl   string // 头像
	Data string //   如果 手机号码 第一绑定 设置为 new ， 送完优惠后， 置空
}

// 传入 openid 从 redis 中获取 用户信息
// 如果redis中数据不存在则创建一个记录
func CreateUser(wxid, openid string) *User {

	rc := redispool.ConnPool.Get()
	defer rc.Close()

	find, _ := redis.Bool(rc.Do("EXISTS", openid))

	var u *User

	if find == true {

		smap, _ := redis.StringMap(rc.Do("HGETALL", openid))
		s, _ := strconv.Atoi(smap["status"])
		sub, _ := strconv.Atoi(smap["sub"])
		v1 := smap["value1"]
		v2 := smap["value2"]
		m1 := smap["mobile1"]
		m2 := smap["mobile2"]
		m3 := smap["mobile3"]
		nv, _ := strconv.Atoi(smap["sex"])
		nm := smap["nickname"]

		u = &User{
			WxId:     wxid,
			OpenId:   openid,
			Status:   s,
			Sub:      sub,
			Value1:   v1,
			Value2:   v2,
			Mobile1:  m1,
			Mobile2:  m2,
			Mobile3:  m3,
			NickName: nm,
			Sex:      nv,
			ImgUrl:   smap["imgurl"],
			Data: smap["data"],
		}
		log.Println("load user from redis", *u)

		if m1 != "" {
			rc.Do("HMSET", m1, "openid", openid)
			rc.Do("EXPIRE", m1, 315360000)
		}
		if m2 != "" {
			rc.Do("HMSET", m2, "openid", openid)
			rc.Do("EXPIRE", m2, 315360000)
		}
		if m3 != "" {
			rc.Do("HMSET", m3, "openid", openid)
			rc.Do("EXPIRE", m3, 315360000)
		}

	} else {
		u = &User{
			WxId:   wxid,
			OpenId: openid,
			Status: UserNormal, 
			Sub: 0, 
		}
	}

	return u
}

func (self *User) Save() {
	rc := redispool.ConnPool.Get()
	defer rc.Close()

	openid := self.OpenId

	rc.Do("HMSET", openid,
		"wxid", self.WxId,
		"status", self.Status,
		"sub", self.Sub, 
		"value1", self.Value1,
		"value2", self.Value2,
		"mobile1", self.Mobile1,
		"mobile2", self.Mobile2,
		"mobile3", self.Mobile3,
		"sex", self.Sex,
		"nickname", self.NickName,
		"imgurl", self.ImgUrl,
		"data", self.Data)
	rc.Do("EXPIRE", openid, 315360000) // 绑定的数据 永久保存


	m1 := self.Mobile1
	m2 := self.Mobile2
	m3 := self.Mobile3

	if m1 != "" {
		rc.Do("HMSET", m1, "openid", openid)
		rc.Do("EXPIRE", m1, 315360000)
	}
	if m2 != "" {
		rc.Do("HMSET", m2, "openid", openid)
		rc.Do("EXPIRE", m2, 315360000)
	}
	if m3 != "" {
		rc.Do("HMSET", m3, "openid", openid)
		rc.Do("EXPIRE", m3, 315360000)
	}
}

// 清除 指定手机号码的绑定 关系
// mobile 查 openid ， 然后 openid 中记录的 这个手机号码 清除
// 然后 mobile  key 删除 
func ClearMobile(wxid, mobile, opid string) bool {
	rc := redispool.ConnPool.Get()
	defer rc.Close()

	log.Println("ClearMobile", mobile)
	openid, _ := redis.String(rc.Do("HGET", mobile, "openid"))

	b := false 

	if openid == "" {
		openid = opid  // 如果通过手机号码 无法找到用户身份， 有可能是 当前微信用户要清理 手机号码
	}

	find, _ := redis.Bool(rc.Do("EXISTS", openid))
	if find {
		log.Println("ClearMobile", mobile, openid)
		smap, _ := redis.StringMap(rc.Do("HGETALL", openid))

		m1 := smap["mobile1"]
		m2 := smap["mobile2"]
		m3 := smap["mobile3"]

		if m3 == mobile {
			rc.Do("HMSET", openid, "mobile3", "")
			b = true
		}
		if m2 == mobile {
			rc.Do("HMSET", openid, "mobile2", m3)
			rc.Do("HMSET", openid, "mobile3", "")
			b = true
		}
		if m1 == mobile {
			rc.Do("HMSET", openid, "mobile1", m2)
			rc.Do("HMSET", openid, "mobile2", m3)
			rc.Do("HMSET", openid, "mobile3", "")
			b = true
		}

		if b {
			rc.Do("DEL", mobile)
		}
	}
	

	return b
}

// 通过 手机号码 反查 openid
func GetOpenidFromMobile(wxid, mobile string) string {
	rc := redispool.ConnPool.Get()
	defer rc.Close()

	openid, _ := redis.String(rc.Do("HGET", mobile, "openid"))
	return openid
}

// 是否存在绑定状态
func (self *User) IsDoBinding() bool {
	if self.Status == UserNormal {
		return false
	} else {
		return true
	}
}

func (self *User) IsWaitMobile() bool {
	if self.Status == UserBindWaitMobile {
		return true
	}
	return false
}

func (self *User) IsWaitSn() bool {
	if self.Status == UserBindWaitSn {
		return true
	}
	return false
}

func (self *User) HaveMobile() bool {
	if self.Mobile1 != "" || self.Mobile2 != "" || self.Mobile3 != "" {
		return true 
	}
	return false 
}

func (self *User) SendBindWelcom() {
	// 判断 是否存在已经绑定的号码
	if self.Mobile1 != "" || self.Mobile2 != "" || self.Mobile3 != "" {
		SendMessageText(self.WxId, self.OpenId, "已为您绑定号码："+
			self.Mobile1+" "+self.Mobile2+" "+self.Mobile3+
			"请输入您的新手机号码：")
	} else {
		SendMessageText(self.WxId, self.OpenId, "请粘贴云喇叭的快递通知短信，或直接输入您的手机号码：")
	}

	self.WaitMobile() // 状态转换到 等待手机号码
}

func (self *User) WaitMobile() {
	rc := redispool.ConnPool.Get()
	defer rc.Close()

	self.Status = UserBindWaitMobile

	openid := self.OpenId

	rc.Do("HMSET", openid,
		"status", UserBindWaitMobile,
		"value1", "",
		"value2", "")

	// 首次绑定  key 只保留 3分钟
	// if self.Mobile1 == "" && self.Mobile2 == "" && self.Mobile3 == "" {
	// 	rc.Do("EXPIRE", openid, 180) // 临时状态维持 2分钟
	// }

	return
}

// 退出 bind 状态
func (self *User) ExitBind() {
	rc := redispool.ConnPool.Get()
	defer rc.Close()

	self.Status = UserNormal

	openid := self.OpenId

	rc.Do("HMSET", openid,
		"status", UserNormal,
		"value1", "",
		"value2", "")

	//首次绑定  key 只保留 3分钟
	if self.Mobile1 == "" && self.Mobile2 == "" && self.Mobile3 == "" {
		rc.Do("EXPIRE", openid, 180) // 临时状态维持 2分钟
	}

	return
}

// 发送验证码
func (self *User) SendSmsSn(mob string) {

	// 判断 手机号码是否ok

	if len(mob) != 11 || mob[0:1] != "1" { // 号码不正确 忽略此次 号码录入
		self.ExitBind()
		return
	}

	// 判断手机号码 是否已经存在

	rand.Seed(time.Now().UnixNano())
	sn := strconv.Itoa(rand.Intn(9999)) 

	if true {
		m, _ := strconv.Atoi(mob[2:6]) 
		n, _ := strconv.Atoi(mob[5:])
		tmp  := strconv.FormatInt( 1989 + int64(m) + int64(n)*23 + (time.Now().Unix()/1000), 10)  //  10分钟 内 保持 不变
		sn = tmp[len(tmp)-4 :]
	}

	SendValidSn(mob, sn) // 发短信

	SendMessageText(self.WxId, self.OpenId, "短信验证码已经发送到您手机，请接收后输入：")

	self.WaitSn(mob, sn) // 转换状态
}

func (self *User) WaitSn(mob, sn string) {
	rc := redispool.ConnPool.Get()
	defer rc.Close()

	self.Status = UserBindWaitSn
	self.Value1 = mob
	self.Value2 = sn

	openid := self.OpenId

	rc.Do("HMSET", openid,
		"status", UserBindWaitSn,
		"value1", mob,
		"value2", sn)

	return
}

func (self *User) ValidSn(sn string) bool {
	if self.Value2 == sn {
		log.Println("Bind ok")
		self.BindClose(true) // 成功绑定
		return true

	} else {
		self.BindClose(false) // 重新绑定
		return false
	}
}

func (self *User) BindClose(success bool) {
	rc := redispool.ConnPool.Get()
	defer rc.Close()

	if success {
		self.Status = UserNormal

		openid := self.OpenId

		m1 := self.Value1
		m2 := self.Mobile1
		m3 := self.Mobile2
		d  := ""

		if m2=="" && m3==""{
			d = "new"  // 新绑定手机号码
		}

		rc.Do("HMSET", openid,
			"status", UserNormal,
			"sub", 1, 
			"value1", "",
			"value2", "",
			"mobile1", m1,
			"mobile2", m2,
			"mobile3", m3,
			"data", d )
		rc.Do("EXPIRE", openid, 315360000) // 绑定的数据 永久保存

		if m1 != "" {
			rc.Do("HMSET", m1, "openid", openid)
			rc.Do("EXPIRE", m1, 315360000)
		}
		if m2 != "" {
			rc.Do("HMSET", m2, "openid", openid)
			rc.Do("EXPIRE", m2, 315360000)
		}
		if m3 != "" {
			rc.Do("HMSET", m3, "openid", openid)
			rc.Do("EXPIRE", m3, 315360000)
		}

		self.Mobile1 = m1
		self.Mobile2 = m2
		self.Mobile3 = m3

		SendMessageText(self.WxId, self.OpenId, "已为您绑定手机："+self.Mobile1)


		// 关联信息 上报到 mq 
		var bind struct {
			Cmd          string `json:"cmd"`       // 扩展
			Mobile       string `json:"mobile"`    // 扩展
			OpenId       string `json:"openid"` 
		}
		bind.Cmd = "WechatMobileBind"
		bind.Mobile = self.Value1
		bind.OpenId = self.OpenId
					
		bs, _ := json.Marshal(bind)
		tmpjson := string(bs)
		// 收到的数据 上报给业务系统

		go func(){
			time.Sleep(50 * time.Millisecond)
			mq.Publish("wechat.out.exchange", "direct", "out.key", tmpjson, true)
			log.Println("userinfo json push to mq:", tmpjson)
		}()
		

	} else {
		SendMessageText(self.WxId, self.OpenId, "诶呀，验证码错误。请重新输入您的手机号码：")
		self.WaitMobile()
	}
	return
}

// 查询包裹   返回 文本信息   接口待优化
func (self *User) GetPackage() (arr []string) {

RTY:
	conn, err := g.GetDbConn("default")
	if err != nil {
		log.Println("[ERROR] get dbConn fail", err)
		time.Sleep(1 * time.Second)
		goto RTY
	}

	// courier_id   created_at
	tmpMob := "\"" + self.Mobile1 + "\""
	if self.Mobile2 != "" {
		tmpMob += ",\"" + self.Mobile2 + "\""
	}
	if self.Mobile3 != "" {
		tmpMob += ",\"" + self.Mobile3 + "\""
	}

	rows, err := conn.Query("select `from`, COALESCE(package_flowcode, '') as package_flowcode, place_name,created_at from send_logs where created_at > ? and `to` in ("+tmpMob+") order by id limit ? ", time.Now().Add(-3*24*time.Hour), 100)
	defer rows.Close()

	if err != nil {
		log.Println("[ERROR] get rows fail", err)
	}

	for rows.Next() {
		var from, code, at string
		var crtdt time.Time
		if e := rows.Scan(&from, &code, &at, &crtdt); e != nil {
			log.Println("[ERROR] get row fail", e)
		} else {
			crtdt8 := crtdt.Add(8 * time.Hour) // 时区
			tmp := "通知时间：" + crtdt8.Format("2006-01-02 15:04:05") + "\n"
			tmp += "快递"
			if code != "" {
				tmp += "编号" + code
			}
			tmp += at
			tmp += "电话" + from
			arr = append(arr, tmp) // 保存id 集合
		}
	}

	return arr
}
