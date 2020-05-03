/*
 * Create By Xinwenjia 2020-02-09
 */

package ossclient

import (
    "context"
    "io"

    "hsb.com/storage"
    "github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OssClient struct {
    cli *oss.Client
}

func init() {
    storage.Register("oss", func() storage.StorageDriver {
        return new(OssClient)
    })
}


// Create OSS Client Handler
func (c *OssClient) Connect(endpoint string, accessKeyID string, secretAccessKey string) error {
    var err error
    c.cli, err = oss.New(endpoint, accessKeyID, secretAccessKey)
    return err
}

// List all objects from a bucket-name with a matching prefix.
func (c *OssClient) GetObjectList(bucketName string, objectPrefix string) (names []string, err error) {
    if lsRes, e := c.cli.ListBuckets(); e != nil {
        err = e
        return
    } else {
        for _, bucket := range lsRes.Buckets {
            names = append(names, bucket.Name)
        }
    }

    return
}

// Uploads Object
func (c *OssClient) PutObject(ctx context.Context, bucketName string, objectName string, reader io.Reader, objectSize int64) (n int, err error) {
    if bucket, e := c.cli.Bucket(bucketName); e != nil {
        err = e
        return
    } else {
        err = bucket.PutObject(objectName, reader)
    }

    return
}

// Returns a stream of the object data. Most of the common errors occur when reading the stream.
func (c *OssClient) GetObject(ctx context.Context, bucketName string, objectName string) (reader io.Reader, err error) {
    if bucket, e := c.cli.Bucket(bucketName); e != nil {
        err = e
        return
    } else {
        reader, err = bucket.GetObject(objectName)
    }

    return
}
