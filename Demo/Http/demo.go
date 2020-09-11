package main

import (
    "log"
    "net/http"
    "io/ioutil"
    "flag"
    "encoding/json"
    "context"
    "fmt"
    //"github.com/golang/glog"

    _ "github.com/go-sql-driver/mysql"
    _ "github.com/toniz/gosuit/queue/mqtt"
    . "github.com/toniz/gosuit/dbproxy"
    . "github.com/toniz/gosuit/queue"
    "github.com/toniz/gosuit/glog"
)

var (
    listenPort = flag.String("listen_port", "8080", "Http Server Listen Port")
    ds *DBProxy
    mq MessageQueuer
)

type Rule struct {
    RuleID   int
    Request  string
    Response string
    Check    string
    Remark   string
}


func MqttReplyServiceHandler(writer http.ResponseWriter, request *http.Request) {
    
    choke := make(chan string)        
    queueName := "test/task_queue"
    err := mq.Worker(queueName, func(s []byte) int{
        choke <- fmt.Sprintln(s)
        return 0
    })
       
    if err != nil {
        log.Printf("Failed To Create Worker: %s", err)
        return
    }

    // 自己发条信息测试下
    err = mq.SendTask(queueName, "Test Message Queue")
    incoming := <- choke
    glog.Info("RECEIVED MESSAGE: %s\n", incoming)

    response := map[string]interface{} {
        "ret": 0,
        "errcode": 0,
        "msg": "Success",
        "data": incoming,
    }

    t, _ := json.Marshal(response)
    writer.Write(t)
}

func HttpReplyServiceHandler(writer http.ResponseWriter, request *http.Request) {
    log.Printf("Echoing back request made to %s to client (%s)", request.URL.Path, request.RemoteAddr)
    _, err := ioutil.ReadAll(request.Body)
    if err != nil {
        log.Printf("Failed To Get Request Body: %s", err)
        return
    }

    ident := "GOIR_t_http_reply_rule_select_byid"
    params := map[string]string{"id": "1"}
    result, err := ds.AutoCommit(context.TODO(), ident, params)
    if err != nil {
        log.Printf("Failed To Read From Mysql: %s", err)
        return
    }

    response := map[string]interface{} {
        "ret": 0,
        "errcode": 0,
        "msg": "Success",
        "data": result,
    }

    t, _ := json.Marshal(response)
    writer.Write(t)
}

func main() {
    flag.Parse()
    defer glog.Flush()

    fmt.Printf("starting server, listening on port %s\n", *listenPort)

    ds = NewDBProxy()
    err := ds.AddDBHandleFromFile("config/db", ".json", "db_*")
    err = ds.AddProxySQLFromFile("config/sql", "json", "sql_*")
    if err != nil {
        log.Printf("Failed To Read DB Configure: %s", err)
        return
    }

    endpoint := "tcp://10.10.240.109:1883"
    mq, _ = NewMessageQueue("mqtt")
    err = mq.Connect(endpoint, "", "")
    if err != nil {
        log.Printf("Connect To Mqtt Failed: %s", err)
        return
    }
    http.HandleFunc("/reply/mqtt", MqttReplyServiceHandler)
    http.HandleFunc("/reply/http", HttpReplyServiceHandler)
    http.ListenAndServe(":"+*listenPort, nil)
}


