# Mysql Multi DB Transaction

## DB Configure File:  
```json
{
    "db_t_account_w" : 
    {
        "DBName" : "accountdb",
        "DBUser" : "account_rw",
        "DBPass" : "123456",
        "ConnString" : "127.0.0.1:4000",
        "ConnMaxIdleTime" : 60,
        "ConnTimeout" : 5,
        "ConnMaxCnt" : 100,
        "ConnMaxLifetime" : 3600,
        "ConnEncoding" : "utf8,utf8mb4"
    },
    "db_t_image_w" : 
    {
        "DBName" : "image",
        "DBUser" : "uds",
        "DBPass" : "uds1234",
        "ConnString" : "127.0.0.1:4000",
        "ConnMaxIdleTime" : 60,
        "ConnTimeout" : 5,
        "ConnMaxCnt" : 100,
        "ConnMaxLifetime" : 3600,
        "ConnEncoding" : "utf8,utf8mb4"
    }
}

```

## GUDP SQL Configure File:  
```json
{
    "t_user_insert_transaction" : 
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
}

```

## Then The GRPC Clien.go Call Like this:  

```go
    s := NewDBProxy()
    err := s.AddDBHandleFromFile("example/db", ".json", "db_*")
    err  = s.AddProxySQLFromFile("example/sql", "json", "sql_*")

    ident := "t_user_insert_transaction"
    gparams := []map[string]string{
        map[string]string{"id": "3", "name": "test3", "group": "1000"},
        map[string]string{"id": "4", "name": "test4", "group": "1001"},
    }
    _, err := s.TransCommit(context.TODO(), ident, gparams)
}
```

## BUILD
go build client.go 
Then client will got the data like follow sql :
```
Result Like: db_t_account_0(accountdb): INSERT INTO t_user(user_id, user_name, type) VALUES(11000, "uet" ,"23");
Result Like: db_t_image_w(image): INSERT INTO t_images(id, name, image) VALUES(10002, "uae" ,"111");
```



