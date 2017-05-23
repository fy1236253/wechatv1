package main

import (
	"flag"
	"fmt"
	"g"
	"log"
	"net/http"
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
	log.Println(*cfg)
	g.ParseConfig(*cfg)
	g.InitWxConfig()
	// log.Println(g.GetWechatConfig("gh_8ac8a8821eb8"))
	// log.Println(g.Config().Debug)
	logTo := g.Config().Logs
	if logTo != "stdout" {
		f, err := os.OpenFile(logTo, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(fmt.Sprintf("open logfile error"))
		}
		defer f.Close()
		log.SetOutput(f)
	}
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetPrefix(fmt.Sprintf("PID.%d ", os.Getpid()))
	log.Println("nice")
	http.ListenAndServe(":8089", nil)
}
