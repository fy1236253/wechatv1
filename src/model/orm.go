package model

import (
	"g"
)

func CreatNewUploadImg(uuid, openid string) {
	conn, _ := g.GetDBConn("default")
	stmt, _ := conn.Prepare("INSERT img_order SET uuid=?,openid=?")
	stmt.Exec(uuid, openid)
}
