/*
 * Create By Xinwenjia 2018-04-15
 * Modify From-https://github.com/toniz/gudp
 */

package etcdloader

import (
	"context"
	"fmt"
	"path"

	"github.com/toniz/gosuit/loader"
	"github.com/toniz/gosuit/parser"
	. "go.etcd.io/etcd/clientv3"
)

type EtcdLoader struct {
	cli *Client
}

func init() {
	loader.Register("etcd", func() loader.Loader {
		return new(EtcdLoader)
	})
}

func (f *EtcdLoader) Connect(connStr string, username string, password string) (err error) {
	f.cli, err = New(Config{Endpoints: []string{connStr}, Username: username, Password: password})
	return
}

// Get ETCD List.
// If ext param is not use.
// The prefix param to operate on the keys with matching prefix.
// For example, 'Get(foo, WithPrefix())' can return 'foo1', 'foo2', and so on.
func (f *EtcdLoader) GetList(namespace string, ext string, prefix string) ([]string, error) {
	fmt.Println("GetList From:", namespace, ext, prefix)
	var fileList []string
	// Add namespace separator
	if namespace[len(namespace)-1] != 47 {
		namespace = namespace + "/"
	}
	//cli.KV = namespace.NewKV(cli.KV, ns)
	resp, _ := f.cli.Get(context.TODO(), namespace+prefix, WithPrefix())
	for _, ev := range resp.Kvs {
		fileList = append(fileList, string(ev.Key))
	}

	return fileList, nil
}

// Read a key, And return key content by []byte
func (f *EtcdLoader) Read(key string) ([]byte, error) {
	var content []byte
	resp, err := f.cli.Get(context.TODO(), key)
	if err != nil {
		return content, nil
	}

	for _, ev := range resp.Kvs {
		content = ev.Value
		break
	}

	return content, nil
}

// Read File And Parse To Struct
func (f *EtcdLoader) Load(key string, l interface{}) error {
	fmt.Println("Load From:", key)
	ext := path.Ext(key)
	if len(ext) == 0 {
		ext = ".json"
	}

	resp, err := f.cli.Get(context.TODO(), key)
	if err != nil {
		return nil
	}

	for _, ev := range resp.Kvs {
		err = parser.Decode(ext, ev.Value, l)
		break
	}

	return err
}

// Read File And Parse To Struct
func (f *EtcdLoader) ReadAll(namespace string, ext string, prefix string) (map[string][]byte, error) {
	fmt.Println("ReadAll From:", namespace, ext, prefix)
	result := make(map[string][]byte)

	// Add namespace separator
	if namespace[len(namespace)-1] != 47 {
		namespace = namespace + "/"
	}

	//cli.KV = namespace.NewKV(cli.KV, ns)
	resp, err := f.cli.Get(context.TODO(), namespace+prefix, WithPrefix())
	if err != nil {
		return result, err
	}

	for _, ev := range resp.Kvs {
		result[string(ev.Key)] = ev.Value
	}

	return result, nil
}

// Read File And Parse To Struct
func (f *EtcdLoader) LoadAll(namespace string, ext string, prefix string, l interface{}) error {

	return nil
}
