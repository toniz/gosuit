# Mysql Read Write Splitting

## DB Grant 
```
GRANT ALL PRIVILEGES ON accountdb.* TO 'account_rw'@'127.0.0.1' IDENTIFIED BY '123456' WITH GRANT OPTION;
GRANT SELECT  ON accountdb.* TO 'account_r'@'127.0.0.1' IDENTIFIED BY '123456' WITH GRANT OPTION;
```

## GUDP DB Configure File:  
```json
[
    {
        "ident": "db_account_w",
        "driver": "mysql",
        "db" : "ibbwhat",
        "user" : "ibbwhat",
        "password" : "123456",
        "endpoint" : "10.107.152.167:3306",
        "encoding" : "utf8,utf8mb4a",
        "connection" : {
            "maxcount" : 100,
            "lifetime" : 3600,
            "timeout" : 5
        }
    },
    {
        "ident" : "db_account_r",
        "driver": "mysql",
        "db" : "ibbwhat",
        "user" : "ibbwhat",
        "password" : "123456",
        "endpoint" : "10.107.152.167:3306",
        "encoding" : "utf8,utf8mb4",
        "connection" : {
            "maxcount" : 100,
            "lifetime" : 3600,
            "timeout" : 5
        }
    }
]
```

## GUDP SQL Configure File:  
```json
[
    {
        "ident":"t_user_insert",
        "sql" : "INSERT IGNORE INTO t_user(Fuser_id, Fname, Frole_group) VALUES($id$, $name$ ,$group$);",
        "noquote": {"id":""},
        "check":   {"id": "^\\d+$"},
        "db" : "db_account_w"
    },
    {
        "ident" :"t_user_update",
        "sql" : "UPDATE t_user SET Fname = $name$ WHERE Fuser_id = $id$;",
        "db" : "db_t_gpsbox_w"
    },
    {
        "ident":"t_user_delete" ,
        "sql" : "DELETE FROM t_user WHERE user_id = $id$",
        "db" : "db_t_gpsbox_w"
    },
    {
        "ident":"t_user_insert_multi",
        "sql" : "REPLACE INTO $table_name$(user_id, user_name, type) VALUES $values$",
        "noquote": {"table_name":"", "values":""},
        "noescape":{"values":""},
        "check":   {"values": "^.*$"},
        "db" : "db_t_gpsbox_w"
    },
    {
        "ident": "t_user_select_by_uid",
        "sql" : "SELECT user_id, user_name, type FROM t_user WHERE user_id>=$limit_start$ ORDER BY user_id ASC LIMIT $limit_end$ ;",
        "noquote" : {"limit_start":"", "limit_end":""},
        "db" : "db_account_r"
    },
    {
        "ident":"t_user_select_by_uids",
        "sql" : "SELECT Fuser_id, Fname, Frole_group FROM t_user WHERE Fuser_id in ($condition$)",
        "noquote": {"condition":""},
        "noescape":{"condition":""},
        "check": {"condition": "^.*$"},
        "db" : "db_account_w"
    }
]

```

## Then The GRPC Clien.go Call Like this:  

```go
    mds := NewMysqlDataService()
    mds.Cfg.DBPath = "./config/db/"
    mds.Cfg.SqlPath = "./config/sql/"
    err := mds.InitMysqlConnection()
    if err != nil {
        log.Fatalf("Connect To Mysql Failed:", err)
    }

    ident = "t_user_select_by_uid"
    params = map[string]string{
        "limit_start": "100",
        "limit_end": "10"
    }
    res, err := mds.AutoCommit(context.TODO(), ident, params)
    if err != nil {
        log.Printf("Select? From DB Failed: %v", err)
    }
```

## BUILD
go build client.go 
Then client will got the data like follow sql :
```
SELECT user_id, user_name, type FROM t_user WHERE user_id>=100 ORDER BY user_id ASC LIMIT 10 ;
```


