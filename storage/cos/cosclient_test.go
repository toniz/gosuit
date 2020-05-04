/*
 * Create By Xinwenjia 2020-02-09
 */

package cosclient_test

import (
	"context"
	"io"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "hsb.com/storage"
	_ "hsb.com/storage/cos"
)

var _ = Describe("COS Client", func() {
	var (
		err        error
		reader     io.Reader
		bucketName string
	)
	c, _ := NewStorageDriver("cos")

	endpoint := "https://<bucket>.cos.ap-guangzhou.myqcloud.com"
	accessKeyID := "AKIDJ"
	secretAccessKey := "tsx0vg"
	bucketName = "toniz-12537"
	err = c.Connect(endpoint, accessKeyID, secretAccessKey)

	It("Should Return No Error", func() {
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Test Tencent COS API", func() {
		Context("Test GetObjectList", func() {
			objectPrefix := "test"
			_, err = c.GetObjectList(bucketName, objectPrefix)

			It("Should Return No Error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("Test GetObject", func() {
			objectName := "test.jpg"
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
			defer cancel()
			reader, err = c.GetObject(ctx, bucketName, objectName)

			It("Should Return No Error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("Test PutObject", func() {
			objectName := "test_new.jpg"
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
			defer cancel()
			_, err = c.PutObject(ctx, bucketName, objectName, reader, -1)

			It("Should Return No Error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
