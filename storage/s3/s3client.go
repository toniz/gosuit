/*
 * Create By Xinwenjia 2020-02-09
 */

package s3client

import (
	"context"
	"io"

	"github.com/minio/minio-go"
	"github.com/toniz/gosuit/storage"
)

type S3Client struct {
	cli *minio.Client
}

func init() {
	storage.Register("s3", func() storage.StorageDriver {
		return new(S3Client)
	})
}

// Create CpehS3 Client Handler
func (c *S3Client) Connect(endpoint string, accessKeyID string, secretAccessKey string) error {
	var err error
	c.cli, err = minio.NewV2(endpoint, accessKeyID, secretAccessKey, false)
	return err
}

// Create a done channel to control 'ListObjects' go routine.
// Indicate to our routine to exit cleanly upon return.
// List all objects from a bucket-name with a matching prefix.
func (c *S3Client) GetObjectList(bucketName string, objectPrefix string) (names []string, err error) {
	doneCh := make(chan struct{})
	defer close(doneCh)

	for object := range c.cli.ListObjectsV2(bucketName, objectPrefix, true, doneCh) {
		if object.Err != nil {
			err = object.Err
			return
		}
		names = append(names, object.Key)
	}

	return
}

// Uploads objects that are less than 128MiB in a single PUT operation.
// For objects that are greater than 128MiB in size, PutObject seamlessly
// uploads the object as parts of 128MiB or more depending on the actual file size.
// The max upload size for an object is 5TB.
func (c *S3Client) PutObject(ctx context.Context, bucketName string, objectName string, reader io.Reader, objectSize int64) (n int, err error) {
	//_, err = c.cli.PutObjectWithContext(ctx, bucketName, objectName, reader, objectSize, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	_, err = c.cli.PutObject(bucketName, objectName, reader, objectSize, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	return
}

// Returns a stream of the object data. Most of the common errors occur when reading the stream.
func (c *S3Client) GetObject(ctx context.Context, bucketName string, objectName string) (reader io.Reader, err error) {
	//reader, err = c.cli.GetObjectWithContext(ctx, bucketName, objectName, minio.GetObjectOptions{})
	reader, err = c.cli.GetObject(bucketName, objectName, minio.GetObjectOptions{})
	return
}
