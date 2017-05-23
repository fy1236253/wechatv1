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
	log.Println(g.GetWechatConfig("gh_8ac8a8821eb8"))
	// log.Println(g.Config().Debug)
	http.ListenAndServe(":8089", nil)
}
