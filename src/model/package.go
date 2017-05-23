package model

import (
	"g"
	"log"
	"time"
	//"math/rand"
	//"strconv"
	"database/sql"
	"encoding/json"
)

type Package struct {
	Uuid        string // 通知记录的唯一标识
	Sn          string
	To          string
	Userdata    string
	Name        string
	CompanyName string
	Note        string
	NotifyTime  time.Time // 通知时间
}

// 返回 某人的  最近 48 小时的 包裹记录 ， 从 redis中获取
func GetPackages(u *User, uuid string) (arr []*Package) {
	nTry := 0
RTY:
	nTry += 1
	conn, err := g.GetDbConn("default")
	if err != nil {
		log.Println("[ERROR] get dbConn fail", err)
		time.Sleep(1 * time.Second)
		if nTry > 3 {
			return
		}
		goto RTY
	}

	// courier_id   created_at
	tmpMob := "\"" + u.Mobile1 + "\""
	if u.Mobile2 != "" {
		tmpMob += ",\"" + u.Mobile2 + "\""
	}
	if u.Mobile3 != "" {
		tmpMob += ",\"" + u.Mobile3 + "\""
	}

	var rows *sql.Rows

	if uuid != "" {
		rows, err = conn.Query("select  message, created_at from send_logs where uuid = ? and `to` in ("+tmpMob+") ", uuid)
	} else {
		rows, err = conn.Query("select  message,created_at from send_logs where created_at > ? and `to` in ("+tmpMob+") order by id limit ? ", time.Now().Add(-2*24*time.Hour), 1)
	}

	defer rows.Close()

	if err != nil {
		log.Println("[ERROR] get rows fail", err)
		return
	}

	for rows.Next() {
		var message string
		var crtdt time.Time
		if e := rows.Scan(&message, &crtdt); e != nil {
			log.Println("[ERROR] get row fail", e)
		} else {
			crtdt8 := crtdt.Add(8 * time.Hour) // 时区
			//tmp := "通知时间：" + crtdt8.Format("2006-01-02 15:04:05") + "\n"
			//tmp += "快递"
			//if code != "" {
			//	tmp += "编号" + code
			//}
			//tmp += at
			//tmp +=  "电话" + from

			var d SysImJson

			e := json.Unmarshal([]byte(message), &d)
			if e != nil {
				log.Println("json 解析失败: %s", e)
				continue
			}

			tmp := &Package{
				Uuid:        d.Uuid,
				Sn:          d.Sn,
				To:          d.To,
				Userdata:    d.Userdata,
				Note:        d.Note,
				NotifyTime:  crtdt8,
				CompanyName: d.CompanyName,
				Name:        d.Name,
			}
			arr = append(arr, tmp) // 保存id 集合
		}
	}

	return arr
}
