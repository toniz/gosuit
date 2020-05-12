/*
 * Create By Xinwenjia 2018-04-25
 */

package kafka_test

import (
    "time"
    "bytes"
    "log"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"

    . "hsb.com/queue"
    _ "hsb.com/queue/kafka"
)

var _ = Describe("Kafka Test", func() {
    queueName := "topic-test"
    mq, _ := NewMessageQueue("kafka")
    Context("Test New Connection", func() {
        endpoint := "10.96.16.9:9092"
        user := ""
        password := ""
        It("Should Return No Error", func() {
            err := mq.Connect(endpoint, user, password)
            Expect(err).NotTo(HaveOccurred())
        })
    })

    Describe("Test Kafka API", func() {
        Context("Test New Task", func() {
            msg := "Test Kafka Task Send"
            It("Should Return No Error", func() {
                err := mq.SendTask(queueName, msg)
                Expect(err).NotTo(HaveOccurred())
            })
        })

        Context("Test Worker", func() {
            It("Should Return No Error", func() {
                err := mq.Worker(queueName, func (s []byte) int{
                    log.Printf("Received a message: %s", s)
                    dot_count := bytes.Count(s, []byte("."))
                    t := time.Duration(dot_count)
                    time.Sleep(t * time.Second)
                    log.Printf("Done")
                    return 0
                })

                Expect(err).NotTo(HaveOccurred())
            })
        })
    })
})
