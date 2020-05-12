package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/toniz/gosuit/loader"
	_ "github.com/toniz/gosuit/loader/fileloader"
	"github.com/toniz/gosuit/parser"
	"github.com/toniz/gosuit/queue"
	"github.com/toniz/gosuit/storage"
	_ "github.com/toniz/gosuit/storage/cos"
	_ "github.com/toniz/gosuit/storage/s3"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

type Object struct {
	Owner  string
	Bucket string
	Prefix string
	Name   string
}

type Message struct {
	Source Object
	Target Object
}

type StorageConfig struct {
	Owner           string
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
}

type MQConfig struct {
	Owner     string
	Endpoint  string
	User      string
	Password  string
	QueueName string
}

var (
	stc map[string]StorageConfig
)

func LoadStorageConfig(p string, ext string, prefix string) (map[string]StorageConfig, error) {
	s := make(map[string]StorageConfig)
	l, err := loader.NewLoader("file")

	fileList, err := l.GetList(p, ext, prefix)
	if err != nil {
		return nil, err
	}

	for _, file := range fileList {
		var result StorageConfig
		err := l.Load(file, &result)
		if err != nil {
			fmt.Println("Load Json File Failed: ", err)
			continue
		}
		s[result.Owner] = result
	}

	//fmt.Println(s)
	return s, nil
}

func LoadMQConfig(p string) (MQConfig, error) {
	var result MQConfig
	l, err := loader.NewLoader("file")
	err = l.Load(p, &result)
	if err != nil {
		fmt.Println("Load Json File Failed: ", err)
	}

	//fmt.Println(result)
	return result, nil
}

func MessageHandler(s []byte) error {
	log.Println("Receive message: ", s)
	stc, err := LoadStorageConfig("./secret/", ".json", "storage*")
	if err != nil {
		log.Printf("Failed To Load Storage Configure: %s", err)
		return nil
	}

	var msg Message
	err = parser.Decode("json", s, &msg)
	if err != nil {
		log.Printf("Failed To Parse RabbitMQ Message: %s", err)
		return nil
	}

	srcOwner := msg.Source.Owner
	src, err := storage.NewStorageDriver(srcOwner)
	if err != nil {
		log.Printf("Failed To New Storage Source: %s", err)
		return nil
	}

	err = src.Connect(stc[srcOwner].Endpoint, stc[srcOwner].AccessKeyID, stc[srcOwner].SecretAccessKey)
	if err != nil {
		log.Printf("Failed To Connect Storage Source: %s", err)
		return nil
	}

	tagOwner := msg.Target.Owner
	tag, err := storage.NewStorageDriver(tagOwner)
	if err != nil {
		log.Printf("Failed To New Storage Target: %s", err)
		return nil
	}

	err = tag.Connect(stc[tagOwner].Endpoint, stc[tagOwner].AccessKeyID, stc[tagOwner].SecretAccessKey)
	if err != nil {
		log.Printf("Failed To Connect Storage Target: %s", err)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	reader, err := src.GetObject(ctx, msg.Source.Bucket, msg.Source.Prefix+msg.Source.Name)
	if err != nil {
		log.Printf("Failed To Get Object From Source: %s", err)
		return nil
	}

	_, err = tag.PutObject(ctx, msg.Target.Bucket, msg.Target.Prefix+msg.Target.Name, reader, -1)
	if err != nil {
		log.Printf("Failed To Put Object To Target: %s", err)
		return nil
	}

	return err
}

func main() {
	mqc, err := LoadMQConfig("./secret/mq_rabbitmq.json")
	failOnError(err, "Failed To Load Message Queue Configure")

	mq, err := queue.NewMessageQueue("rabbitmq")
	failOnError(err, "Failed To Connect RabbitMQ")

	err = mq.Connect(mqc.Endpoint, mqc.User, mqc.Password)
	failOnError(err, "Failed To Connect RabbitMQ")

	err = mq.Worker(mqc.QueueName, MessageHandler)
	failOnError(err, "Failed To Read From RabbitMQ")

	forever := make(chan bool)
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
