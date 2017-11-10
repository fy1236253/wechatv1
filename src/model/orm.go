package model

import (
	"database/sql"
	"g"
	"log"
)

func CreatNewUploadImg(uuid, openid string) {
	conn, _ := g.GetDBConn("default")
	defer conn.Close()
	stmt, _ := conn.Prepare("INSERT img_order SET uuid=?,openid=?")
	stmt.Exec(uuid, openid)
}

type ImgUuid struct {
	UUID string `json:"uuid"`
}

// GetUploadImgInfo 从数据库中提取需要人工的单子
func GetUploadImgInfo() (arr []string) {
	conn, _ := g.GetDBConn("default")
	var rows *sql.Rows
	rows, _ = conn.Query("select uuid from img_order")
	defer rows.Close()
	for rows.Next() {
		var uuid string
		if e := rows.Scan(&uuid); e != nil {
			log.Println("[ERROR] get row fail", e)
		} else {
			log.Println(uuid)
			arr = append(arr, uuid) // 保存id 集合
		}
	}
	return arr
}
