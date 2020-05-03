/*
 * Create By Xinwenjia 2020-02-09
 */

package cosclient

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"

        "hsb.com/storage"
	"github.com/tencentyun/cos-go-sdk-v5"
)

type CosClient struct {
	endpoint  string
	secretID  string
	secretKey string
}

func init() {
    storage.Register("cos", func() storage.StorageDriver {
        return new(CosClient)
    })
}

// Create Tencent COS Client Handler
//   endpoint: https://<bucket>.cos.ap-guangzhou.myqcloud.com
//   accessKeyID: cos secretID
//   secretAccessKey: cos secretKey
func  (c *CosClient) Connect(endpoint string, accessKeyID string, secretAccessKey string)  error {
	c.endpoint = endpoint
	c.secretID = accessKeyID
	c.secretKey = secretAccessKey

	return nil
}

// List all Object from a bucket-name with a matching prefix.
func (c *CosClient) GetObjectList(bucketName string, objectPrefix string) (names []string, err error) {
	// Not Found API

	return
}

// Uploads Object Using Http
func (c *CosClient) PutObject(ctx context.Context, bucketName string, objectName string, reader io.Reader, objectSize int64) (n int, err error) {
	u, _ := url.Parse(strings.Replace(c.endpoint, "<bucket>", bucketName, 1))
	b := &cos.BaseURL{BucketURL: u}
	cli := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  c.secretID,
			SecretKey: c.secretKey,
		},
	})

	_, err = cli.Object.Put(context.Background(), objectName, reader, nil)

	return
}

// Download Object Using Http
func (c *CosClient) GetObject(ctx context.Context, bucketName string, objectName string) (reader io.Reader, err error) {
	u, _ := url.Parse(strings.Replace(c.endpoint, "<bucket>", bucketName, 1))
	b := &cos.BaseURL{BucketURL: u}
	cli := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  c.secretID,
			SecretKey: c.secretKey,
		},
	})

	// Download object into ReadCloser(). the body needs to be closed
	if res, e := cli.Object.Get(context.Background(), objectName, nil); e != nil {
		err = e
		return
	} else {
		reader = res.Body
	}
	return
}
