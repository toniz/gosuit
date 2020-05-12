## 简介
优点是切换使用不同的对象存储不需要修改代码，只需要改配置文件。  
几行代码就能很方便的在不同对象存储间同步文件。  
* 目前支持了阿里云oss, 腾讯云cos和ceph s3.  


## 创建对象：
* 创建cos对象：  
``` go
import (
    . "github.com/toniz/gosuit/storage"
    _ "github.com/toniz/gosuit/storage/cos"
)

c, _ := NewStorageDriver("cos")
```

* 创建oss对象:  
``` go
import (
    . "github.com/toniz/gosuit/storage"
    _ "github.com/toniz/gosuit/storage/oss"
)
c, _ := NewStorageDriver("oss")
```

* 创建ceph s3对象：  
``` go
import (
    . "github.com/toniz/gosuit/storage"
    _ "github.com/toniz/gosuit/storage/s3"
)
c, _ := NewStorageDriver("s3")
```

## 建立链接:

* 链接cos  
``` go
endpoint := "https://<bucket>.cos.ap-guangzhou.myqcloud.com"
accessKeyID := "AKID"
secretAccessKey := "tsx0vg"
err = c.Connect(endpoint, accessKeyID, secretAccessKey)
```

* 链接oss  
``` go
endpoint := "http://oss-cn-beijing.aliyuncs.com"
accessKeyID := "LTAIJ"
secretAccessKey := "ADZRvf"
err = c.Connect(endpoint, accessKeyID, secretAccessKey)
```

* 链接ceph s3(腾讯云也支持s3协议):  
```
endpoint := "master01:7480"
accessKeyID := "42II6092AF4I2OGA5TP9"
secretAccessKey := "9MP7VuzkMFVpzIDzL5ueXubdB254RXgezQm5hN3W"
//注释掉的使用s3链接腾讯云。
//endpoint := "cos.ap-guangzhou.myqcloud.com"
//accessKeyID := "AKt1wh"
//secretAccessKey := "tsx0vg"
err = c.Connect(endpoint, accessKeyID, secretAccessKey)
```

## 获取列表
* 对象存储上的文件列表:   
``` go
bucketName = "image"
objectPrefix := "test"
_, err = c.GetObjectList(bucketName, objectPrefix)
```

## 上传文件：
* 上传文件到对象存储:  
``` go
bucketName = "image"
objectName := "beimian/1560431791881-2_new.txt"
reader := strings.NewReader("abcde")
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
_, err = c.PutObject(ctx, bucketName, objectName, reader, -1)
```

## 下载文件：
``` go
var reader io.Reader
objectName := "beimian/1560431791881-2.jpeg"
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
reader, err = c.GetObject(ctx, bucketName, objectName)
```

## 同步cos文件到ceph s3：
``` go
// 初始化和建立链接 
src := NewStorageDriver("cos")
src.Connect(endpoint, accessKeyID, secretAccessKey)
tag := NewStorageDriver("s3")
tag.Connect(endpoint, accessKeyID, secretAccessKey)

ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

// 下载到io.reader
reader, err := src.GetObject(ctx, msg.Source.Bucket, msg.Source.Prefix + msg.Source.Name)
if err != nil {
    log.Printf("Failed To Get Object From Source: %s", err)
    return nil
}

// 上传这个io.reader
_, err = tag.PutObject(ctx, msg.Target.Bucket, msg.Target.Prefix + msg.Target.Name, reader, -1)
if err != nil {
    log.Printf("Failed To Put Object To Target: %s", err)
    return nil
}

```

### 使用例子可以参考  
[cos 测试用例](cos/cosclient_test.go)  
[oss 测试用例](oss/ossclient_test.go)  
[s3 测试用例](s3/s3client_test.go)  


