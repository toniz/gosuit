package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"

	. "github.com/toniz/gosuit/Service/Storage/Worker"
	"github.com/toniz/gosuit/queue"
	_ "github.com/toniz/gosuit/queue/rabbitmq"
)

var _ = Describe("Test Storage Worker Service", func() {
	var (
		mqc MQConfig
		err error
		mq  queue.MessageQueuer
	)

	Context("Test Load Message Queue Configure", func() {
		It("Should Return No Error", func() {
			mqc, err = LoadMQConfig("./secret/mq_rabbitmq.json")
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("Test Load Storage Configure", func() {
		It("Should Return No Error", func() {
			_, err = LoadStorageConfig("./secret/", ".json", "storage*")
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("Test Connect Queue", func() {
		It("Should Return No Error", func() {
			mq, _ = queue.NewMessageQueue("rabbitmq")
			err = mq.Connect(mqc.Endpoint, mqc.User, mqc.Password)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("Test Send Task", func() {
		str := `{
            "Source":{
                "Owner":"s3",
                "Bucket":"image",
                "Prefix":"beimian/",
                "Name":"1560431791881-2.jpeg"
            },
            "Target":{
                "Owner":"cos",
                "Bucket":"ibbwhat",
                "Prefix":"beimian/",
                "Name":"1560431791881-storage.jpeg"
            }
        }`

		It("Should Return No Error", func() {
			err = mq.SendTask(mqc.QueueName, str)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("Test Send Task", func() {
		It("Should Return No Error", func() {
			err = mq.Worker(mqc.QueueName, MessageHandler)
			time.Sleep(time.Duration(3) * time.Second)
		})
	})
})
