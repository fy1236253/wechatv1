package model

import (
	"bufio"
	"database/sql"
	"g"
	"io"
	"log"
	"os"
	"strings"

	"github.com/toolkits/file"
)

func CreatNewUploadImg(uuid, openid string) {
	conn, _ := g.GetDBConn("default")
	stmt, err := conn.Prepare("INSERT img_order SET uuid=?,openid=?")
	log.Println(err)
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
	if rows == nil {
		log.Println("rows nil ")
		return
	}
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

func ImportDatbase() {
	conn, _ := g.GetDBConn("default")
	f, _ := os.Open("data/m.txt")
	r := bufio.NewReader(f)
	for {
		line, err := file.ReadLine(r)
		if err == io.EOF {
			break
		}
		str := string(line)
		str = strings.Trim(str, " ")
		log.Println(str)
		fields := strings.Split(str, ",")
		stmt, _ := conn.Prepare("INSERT medicine_info SET name=?,province=?,city=?,origin=?,address=?,method=?")
		_, e := stmt.Exec(fields[0], fields[1], fields[2], fields[3], fields[4], fields[5])
		log.Println(e)
	}
}
