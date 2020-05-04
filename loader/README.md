## 简介
从数据源读取数据，并调用了parse模块, 转换成GO的数据结构。  
目前支持了从文件或ETCD加载数据。  

## 创建对象：
* 创建file loader对象：  
``` go
import (
    "hsb.com/loader"
    _ "hsb.com/loader/fileloader"
)

l, err := loader.NewLoader("file")
```

* 创建etcd loader对象：
``` go
import (
    "hsb.com/loader"
    _ "hsb.com/loader/etcdloader"
)
l, err := loader.NewLoader("etcd")
```

## 建立链接:  
* 链接etcd  
``` go
endpoint := "10.106.210.224:2379"
user := "root"
password := "e3jSlAsGNw"
err = l.Connect(endpoint, user, password)
```

* file不需要connect  


## 获取列表  
### file loader 参数说明  
* path: 文件路径
* ext: 文件后缀名
* prefix: 文件名前置

### etcd loader 参数说明
* path: namespace
* ext: 无效
* prefix: key前缀

``` go
list, err := l.GetList(p, ".json", "t*")
```


## 加载数据：
* 会根据key或者文件名的后缀来判断文件格式，并加载数据到对应的数据结构  
``` go
type testJson struct{
    Name  string `json:"name"`
    Test  int    `json:"test"`
    Array []string `json:"array"`
}
p := "example/test.json"
var result testJson
err = l.Load(p, &result)
```

### 使用例子可以参考  
[file loader 测试用例](fileloader/loader_test.go)   
[etcd loader 测试用例](etcdloader/loader_test.go)   


