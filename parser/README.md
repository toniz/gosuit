## 简介
将传入的字符传按照指定的格式，转换成为GO的数据结构。  
目前格式支持： json, xml, yaml  
添加格式支持，大部分格式只需要在parser.go头文件对应的库即可。  


## json解析:  
``` go
import . "hsb.com/parser"

type testJson struct{
    Name  string `json:"name"`
    Test  int    `json:"test"`
    Array []string `json:"array"`
}

bytes := []byte(`{"name":"boouoo","test":3,"array":["1","3","4"]}`)
var j testJson
err := Decode("json", bytes, &j)
```

## XML解析：  
``` go
var zoo struct {
    Animals []string `xml:"animal"`
}

bytes := []byte(`<animals>
                    <animal>gopher</animal>
                    <animal>armadillo</animal>
                    <animal>zebra</animal>
                    <animal>unknown</animal>
                    <animal>gopher</animal>
                    <animal>bee</animal>
                    <animal>gopher</animal>
                    <animal>zebra</animal>
                </animals>`)
err := Decode("xml", bytes, &zoo)
```

## yaml解析：
``` go
type StructA struct {
    A string `yaml:"a"`
}

type StructB struct {
    StructA `yaml:",inline"`
    B string `yaml:"b"`
}

bytes := []byte("a: a string from struct A\nb: a string from struct B")
var b StructB
err := Decode("yaml", bytes, &b)
```

### 使用例子可以参考  
[parse 测试用例](parse_test.go)  


