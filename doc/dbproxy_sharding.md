# Mysql DB Sharding

## DB Configure File:  
```json
{
    "db_t_account_0" : 
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
    "db_t_account_1" : 
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
    "t_user_sharding":
    {
        "sql" : "SELECT user_id, user_name, type FROM t_user;",
	"sharding": {"dbseq": ""},
        "db" : "db_t_account_$dbseq$"
    }
}

```

## Then The GRPC Clien.go Call Like this:  

```go
s := NewDBProxy()
err := s.AddDBHandleFromFile("example/db", ".json", "db_*")
err := s.AddProxySQLFromFile("example/sql", "json", "sql_*")
uids := []int{100, 1003456, 2004000}
for _, uid := range uids {
    shardnum := uid / 1000000
    ident := "t_user_sharding",
    params := map[string]string{
        "dbseq": strconv.Itoa(shardnum),
    }

    res, err := mds.AutoCommit(context.TODO(), ident, params)
    if err != nil {
        log.Printf("Select? From DB Failed: %v", err)
    }
}
```

## BUILD
go build client.go 
Then client will got the data result:
```
Result Like: db_t_account_0(accountdb): SELECT user_id, user_name, type FROM t_user
Result Like: db_t_account_1(image):SELECT user_id, user_name, type FROM t_user
error: code = Unknown desc = Error: [db_t_account_2]Not Found In Configure..

```




