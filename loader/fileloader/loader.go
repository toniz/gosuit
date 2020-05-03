/*
 * Create By Xinwenjia 2018-04-15
 * Modify From-https://github.com/toniz/gudp
 */

package fileloader

import (
    "os"
    "path"
    "io/ioutil"

    "hsb.com/loader"
    "hsb.com/parser"
)

type FileLoader struct {
}

func init() {
    loader.Register("file", func() loader.Loader {
        return new(FileLoader)
    })
}

func  (f *FileLoader) Connect(name string, user string, password string) error {
    return nil
}

// Get File List. 
// If ext param is nil, get all file within Path
// The perfix param is shell pattern. eg: 'f*' match 'files.t' 'file.txt'.
func  (f *FileLoader) GetList(p string, ext string, prefix string) ([]string, error) {
    s, err := os.Stat(p)
    if err != nil {
        return nil, err
    }

    var fileList []string
    if s.IsDir() {
        files, err := ioutil.ReadDir(p)
        if err != nil {
            return nil, err
        }

        // Extension Name Add '.'
        if c := ext[0:1]; c != "." {
            ext = "." + ext
        }

        // Filter Out Subdirectory
        for _, sf := range files {
            if len(prefix) > 0 {
                if ok, _ := path.Match(prefix, sf.Name()); !ok {
                    continue
                }
            }

            if !sf.IsDir() {
                if len(ext) == 0  {
                    fileList = append(fileList, path.Clean(p) +"/" + sf.Name())
                } else {
                    if sfext := path.Ext(sf.Name()); sfext == ext {
                        fileList = append(fileList, path.Clean(p) +"/" + sf.Name())
                    }
                }
            }
        }
    } else {
        fileList = append(fileList, p )
    }
    return fileList, nil
}

// Read a file, And return file content by []byte
func  (f *FileLoader) Read(file string) ([]byte, error) {
    s, err := os.Stat(file)
    if err != nil {
        return nil, err
    }

    var content []byte
    if !s.IsDir() {
        fileHandler, err := os.Open(file)
        if err != nil {
            return nil, err
        }

        defer fileHandler.Close()
        content, err = ioutil.ReadAll(fileHandler)
    }
    return content, err
}

// Read File And Parse To Struct
func (f *FileLoader) Load(file string, l interface{}) error {
    c, err := f.Read(file)
    if err != nil {
        return err
    }

    err = parser.Decode(path.Ext(file), c, l)
    //fmt.Println("Loading Configure From Json File: ", l)
    return err
}

func (f *FileLoader) LoadAll(p string, ext string, prefix string, l interface{}) error {

    return nil
}

func (f *FileLoader) ReadAll(p string, ext string, prefix string) (map[string][]byte, error) {
    result := make(map[string][]byte)

    s, err := os.Stat(p)
    if err != nil {
        return nil, err
    }

    if s.IsDir() {
        files, err := ioutil.ReadDir(p)
        if err != nil {
            return nil, err
        }

        // Extension Name Add '.'
        if c := ext[0:1]; c != "." {
            ext = "." + ext
        }

        // Filter Out Subdirectory
        for _, sf := range files {
            if len(prefix) > 0 {
                if ok, _ := path.Match(prefix, sf.Name()); !ok {
                    continue
                }
            }

            if !sf.IsDir() {
                if len(ext) == 0  {
                    c, err := f.Read(path.Clean(p) +"/" + sf.Name())
                    if err != nil {
                        return nil, err
                    }
                    result[path.Clean(p) +"/" + sf.Name()] = c
                } else {
                    if sfext := path.Ext(sf.Name()); sfext == ext {
                         c, err := f.Read(path.Clean(p) +"/" + sf.Name())
                         if err != nil {
                             return nil, err
                         }
                         result[path.Clean(p) +"/" + sf.Name()] = c
                    }
                }
            }
        }
    } else {
        c, err := f.Read(p)
        if err != nil {
            return nil, err
        }
        result[p] = c
    }
    return result, nil
}


