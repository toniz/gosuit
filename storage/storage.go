/*
 * Create By Xinwenjia 2020-02-09
 */

package storage

import (
    "context"
    "io"
    "errors"
    "fmt"
)

type StorageDriver interface {
    Connect(endpoint string, accessKeyID string, secretAccessKey string) error
    GetObjectList(bucketName string, objectPrefix string) (names []string, err error)
    PutObject(ctx context.Context, bucketName string, objectName string, reader io.Reader, objectSize int64) (n int, err error)
    GetObject(ctx context.Context, bucketName string, objectName string) (reader io.Reader, err error)
}

var (
    storageDrivers = make(map[string]func() StorageDriver)
)

func Register(name string, driver func() StorageDriver) {
    storageDrivers[name] = driver
}

func NewStorageDriver(driverName string) (s StorageDriver, err error) {
    if f, ok := storageDrivers[driverName]; ok {
        s = f()
    } else {
        err = errors.New(fmt.Sprintf("Storage type %s not recognized!", driverName))
    }
    return
}

