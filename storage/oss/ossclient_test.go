/*
 * Create By Xinwenjia 2020-02-09
 */

package ossclient_test

import (
    "context"
    "io"
    "time"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"

    . "hsb.com/storage"
    _ "hsb.com/storage/oss"
)

var _ = Describe("OSS Client", func() {
    var (
        err        error
        reader     io.Reader
        bucketName string
    )
    c, _ := NewStorageDriver("oss")

    endpoint := "http://oss-cn-beijing.aliyuncs.com"
    accessKeyID := "LTAI4FgJC35ZB9Sk7dZFJp6o"
    secretAccessKey := "ADZRavgKfD3NQMGFpkHECgHcVA1xJf"
    bucketName = "ibbwhat"

    Context("Test GetObjectList", func() {
        err = c.Connect(endpoint, accessKeyID, secretAccessKey)
        It("Should Return No Error", func() {
            Expect(err).NotTo(HaveOccurred())
        })
    })

    Describe("Test Aliyun OSS API", func() {
        Context("Test GetObjectList", func() {
            objectPrefix := "android"
            _, err = c.GetObjectList(bucketName, objectPrefix)
            It("Should Return No Error", func() {
                Expect(err).NotTo(HaveOccurred())
            })
        })

        Context("Test GetObject", func() {
            objectName := "android/20190920165652.jpg"
            ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
            defer cancel()
            reader, err = c.GetObject(ctx, bucketName, objectName)

            It("Should Return No Error", func() {
                Expect(err).NotTo(HaveOccurred())
            })
        })

        Context("Test PutObject", func() {
            objectName := "android/20190920165652_new.jpg"
            ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
            defer cancel()
            _, err = c.PutObject(ctx, bucketName, objectName, reader, -1)

            It("Should Return No Error", func() {
                Expect(err).NotTo(HaveOccurred())
            })
        })
    })
})
