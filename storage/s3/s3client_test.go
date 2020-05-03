/*
 * Create By Xinwenjia 2020-02-09
 */

package s3client_test

import (
    "context"
    "io"
    "time"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"

    . "hsb.com/storage"
    _ "hsb.com/storage/s3"
)

var _ = Describe("Ceph S3 Client", func() {
    var (
        err        error
        reader     io.Reader
        bucketName string
    )
    c, _ := NewStorageDriver("s3")

    Context("Test GetObjectList", func() {
        endpoint := "master01:7480"
        accessKeyID := "42II6092AF4I2OGA5TP9"
        secretAccessKey := "9MP7VuzkMFVpzIDzL5ueXubdB254RXgezQm5hN3W"
        bucketName = "image"
        //endpoint := "cos.ap-guangzhou.myqcloud.com"
        //accessKeyID := "AKIDtNkXxNaS2azJaprmcRxRnEfxVwKgt1wh"
        //secretAccessKey := "tsx0vgH39ezRllcFi8R5yX3RegAZDj4G"
        //bucketName = "toniz-1253750834"
        err = c.Connect(endpoint, accessKeyID, secretAccessKey)
        It("Should Return No Error", func() {
            Expect(err).NotTo(HaveOccurred())
        })
    })

    Describe("Test Ceph S3 API", func() {
        Context("Test GetObjectList", func() {
            objectPrefix := "test"
            _, err = c.GetObjectList(bucketName, objectPrefix)

            It("Should Return No Error", func() {
                Expect(err).NotTo(HaveOccurred())
            })
        })

        Context("Test GetObject", func() {
            objectName := "beimian/1560431791881-2.jpeg"
            ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
            defer cancel()
            reader, err = c.GetObject(ctx, bucketName, objectName)

            It("Should Return No Error", func() {
                Expect(err).NotTo(HaveOccurred())
            })
        })

        Context("Test PutObject", func() {
            objectName := "beimian/1560431791881-2_new.jpeg"
            ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
            defer cancel()
            _, err = c.PutObject(ctx, bucketName, objectName, reader, -1)

            It("Should Return No Error", func() {
                Expect(err).NotTo(HaveOccurred())
            })
        })
    })
})
