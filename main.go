package main

import (
	"cron"
	"flag"
	"fmt"
	"g"
	"http"
	"log"
	"mq"
	"os"
	"os/signal"
	"proc"
	"redis"
	"syscall"
)

func main() {
	cfg := flag.String("c", "cfg.json", "specify config file")
	version := flag.Bool("v", false, "show version")
	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	// config
	g.ParseConfig(*cfg)
	g.InitRootDir()
	g.InitWxConfig() // 初始化微信服务号的配置信息

	// //  log.Println("debug", g.GetWechatConfig("gh_8ac8a8821eb9").AppId  )

	log.Println("starting...")

	logTo := g.Config().Logs
	if logTo != "stdout" {
		f, err := os.OpenFile(logTo, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			g.FailOnError(err, "error opening file")
		}
		defer f.Close()
		log.SetOutput(f)
	}
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetPrefix(fmt.Sprintf("PID.%d ", os.Getpid()))

	//  模块开始加载

	g.InitDB()
	redis.InitConnPool() // redis 连接池
	mq.InitConnPool()    // mq 的连接池

	// proc 统计模块启动
	go proc.Start()

	// http  接口启动
	go http.Start()

	// workder 启动
	go cron.Start()

	// 停止 处理
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		log.Println("all service stopping...")

		cron.Stop()
		http.Stop()
		proc.Stop()

		mq.ConnPool.Close() // 关闭连接池
		redis.ConnPool.Close()

		log.Println("all service stop ok ")
		os.Exit(0)
	}()

	select {}
}
