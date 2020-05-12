package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	//   rb "github.com/toniz/gosuit/mq/rabbitmq"
	. "github.com/toniz/gosuit/storage/oss"
	. "ibbwhat.com/mysql"
)

func handleError(err error) {
	fmt.Println(err)
	os.Exit(-1)
}

const (
	// RabbitMQ configure
	endpointRB  = "10.108.197.89:5672"
	userRB      = "user"
	passwordRB  = "AMeXcxWc4d"
	queueNameRB = "image_storage"

	// aliyun oss configure
	endpointOss        = "http://oss-cn-beijing.aliyuncs.com"
	accessKeyIDOss     = "LTAI"
	secretAccessKeyOss = "ADZRa"
)

var (
	listenPort = flag.String("listen_port", "8080", "Http Server Listen Port")
	mds        *MysqlDataService
)

type object struct {
	Owner  string
	Bucket string
	Prefix string
	Name   string
}

type Message struct {
	Source object
	Target object
}

type resJson struct {
	Ret     string `json:"ret"`
	ErrCode string `json:"errCode"`
	ErrStr  string `json:"errStr"`
}

func FileServiceHandler(writer http.ResponseWriter, request *http.Request) {
	log.Printf("Echoing back request made to %s to client (%s)", request.URL.Path, request.RemoteAddr)
	b, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Printf("Failed To Get Request Body: %s", err)
		return
	}

	var rq reqJson
	json.Unmarshal(b, &rq)

	ident = "STORAGE_t_storage_sync_insert"
	params = map[string]string{"srcOwner": rq.Source.Owner,
		"srcBucket": rq.Source.Bucket,
		"srcPrefix": rq.Source.Prefix,
		"srcName":   rq.Source.Name,
		"tagOwner":  rq.Target.Owner,
		"tagBucket": rq.Target.Bucket,
		"tagPrefix": rq.Target.Prefix,
		"tagName":   rq.Target.Name}

	_, err := mds.AutoCommit(context.TODO(), ident, params)

	t, _ := json.Marshal(resp.Credentials)
	writer.Write(t)
}

func OssDirServiceHandler(writer http.ResponseWriter, request *http.Request) {
	log.Printf("Echoing back request made to %s to client (%s)", request.URL.Path, request.RemoteAddr)
	b, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Printf("Failed To Get Request Body: %s", err)
		return
	}

	var rq reqJson
	err := json.Unmarshal(b, &rq)
	if err != nil {
		log.Printf("Failed To Parse RabbitMQ Message: %s", err)
		return
	}

	src, err := NewStorage(msg.Source)
	if err != nil {
		log.Printf("Failed To New Storage Source: %s", err)
		return
	}

	names, err := c.GetObjectList(bucketName, objectPrefix)
	if err != nil {
		log.Printf("Failed To Get Object List From Storage Source: %s", err)
		return
	}

	if len(names) < 1 {
		log.Printf("Object List Size: %d", len(names))
		return
	}

	var paramString string
	for i, name := range names {
		if i != 0 {
			paramString = "," + paramString
		}

		paramString = paramString + "(\"" + rq.Source.Owner + "\", "
		paramString = paramString + "(\"" + rq.Source.Bucket + "\", "
		paramString = paramString + "(\"" + rq.Source.Prefix + "\", "
		paramString = paramString + "(\"" + rq.Source.Name + "\", "
		paramString = paramString + "(\"" + rq.Target.Owner + "\", "
		paramString = paramString + "(\"" + rq.Target.Bucket + "\", "
		paramString = paramString + "(\"" + rq.Target.Prefix + "\", "
		paramString = paramString + "(\"" + rq.Target.Name + "\", NOW(), NOW(), 0, 0) "
	}

	ident = "STORAGE_t_storage_sync_insert_multi"
	params = map[string]string{"values": paramString}
	_, err := mds.AutoCommit(context.TODO(), ident, params)

	t, _ := json.Marshal(resp.Credentials)
	writer.Write(t)
}

func taskLoader() {
	for {
		log.Println("Loading Task From DB.")
		ident = "STORAGE_t_storage_sync_select_unsync"
		params = map[string]string{}
		res, err := mds.AutoCommit(context.TODO(), ident, params)
		if err != nil {
			log.Printf("Loading Task From DB Failed: %v", err)
		}

		if len(res) > 0 {
			//mq, err := rb.NewRabbitMQ(endpointRB, userRB, passwordRB)
			//if err != nil {
			//    log.Printf("Connect To RabbitMQ Failed: %v", err)
			//}
		}
		log.Println(res)
		time.Sleep(time.Second * 10)
	}
}

func main() {
	flag.Parse()
	fmt.Printf("starting server, listening on port %s\n", *listenPort)

	mds = NewMysqlDataService()
	mds.Cfg.DBPath = "./config/db/"
	mds.Cfg.SqlPath = "./config/sql/"
	err := mds.InitMysqlConnection()
	if err != nil {
		log.Fatalf("Connect To Mysql Failed:", err)
	}

	http.HandleFunc("/file", FileServiceHandler)
	http.HandleFunc("/dir/oss", OssDirServiceHandler)
	//http.HandleFunc("/sts", StsDefaultHandler)
	http.ListenAndServe(":"+*listenPort, nil)
}
