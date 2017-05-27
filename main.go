package main

import (
	"cron"
	"flag"
	"fmt"
	"g"
	"http"
	"log"
	"os"
)

func main() {
	cfg := flag.String("c", "cfg.json", "specify config file")
	version := flag.Bool("v", false, "show version")
	flag.Parse()
	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}
	g.ParseConfig(*cfg) //配置文件
	g.InitWxConfig()    //微信相关参数
	log.Println(g.GetWechatConfig("123424"))
	g.InitDB()      //db池
	g.InitRootDir() //全局参数

	logTo := g.Config().Logs
	if logTo != "stdout" {
		f, err := os.OpenFile(logTo, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(fmt.Sprintf("open logfile error"))
		}
		defer f.Close()
		log.SetOutput(f)
	}
	// 日志追加pid和时间
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetPrefix(fmt.Sprintf("PID.%d ", os.Getpid()))

	go http.Start()
	go cron.Start()
}
