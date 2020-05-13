/*
 * Create By Xinwenjia 2018-04-15
 * Modify From-https://github.com/toniz/gudp
 */

package dbproxy

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    "regexp"
    "strconv"
    "strings"
    "time"

    "github.com/golang/glog"

    . "github.com/toniz/gosuit/loader"
    _ "github.com/toniz/gosuit/loader/fileloader"
)

type ProxyDB struct {
    Ident      string       `json:"ident"`
    Balancer   string       `json:"balancer"`
    Driver     string       `json:"driver"`
    DB         string       `json:"db"`
    User       string       `json:"user"`
    Password   string       `json:"password"`
    Variables  string       `json:"variables"`
    Endpoint   string       `json:"endpoint"`
    Encoding   string       `json:"encoding"`
    Connection DBConnection `json:"connection"`
}

type DBConnection struct {
    MaxIdleTime int `json:"idletime"`
    Timeout     int `json:"timeout"`
    MaxCount    int `json:"maxcount"`
    MaxLifetime int `json:"lifetime"`
}

type Statement struct {
    SQL      string            `json:"sql"`
    DB       string            `json:"db"`
    NoQuote  map[string]string `json:"noquote"`
    NoEscape map[string]string `json:"noescape"`
    Check    map[string]string `json:"check"`
    Sharding map[string]string `json:"sharding"`
}

type ProxySQL struct {
    Ident string `json:"ident"`
    Statement
    SQLGroup []Statement `json:"sqlgroup"`
}

type DBCheck struct {
    Field string
    Regex string
}

type DBProxy struct {
    dbh map[string]*sql.DB
    sc  map[string]ProxySQL
}

type RowData map[string]string

// NewDBProxy creates a DBProxy.
// Call loadConfigure To Load Configure Data
func NewDBProxy() *DBProxy {
    s := &DBProxy{dbh: make(map[string]*sql.DB), sc: make(map[string]ProxySQL)}
    return s
}

// Read Proxy DB Configure From File.
// Initialize the DB Connection.
// DBHandle Will Being Cover When DB Ident Is The Same.
func (s *DBProxy) AddDBHandleFromFile(p string, ext string, prefix string) error {
    l, err := NewLoader("file")
    fmt.Println(err)
    fileList, err := l.GetList(p, ext, prefix)
    if err != nil {
        return err
    }

    for _, file := range fileList {
        var result []ProxyDB
        err := l.Load(file, &result)
        if err != nil {
            glog.Warningln("Load Json File Failed: ", err)
            continue
        }

        for _, c := range result {
            if len(c.Ident) > 0 {
                // Construct Connect String
                var connStr string
                switch c.Driver {
                case "postgres":
                    {
                        connStr = "postgres://" + c.User + ":" + c.Password + "@" + c.Endpoint + "/" + c.DB + "?sslmode=disable"
                    }
                default:
                    {
                        c.Driver = "mysql"
                        connStr = c.User + `:` + c.Password + `@tcp(` + c.Endpoint + `)/` + c.DB + `?charset=` + c.Encoding + c.Variables
                    }
                }
                glog.Infoln("Loading Config Connect To: ", connStr)

                // Open DB Connection
                db, err := sql.Open(c.Driver, connStr)
                if err != nil {
                    glog.Warningf("DB[%s] Connect Failed [%s]: %v", c.Ident, connStr, err)
                    continue
                }

                // DB Ping
                if err = db.Ping(); err != nil {
                    glog.Warningf("DB[%s] Ping Failed [%s]: %v", c.Ident, connStr, err)
                    continue
                }

                // Check Is Exists.
                if _, ok := s.dbh[c.Ident]; ok {
                    glog.Warningln("ProxyDB Handle Has Being Conver: ", c)
                }
                s.dbh[c.Ident] = db

                // Set ConnMaxCnt
                if c.Connection.MaxCount != 0 {
                    s.dbh[c.Ident].SetMaxOpenConns(c.Connection.MaxCount)
                }

                // Set ConnMaxLifetime
                if c.Connection.MaxLifetime != 0 {
                    duration := strconv.Itoa(c.Connection.MaxLifetime)
                    if d, e := time.ParseDuration(duration + "s"); e == nil {
                        s.dbh[c.Ident].SetConnMaxLifetime(d)
                    }
                }
            } else {
                glog.Warningln("Load Json File Not Found Ident: ", c)
            }
        }
    }

    // fmt.Println(s.dbh)
    return nil
}

// Add Proxy SQL Configure From File.
// SQL Configure Will Being Cover. When SQL Ident Is The Same.
func (s *DBProxy) AddProxySQLFromFile(p string, ext string, prefix string) error {
    l, err := NewLoader("file")
    fileList, err := l.GetList(p, ext, prefix)
    if err != nil {
        return err
    }

    for _, file := range fileList {
        var result []ProxySQL
        err := l.Load(file, &result)
        if err != nil {
            glog.Warningln("Load Json File Failed: ", err)
            continue
        }

        for _, c := range result {
            if len(c.Ident) > 0 {
                if _, ok := s.sc[c.Ident]; ok {
                    glog.Warningln("ProxySQL Configure Has Being Conver: ", c)
                }
                s.sc[c.Ident] = c
            } else {
                glog.Warningln("Load Json File Not Found Ident: ", c)
            }
        }
    }

    // fmt.Println(s.sc)
    return nil
}

// Get DB Handle By dbname
func (s *DBProxy) GetDBHandle(dbname string) (dbh *sql.DB, err error) {
    dbh = s.dbh[dbname]
    if dbh == nil {
        err = errors.New("Not Found DB Handle")
    }
    return
}

// EscapeString. Golang don`t have mysql_real_escape_string function.
// Using this funtion to make escape string.
func EscapeString(value string) string {
    value = strings.Replace(value, `\`, `\\`, -1)
    value = strings.Replace(value, `"`, `\"`, -1)
    return value
}

// Transform *sql.Rows to []map[string]string.
func TransformRowData(rows *sql.Rows) ([]RowData, error) {
    // Get column names
    columns, err := rows.Columns()
    if err != nil {
        return nil, err
    }

    // Make a slice for the values
    values := make([]string, len(columns))

    scanArgs := make([]interface{}, len(values))
    for i := range values {
        scanArgs[i] = &values[i]
    }

    // Fetch rows
    var res []RowData
    for rows.Next() {
        err = rows.Scan(scanArgs...)
        if err != nil {
            return nil, err
        }

        // Now do something with the data.
        // Here we just print each column as a string.
        data := make(RowData)
        for i, col := range values {
            data[columns[i]] = col
        }

        if rows.Err() == nil {
            res = append(res, data)
        } else {
            err = rows.Err()
        }
    }

    return res, err
}

func (s *DBProxy) ReplaceParameters(st Statement, params map[string]string) (string, string, error) {
    sqlc := st.SQL
    if len(sqlc) < 1 {
        return "", "", errors.New("Error: Not Found This Sql ID On AutoCommit Mode.")
    }

    dbname := st.DB
    for k, v := range params {
        // Check the parameter using regex
        // eg: "check":   {"id": "^\\d+$"}
        // The id parameter must be number string.
        if rex, ok := st.Check[k]; ok {
            var validParam = regexp.MustCompile(rex)
            if match := validParam.MatchString(v); !match {
                return "", "", errors.New(fmt.Sprintf("Param Check Failed regex[%s] param[%s]", rex, v))
            }
        }

        // Escape query string
        // eg: "noescape":{"values":""}
        // The `values`parameter don`t need Escape
        var val string
        if _, ok := st.NoEscape[k]; ok {
            val = v
        } else {
            val = EscapeString(v)
        }

        // Replace parameter
        if _, ok := st.NoQuote[k]; ok {
            sqlc = strings.Replace(sqlc, "$"+k+"$", val, -1)
        } else {
            sqlc = strings.Replace(sqlc, "$"+k+"$", "'"+val+"'", -1)
        }

        // DB Sharding Support
        if _, ok := st.Sharding[k]; ok {
            dbname = strings.Replace(dbname, "$"+k+"$", val, -1)
        }
    }

    return sqlc, dbname, nil
}

func (s *DBProxy) AutoCommit(ctx context.Context, ident string, params map[string]string) ([]RowData, error) {
    if _, ok := s.sc[ident]; !ok {
        err := errors.New(fmt.Sprintf("Error: [%s]Not Found In Configure..", ident))
        return nil, err
    }

    var res []RowData
    sqlc, dbname, err := s.ReplaceParameters(s.sc[ident].Statement, params)
    if err != nil {
        return nil, err
    }

    glog.Infof("Ident[%s] Sql[%s] Dbname[%s]", ident, sqlc, dbname)
    if dbh := s.dbh[dbname]; dbh != nil {
        rows, err := dbh.Query(sqlc)
        if err != nil {
            return nil, err
        }
        defer rows.Close()

        res, err = TransformRowData(rows)
        if err != nil {
            glog.Warningf(" %v", err)
        }
    } else {
        err = errors.New(fmt.Sprintf("Error: [%s]Not Found In Configure..", dbname))
    }
    return res, err
}

func (s *DBProxy) TransCommit(ctx context.Context, ident string, gparams []map[string]string) ([][]RowData, error) {
    if _, ok := s.sc[ident]; !ok {
        err := errors.New(fmt.Sprintf("Error: [%s]Not Found In Configure..", ident))
        return nil, err
    }
    var txRes [][]RowData
    var err error

    var dbs []string
    var sqls []string
    for i, st := range s.sc[ident].SQLGroup {
        sqlc, dbname, err := s.ReplaceParameters(st, gparams[i])
        if err != nil {
            return nil, err
        }

        dbs = append(dbs, dbname)
        sqls = append(sqls, sqlc)
    }

    // Get DB Handle, Then Set AutoCommit = false
    rollback := false
    txs := make(map[string]*sql.Tx)
    for _, dbname := range dbs {
        // Unique
        if _, ok := txs[dbname]; !ok {
            dbh := s.dbh[dbname]
            if dbh == nil {
                err = errors.New(fmt.Sprintf("Error: [" + ident + "]No Such DB Handle." + dbname))
                rollback = true
                break
            }

            tx, err := dbh.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault})
            if err != nil {
                err = errors.New(fmt.Sprintf("Warn: Set AutoCommit=0 Failed: " + dbname))
                rollback = true
                break
            }
            txs[dbname] = tx
        }
    }

    // Exec db query
    if !rollback {
        for i, q := range sqls {
            glog.Infof("Ident[%s] Exec [%s] Sql[%s] ", ident, dbs[i], q)

            rows, err := txs[dbs[i]].Query(q)
            if err != nil {
                return nil, err
            }
            defer rows.Close()

            res, err := TransformRowData(rows)
            if err != nil {
                glog.Warningf("DB Query Failed: seq[%d] sql[%s] db[%s]: %v", i, q, dbs[i], err)
                rollback = true
                break
            }
            txRes = append(txRes, res)
        }
    }

    for _, tx := range txs {
        if rollback {
            if err = tx.Rollback(); err != nil {
                glog.Warningf("Warn: Rollback Transcation Failed: %s [%v]", ident, err)
            }
        } else {
            if err = tx.Commit(); err != nil {
                glog.Warningf("Warn: Commit Transcation Failed: %s [%v]", ident, err)
            }
        }
    }

    return txRes, err
}

func (s *DBProxy) MultiInsert(ctx context.Context, ident string, gparams []RowData) ([]RowData, error) {
    var res []RowData
    var err error

    return res, err
}
