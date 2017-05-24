package http

import (
	"g"
	"log"
	"net/http"
	"path/filepath"
)

func Start() {
	// 静态资源请求
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.FileServer(http.Dir(filepath.Join(g.Root, "/public"))).ServeHTTP(w, r)
	})
	WebHTTP()
	// start http server
	addr := g.Config().HTTP.Listen

	s := &http.Server{
		Addr:           addr,
		MaxHeaderBytes: 1 << 30,
	}

	log.Println("http.Start ok, listening on", addr)
	log.Fatalln(s.ListenAndServe())
}
