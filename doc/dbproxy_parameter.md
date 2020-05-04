# SQL File Parameter Instructions

## 1 Common Configure:
```
"ACCOUNT_t_user_select_by_uids" :
{
    "sql" : "SELECT user_id, user_name, type FROM t_user WHERE user_id in ($condition$)",
    "noquote": {"condition":""},
    "noescape":{"condition":""},
    "check": {"condition": "^.*$"},
    "db" : "db_t_gpsbox_w"
}
```
* **noquote:{"table_name":"", "values":""}** 
```
not quoted when replacing the parameters.  
替换table_name和values这两个参数的时候,不在两边加引号.   
```

### 1.2 **eg:**   
* sql configure:   
```
"example1" :
{
    "sql": "select * from $table_name$ where uid = $uid$"
     noquote:{"table_name":""}
}
```

* client pass value:  
``` 
Ident:  "example1",
Params: map[string]string{
    "table_name": "t_user",
    "uid": "abc",
}
``` 

* real sql:  
```
select * from t_user where uid = "abc";  
```

* **noescape: {"condition":""}**  
```
not escape string when replacing the parameters.
决定加不加转义字符.需要往数据库写入引号的时候,要加上.
```
* **check: {"condition": "^.*$"}**  
```
使用正则表达式校验client传过来的参数是否符合要求.
Use regular expressions to check whether the client parameters match the rule.
eg: "check":   {"id": "^\\d+$"}
The id parameter must be number string.
```

* **"db": "db_t_gpsbox_w"**  
```
db_t_gpsbox_w对应的数据库配置在DB配置文件中.
Connect to this database: db_t_gpsbox_w
db_t_gpsbox_w is definded in DBConfigure.
```

## 2. Sharding Configure:  
```
"ACCOUNT_t_user_sharding":
{   
    "sql" : "SELECT user_id, user_name, type FROM t_user;",
    "sharding": {"dbseq": ""},
    "db" : "db_t_account_$dbseq$"
} 
```
* **sharding: {"dbseq": ""}**  
```
使用client传过来的dbseq值,替换dbname里面的“$dbseq$”。
Replace the value 'dbseq' in dbname. 
```

### 2.2 **eg:**   
* Mysql Sharding Example   
[MYSQL数据库分片实现](dbproxy_sharding.md)。  


## 3. Trancation Configure  
```
"ACCOUNT_t_user_insert_transaction" : 
{   
"sqlgroup": 
    [   
        {   
            "sql" : "INSERT INTO t_user(user_id, user_name, type) VALUES($id$, $name$ ,$type$);",
            "noquote": {"id":""},
            "check":   {"id": "^\\d+$"},
            "db" : "db_t_account_w"
        },  
        {   
            "sql" : "INSERT INTO t_images(id, name, image) VALUES($id$, $name$ ,$image$);",
            "noquote": {"id":""},
            "check":   {"id": "^\\d+$"},
            "db" : "db_t_image_w"
        }   
    ]   
}
```

* **Put the sql configure in "sqlgroup". It will execute with transaction.**   

### 3.2 ** eg: **   
* Mysql Multi DB Transcation  
[MYSQL多数据库事务实现](dbproxy_multi_db_transaction.md)。 


