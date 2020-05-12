/*
 * Create By Xinwenjia 2018-04-15
 * Modify From-https://github.com/toniz/gudp
 */

package etcdloader_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/toniz/gosuit/loader"
	_ "github.com/toniz/gosuit/loader/etcdloader"
	"github.com/toniz/gosuit/parser"
	"go.etcd.io/etcd/clientv3"
)

var _ = Describe("Etcd Loader", func() {
	endpoint := "10.106.210.224:2379"
	user := "root"
	password := "e3jSlAsGNw"

	Context("Init Etcd Data", func() {
		testDBCnfName1 := "/mysql/db/db_account_w.json"
		testDBCnfValue1 := `{
            "DBName" : "accountdb",
            "DBUser" : "account_rw",
            "DBPass" : "123456",
            "ConnString" : "127.0.0.1:4000",
            "ConnMaxIdleTime" : 60,
            "ConnTimeout" : 5,
            "ConnMaxCnt" : 100,
            "ConnMaxLifetime" : 3600,
            "ConnEncoding" : "utf8,utf8mb4"
        }`
		testDBCnfName2 := "/mysql/db/db_account_r.json"
		testDBCnfValue2 := `{
            "DBName" : "accountdb",
            "DBUser" : "account_r",
            "DBPass" : "123456",
            "ConnString" : "127.0.0.1:4000",
            "ConnMaxIdleTime" : 60,
            "ConnTimeout" : 5,
            "ConnMaxCnt" : 100,
            "ConnMaxLifetime" : 3600,
            "ConnEncoding" : "utf8,utf8mb4"
        }`
		testSqlsCnfName := "/mysql/sqls/account/t_user_insert.json"
		testSqlsCnfValue := `{
            "sql" : "INSERT IGNORE INTO t_user(Fuser_id, Fname, Frole_group) VALUES($id$, $name$ ,$group$);",
            "noquote": {"id":""},
            "check":   {"id": "^\\d+$"},
            "db" : "db_account_w"
        }`

		It("Should Return Not Errar", func() {
			cli, err := clientv3.New(clientv3.Config{Endpoints: []string{endpoint}, Username: user, Password: password})
			cli.Put(context.TODO(), testDBCnfName1, testDBCnfValue1)
			cli.Put(context.TODO(), testDBCnfName2, testDBCnfValue2)
			cli.Put(context.TODO(), testSqlsCnfName, testSqlsCnfValue)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("Test Json Etcd List", func() {
		p := "/mysql/db/"
		l, err := loader.NewLoader("etcd")
		It("Should Return Loader Interface", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should Return Not Errar", func() {
			err = l.Connect(endpoint, user, password)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should Return Two Json File", func() {
			filelist, err := l.GetList(p, ".json", "db")
			Expect(len(filelist)).To(Equal(2))
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("Test Json Etcd Load", func() {
		l, err := loader.NewLoader("etcd")
		It("Should Return A Json File", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should Return Not Errar", func() {
			err = l.Connect(endpoint, user, password)
			Expect(err).NotTo(HaveOccurred())
		})

		type testJson struct {
			DBUser string
			DBPass string
		}
		p := "/mysql/db/db_account_w.json"
		var s testJson
		It("Should Return A Json File", func() {
			err = l.Load(p, &s)
			Expect(err).NotTo(HaveOccurred())
			Expect(s.DBUser).To(Equal("account_rw"))
			Expect(s.DBPass).To(Equal("123456"))
		})
	})

	Context("Test Json Etcd LoadAll", func() {
		l, err := loader.NewLoader("etcd")
		It("Should Return Two Json File", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should Return Not Errar", func() {
			err = l.Connect(endpoint, user, password)
			Expect(err).NotTo(HaveOccurred())
		})

		type testJson struct {
			DBUser string
			DBPass string
		}
		p := "/mysql/db/"
		s := make(map[string]testJson)
		It("Should Return Two Json File", func() {
			sc, err := l.ReadAll(p, "", "db_")
			for k, v := range sc {
				var l testJson
				err = parser.Decode(".json", v, &l)
				s[k] = l
			}
			Expect(err).NotTo(HaveOccurred())
			Expect(s["/mysql/db/db_account_w.json"].DBUser).To(Equal("account_rw"))
			Expect(s["/mysql/db/db_account_w.json"].DBPass).To(Equal("123456"))
			Expect(s["/mysql/db/db_account_r.json"].DBUser).To(Equal("account_r"))
		})
	})

	Context("Delete Etcd File", func() {
		It("Should Delete Three Json File", func() {
			cli, err := clientv3.New(clientv3.Config{Endpoints: []string{endpoint}, Username: user, Password: password})
			Expect(err).NotTo(HaveOccurred())

			testDBCnfName1 := "/mysql/db/db_account_w.json"
			cli.Delete(context.TODO(), testDBCnfName1)

			testDBCnfName2 := "/mysql/db/db_account_r.json"
			cli.Delete(context.TODO(), testDBCnfName2)

			testSqlsCnfName := "/mysql/sqls/account/t_user_insert.json"
			cli.Delete(context.TODO(), testSqlsCnfName)
		})
	})
})
