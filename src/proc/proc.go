package proc

import (
	nproc "github.com/toolkits/proc"
	"log"

	"github.com/garyburd/redigo/redis"
	redispool "redis"
)

// 统计指标的整体数据
var (
	PVCnt = nproc.NewSCounterBase("PVCnt") // mq 数据包 计数
)

func Start() {
	log.Println("proc.Start, ok")

	//PVCnt.PutOther("Cnt", int64(0))
	rc := redispool.ConnPool.Get()
	defer rc.Close()

	cnt, _ := redis.Int64(rc.Do("HGET", "nproc", "PVCnt"))

	PVCnt.SetCnt(cnt)
}

func Stop() {
	log.Println("proc.Stop ok")

	rc := redispool.ConnPool.Get()
	defer rc.Close()

	cnt, _ := redis.Int64(rc.Do("HGET", "nproc", "PVCnt"))

	rc.Do("HMSET", "nproc", "PVCnt", PVCnt.Cnt+cnt) // 当前进程退出了， 需要把计数累加到 redis中
}

func GetAll() []interface{} {
	ret := make([]interface{}, 0)

	ret = append(ret, PVCnt.Get())

	return ret
}
