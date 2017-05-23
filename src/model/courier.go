package model

import (
    "database/sql"
    "g"
    "log"
    "strconv"
    _ "github.com/go-sql-driver/mysql"
    "github.com/toolkits/net/httplib"
    "encoding/json"
    //"time"
    //"bytes"
    "github.com/garyburd/redigo/redis"
    redispool "redis"
)

type Courier struct {
    Courier_name    string `json:"name,omitempty"` //快递员编号
    Courier_phone   int `json:"phone,omitempty"`//geohash值
    Courier_company    string `json:"company,omitempty"` //快递公司编号
}

func GetCouriesId(geo string)(arr []int) {
    DB, err := g.GetDbConn("default")
    checkErr(err)
    var rows *sql.Rows
    if geo != "" {
        rows, err = DB.Query("select courier_id from courier_geos,couriers where courier_geos.geo_key=? and courier_geos.courier_id=couriers.id order by couriers.updated_at desc limit 10", geo)
        checkErr(err)
    }
    defer rows.Close()
    for rows.Next() {
        var courier_id int
        if e := rows.Scan(&courier_id); e != nil {
            log.Println("[ERROR] get row fail", e)
        }

        arr = append(arr, courier_id) // 保存id 集合
    }
    return arr

}

func GetCouriesCompany(express_company_id string)(company_name string) {
    DB, err := g.GetDbConn("default")
    checkErr(err)
    var rows *sql.Rows 
    rows, err = DB.Query("select name from express_companies where id=?", express_company_id)
    checkErr(err)
    defer rows.Close() 
    for rows.Next(){
        var name string
        if e := rows.Scan(&name); e != nil {
            log.Println("[ERROR] get row fail", e)  
        } 
    company_name = name

    }
    return company_name
}

func GetCouries(geo string) (arrn []*Courier) {
    DB, err := g.GetDbConn("default")
    checkErr(err)
    //defer DB.Close()
    var rows *sql.Rows
    var  arr []int
    arr = GetCouriesId(geo)
    for i:=0;i<len(arr);i++{
        rows, err = DB.Query("select name,phone,express_company_id from couriers where id=?", arr[i])
        checkErr(err)
        
        for rows.Next(){
            var phone               int
            var name                string
            var express_company_id  string

            if e := rows.Scan(&name,&phone,&express_company_id); e != nil {
                log.Println("[ERROR] get row fail", e)
            }
            
            company_name := GetCouriesCompany(express_company_id)

            tmp := &Courier{
                Courier_name:       name,
                Courier_phone:      phone,
                Courier_company:    company_name,
            }
            arrn = append(arrn, tmp)
            
        }
        rows.Close()
    }

    return arrn
}

func checkErr(err error) {
    if err != nil {
        log.Println("ERROR",err)
        //panic(err)
    }
}

//动态查询快递
type Msg struct {
    COMCODE   string `json:"comcode"`
    ID        string `json:"id"`
    NOCOUNT   string `json:"nocount"`
    NOPRE     string `json:"nopre"`
    STARTTIME string `json:"starttime"`
}


type ExpressInfo struct {
    DATA []*ExpressData `json:"data"`
}

type ExpressData struct {
    TIME    string `json:"time"`
    CONTEXT string `json:"context"`
}

type PMessage struct {
    Msg  string `json:"msg"`
    NUM  string `json:"num"`
    TIME string `json:"time"`
}

func GetCompany(num string) (comcode string) {
    //num := "3101066983900"
    incompleteURL := "http://m.kuaidi100.com/autonumber/auto?num=" + num
    req := httplib.Get(incompleteURL)
    resp, err := req.String()
    //log.Println(resp)
    if err != nil {
        log.Println("[ERROR]", err)
    }
    var result []*Msg
    err = json.Unmarshal([]byte(resp), &result)
    //log.Println(result)
    for _, v := range result {
        comcode = v.COMCODE
    }
    return comcode
}

func GetCouryInfo(num string, openid string) (re []*ExpressData){
    //num := "3101066983900"
    company := GetCompany(num)
    compUrl := "http://m.kuaidi100.com/query?type=" + company + "&postid=" + num
    req := httplib.Get(compUrl)
    resp, err := req.String()
    //log.Println(resp)

    if err != nil {
        log.Println("[ERROR]", err)
    }

    var result ExpressInfo
    err = json.Unmarshal([]byte(resp), &result)

    //log.Println(result.DATAS)
    re = result.DATA
    amount := len(re)
    if amount == 0 {
        log.Println("没有查询到信息")
    }
    //log.Println(amount)
    rc := redispool.ConnPoolLocalNet.Get()
    defer rc.Close()
    find, _ := redis.Bool(rc.Do("EXISTS", num))
    
    if find == true {
        rc.Do("HMSET",num,
        "amount",strconv.Itoa(amount), 
        "infos",resp)
        rc.Do("EXPIRE", num, 604800)//保存七天的单号信息
        log.Println("多次查询，不入队列")
    }else{
        log.Println("new boy")
        rc.Do("HMSET",num,
        "amount",strconv.Itoa(amount),
        "infos",resp,
        "openid",openid,
        "company",company)
        rc.Do("lpush","allnum",num )  //队列暂时一直保存
    }

    return 
}
