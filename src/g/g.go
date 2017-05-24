package g

import (
	"log"
	"os"
)

const (
	// VERSION 版本号
	VERSION = "wechatv1 0.1.0"
)

var Root string

func InitRootDir() {
	var err error
	Root, err = os.Getwd()
	if err != nil {
		log.Fatalln("getwd fail:", err)
	}
}
