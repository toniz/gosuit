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
    "sync"
    "reflect"

    //"github.com/golang/glog"
    . "github.com/toniz/gosuit/loader"
    _ "github.com/toniz/gosuit/loader/fileloader"

    "github.com/toniz/gosuit/glog"
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
    IsoLevel string `json:"isolevel"`
    Statement
    SQLGroup []Statement `json:"sqlgroup"`
}

type DBCheck struct {
    Field string
    Regex string
}

type DBProxy struct {
    rwlock *sync.RWMutex
    dbh map[string]*sql.DB
    sc  map[string]ProxySQL
    dbc  map[string]ProxyDB
}

type RowData map[string]string

// NewDBProxy creates a DBProxy.
// Call loadConfigure To Load Configure Data
func NewDBProxy() *DBProxy {
    rwl := &sync.RWMutex{}
    s := &DBProxy{
        dbh: make(map[string]*sql.DB),
        dbc: make(map[string]ProxyDB),
        sc: make(map[string]ProxySQL),
        rwlock: rwl}
    return s
}

// Close DBProxy Handle.
func (s *DBProxy) Close() {
    for _, dbh := range s.dbh {
        dbh.Close()
    }
}

// Read Proxy DB Configure From File.
// Initialize the DB Connection.
// DBHandle Will Being Cover When DB Ident Is The Same.
func (s *DBProxy) AddDBHandleFromFile(p string, ext string, prefix string) error {
    l, err := NewLoader("file")
    //fmt.Println(err)
    fileList, err := l.GetList(p, ext, prefix)
    if err != nil {
        return err
    }

    dbConfig := make(map[string]ProxyDB)
    for _, file := range fileList {
        var result []ProxyDB
        err := l.Load(file, &result)
        if err != nil {
            glog.Warningf("Load Json File[%v] Failed: %v", file, err)
            continue
        }

        for _, cnf := range result {
            if len(cnf.Ident) > 0 {
                dbConfig[cnf.Ident] = cnf
            } else {
                glog.Warningln("Load Json File Not Found Ident: ", cnf)
            }
        }
    }

    updateDBH := false
    s.rwlock.RLock()
    if !reflect.DeepEqual(s.dbc, dbConfig) {
        for k, v := range dbConfig {
            s.dbc[k] = v
        }
        updateDBH = true
    }
    s.rwlock.RUnlock()

    if updateDBH {
        dbHandle := make(map[string]*sql.DB)
        for k, c := range s.dbc {
            // Construct Connect String
            var connStr string
            switch c.Driver {
                case "postgres": {
                    connStr = "postgres://" + c.User + ":" + c.Password + "@" + c.Endpoint + "/" + c.DB + "?sslmode=disable"
                }
                default: {
                    c.Driver = "mysql"
                    connStr = c.User + `:` + c.Password + `@tcp(` + c.Endpoint + `)/` + c.DB + `?charset=` + c.Encoding + c.Variables
                }
            }
            glog.V(11).Infof("Loading Config Connect To [%s]: %v ", k, connStr)

            // Open DB Connection
            db, err := sql.Open(c.Driver, connStr)
            if err != nil {
                glog.Warningf("DB[%s] Connect Failed [%s]: %v", k, connStr, err)
                continue
            }

            // DB Ping
            if err = db.Ping(); err != nil {
                glog.Warningf("DB[%s] Ping Failed [%s]: %v", k, connStr, err)
                continue
            }

            dbHandle[k] = db
            // Set ConnMaxCnt
            if c.Connection.MaxCount != 0 {
                dbHandle[k].SetMaxOpenConns(c.Connection.MaxCount)
            }

            // Set ConnMaxLifetime
            if c.Connection.MaxLifetime != 0 {
                duration := strconv.Itoa(c.Connection.MaxLifetime)
                if d, e := time.ParseDuration(duration + "s"); e == nil {
                    dbHandle[k].SetConnMaxLifetime(d)
                }
            }
        }

        s.rwlock.Lock()
        defer s.rwlock.Unlock()
        glog.V(11).Infof("Load DB Configure From File[%d] dbh[%d]", len(dbHandle), len(s.dbh))
        if len(dbHandle) > 0 {
            for _, dbh := range s.dbh {
                dbh.Close()
            }
            s.dbh = dbHandle
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

    sqlConfig := make(map[string]ProxySQL)
    for _, file := range fileList {
        var result []ProxySQL
        err := l.Load(file, &result)
        if err != nil {
            glog.Warningf("Load Json File[%v] Failed: %v", file, err)
            continue
        }

        for _, cnf := range result {
            if len(cnf.Ident) > 0 {
                sqlConfig[cnf.Ident] = cnf
            } else {
                glog.Warningln("Load Json File Not Found Ident: ", cnf)
            }
        }
    }

    s.rwlock.Lock()
    defer s.rwlock.Unlock()
    if !reflect.DeepEqual(s.sc, sqlConfig) {
        for k, v := range sqlConfig {
            s.sc[k] = v
        }
        keys := reflect.ValueOf(s.sc).MapKeys()
        glog.V(10).Infof("SQL Configure Refresh From File[%d] sc[%d]: %v", len(sqlConfig), len(s.sc), keys)
    }
    // fmt.Println(s.sc)
    return nil
}

// Get DB Handle By dbname
func (s *DBProxy) GetDBHandle(dbname string) (dbh *sql.DB, err error) {
    s.rwlock.RLock()
    defer s.rwlock.RUnlock()

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
    values := make([]sql.NullString, len(columns))

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
            if col.Valid {
                data[columns[i]] = col.String
            } else {
                data[columns[i]] = ""
            }
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
    s.rwlock.RLock()
    defer s.rwlock.RUnlock()

    if s == nil {
        err := errors.New(fmt.Sprintf("Error: [%s] Configure Not Init.", ident))
        return nil, err
    }

    if _, ok := s.sc[ident]; !ok {
        err := errors.New(fmt.Sprintf("Error: [%s]Not Found In Configure..", ident))
        return nil, err
    }

    var res []RowData
    sqlc, dbname, err := s.ReplaceParameters(s.sc[ident].Statement, params)
    if err != nil {
        return nil, err
    }

    glog.V(11).Infof("Ident[%s] Sql[%s] Dbname[%s]", ident, sqlc, dbname)
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
    s.rwlock.RLock()
    defer s.rwlock.RUnlock()

    if s == nil {
        err := errors.New(fmt.Sprintf("Error: [%s] Configure Not Init.", ident))
        return nil, err
    }

    if _, ok := s.sc[ident]; !ok {
        err := errors.New(fmt.Sprintf("Error: [%s]Not Found In Configure..", ident))
        return nil, err
    }

    isoLevel := sql.LevelDefault
    if len(s.sc[ident].IsoLevel) != 0 {
        isoInt, _ := strconv.Atoi(s.sc[ident].IsoLevel)
        isoLevel = sql.IsolationLevel(isoInt)
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
    glog.V(12).Infof("Ident[%s] Exec DBs[%v] Sqls[%v] ", ident, dbs, sqls)

    // Get DB Handle, Then Set AutoCommit = false
    rollback := false
    txs := make(map[string]*sql.Tx)
    for _, dbname := range dbs {
        // Unique
        if _, ok := txs[dbname]; !ok {
            dbh := s.dbh[dbname]
            if dbh == nil {
                err = errors.New(fmt.Sprintf("Error: [" + ident + "]No Such DB Handle." + dbname))
                glog.Warningf("Ident[%s] Exec Failed: %v", ident, err)
                rollback = true
                break
            }

            tx, err := dbh.BeginTx(ctx, &sql.TxOptions{Isolation: isoLevel})
            if err != nil {
                err = errors.New(fmt.Sprintf("Warn: Set AutoCommit=0 Failed: " + dbname))
                glog.Warningf("Ident[%s] Exec Failed: %v", ident, err)
                rollback = true
                break
            }
            txs[dbname] = tx
        }
    }

    // Exec db query
    if !rollback {
        for i, q := range sqls {
            glog.V(11).Infof("Ident[%s] Exec [%s] Sql[%s] ", ident, dbs[i], q)

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

func (s *DBProxy) MultiInsert(ctx context.Context, ident string, values [][]string) error {
    var insRows string
    for i, row := range values {
        rowStr := strings.Join(row,"', '")
        rowStr = "('" + rowStr + "')"
        if i != 0  {
            insRows = insRows + ", "
        }
        insRows = insRows + rowStr
    }
    params := map[string]string{"VALUES": insRows}
    _, err := s.AutoCommit(ctx, ident, params)
    return err
}


