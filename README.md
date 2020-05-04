# gosuit
编写该项目的目的是为了适配多种开源组件，减少业务代码复杂度, 并且在切换同类组件时候，不需要修改代码。  

### 目前支持的组件
* ***数据库：*** 目前支持mysql和pgsql,并且很容易添加其它数据库支持。从mysql切换到pgsql，只需要修改对应的配置文件即可。  
* ***对象存储：*** 目前支持cos, oss和s3，存储之间同步或者切换不需要修改框架代码.  
* ***队列：*** 目前支持kafka和rabbitmq。


### 使用
[dbproxy](dbproxy/README.md)  
[queue](queue/README.md)  
[storage](storage/README.md)  
[loader](loader/README.md)  
[parser](parse/README.md)  




