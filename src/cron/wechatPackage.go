package cron

import (
    "g"
    "log"
    "strconv"
    _ "github.com/go-sql-driver/mysql"
    "github.com/toolkits/net/httplib"
    "encoding/json"
    "time"
    "model"
    "bytes"
    "github.com/garyburd/redigo/redis"
    redispool "redis"
)

func StartPackage() {
    go monitorPackage()
}

func monitorPackage() {
    for{
        wxid := "gh_f353e8a82fe5"
        //wxcfg = g.GetWechatConfig(wxid)
        rc := redispool.ConnPoolLocalNet.Get()
        defer rc.Close()
        ss,_ := redis.Int(rc.Do("LLEN","allnum")) 
        for i := 0; i < ss; i++ {
            tnum,_:= rc.Do("RPOP","allnum")
            querynum := string(tnum.([]byte))
            log.Println(querynum)
            find, _ := redis.Bool(rc.Do("EXISTS", querynum))
            if find==true {
                log.Println("nice boy")
                smap, _ := redis.StringMap(rc.Do("HGETALL", querynum))
                company := smap["company"]
                amount_last := smap["amount"]
                openid := smap["openid"]
                infos := smap["infos"]
                log.Println(openid)

                var result_ss model.ExpressInfo
                err := json.Unmarshal([]byte(infos), &result_ss)
                if err != nil {
                    log.Println("[ERROR]", err)
                }
                amounts, err := strconv.Atoi(amount_last)
                compUrl := "http://m.kuaidi100.com/query?type=" + company + "&postid=" + querynum
                req := httplib.Get(compUrl)
                resp, err := req.String()
                var result model.ExpressInfo
                err = json.Unmarshal([]byte(resp), &result)
                re := result.DATA
                amount := len(re)
                //log.Println(re)
                //var msgjson Message
                var msg model.PMessage
                if amount>amounts {
                    rc.Do("HMSET",querynum,
                    "amount",strconv.Itoa(amount), 
                    "infos",resp)
                    log.Println("new package message")
                    for k, v := range result_ss.DATA {
                        if k==0{
                            msg.Msg = v.CONTEXT
                            msg.TIME = v.TIME
                        }
                    }
                    msg.NUM = querynum
                    buf := bytes.NewBuffer(make([]byte, 0, 16<<10))
                    buf.Reset()
                    json.NewEncoder(buf).Encode(msg)
                    tmpjson := buf.String()
                    //log.Println(tmpjson)
                    model.SendSmsNotifies(openid, g.GetWechatAccessToken(wxid), tmpjson, "")
                    return
                }
                rc.Do("lpush","allnum",querynum)
            }else{
                log.Println("bad boy")

            }
        }
        time.Sleep(12 * time.Hour)
        continue
    }
}
