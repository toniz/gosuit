/*
 * Create By Xinwenjia 2018-04-15
 * Modify From-https://github.com/toniz/gudp
 */

package loader

import (
    "errors"
    "fmt"
)

type Loader interface {
    Connect(name string, username string, password string) error
    GetList(p string, ext string, prefix string) ([]string, error)
    Read(file string) ([]byte, error)
    ReadAll(p string, ext string, prefix string) (map[string][]byte, error)
    Load(file string, l interface{}) error
    LoadAll(p string, ext string, prefix string, l interface{}) error
}

var (
    loaders = make(map[string]func() Loader)
)

func Register(name string, l func() Loader) {
    loaders[name] = l
}

func NewLoader(t string) (s Loader, err error) {
    if f, ok := loaders[t]; ok {
        s = f()
    } else {
        err = errors.New(fmt.Sprintf("Loader type %s not recognized!", t))
    }
    return
}

