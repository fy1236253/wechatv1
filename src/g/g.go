package g

import (
	"fmt"
	"log"
	"runtime"
)

const (
	VERSION = "0.1.0"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", err, msg)
		panic(fmt.Sprintf("%s: %s", err, msg))
	}
}
