/*
 * Create By Xinwenjia 2018-04-15
 * Modify From-https://github.com/toniz/gudp
 */

package parser

import (
    "errors"
    glog "github.com/golang/glog"

    "strings"
    "encoding/json"
    "encoding/xml"
    "gopkg.in/yaml.v2"
)

func Decode(name string, text []byte, l interface{}) error {
    switch name {
    case "json", ".json":
    {
        t := string(text)
        t = strings.Replace(t, "\n", " ", -1)
        t = strings.Replace(t, "\r", " ", -1)
        err := json.Unmarshal([]byte(t), l)
        if err != nil {
            glog.Warningf("Decode Parse Failed[%v]: %v", err, t)
        }
        return err
    }
    case "xml", ".xml":
    {
        err := xml.Unmarshal(text, l)
        if err != nil {
            glog.Warningf("Decode Parse Failed[%v]: %v", err, text)
        }
        return err
    }
    case "yaml", ".yaml":
    {
        err := yaml.Unmarshal(text, l)
        if err != nil {
            glog.Warningf("Decode Parse Failed[%v]: %v", err, text)
        }
        return err
    }
    default:
    {
        glog.Warningf("Decode Parse No Decoder Being Add")
        return errors.New("No Decoder Being Add.")
    }
    }
}
