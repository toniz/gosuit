/*
 * Create By Xinwenjia 2018-04-15
 * Modify From-https://github.com/toniz/gudp
 */

package parser

import (
    "errors"

    "encoding/json"
    "encoding/xml"
    "gopkg.in/yaml.v2"
)

func Decode(name string, text []byte, l interface{}) error {
    switch name {
    case "json", ".json":
    {
        text = strings.Replace(text, "\n", " ", -1)
        text = strings.Replace(text, "\r", " ", -1)
        err := json.Unmarshal(text, l)
        return err
    }
    case "xml", ".xml":
    {
        err := xml.Unmarshal(text, l)
        return err
    }
    case "yaml", ".yaml":
    {
        err := yaml.Unmarshal(text, l)
        return err
    }
    default:
    {
        return errors.New("No Decoder Being Add.")
    }
    }
}
