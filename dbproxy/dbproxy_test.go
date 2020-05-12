/*
 * Create By Xinwenjia 2018-04-15
 * Modify From-https://github.com/toniz/gudp
 */

package dbproxy_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	. "github.com/toniz/gosuit/dbproxy"
)

var _ = Describe("Test Database Proxy", func() {
	s := NewDBProxy()

	// Load Mysql Configure
	Context("Test Add Proxy DB Handler", func() {
		err := s.AddDBHandleFromFile("example/db", ".json", "db_*")
		It("Should Return No Error", func() {
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("Test Add Proxy SQL Configure", func() {
		err := s.AddProxySQLFromFile("example/sql", "json", "sql_*")
		It("Should Return No Error", func() {
			Expect(err).NotTo(HaveOccurred())
		})
	})

	// Load PostgreSQL Configure
	Context("Test Add Proxy DB Handler", func() {
		err := s.AddDBHandleFromFile("example/db", ".json", "pg_*")
		It("Should Return No Error", func() {
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("Test Add Proxy SQL Configure", func() {
		err := s.AddProxySQLFromFile("example/sql", "json", "pgsql_*")
		It("Should Return No Error", func() {
			Expect(err).NotTo(HaveOccurred())
		})
	})

	// ...Load Any Other Drivers.
	// ....
	// ........................

	Describe("Test Mysql Proxy", func() {

		Context("Test Prepare: Create Table: ", func() {
			dbh, err := s.GetDBHandle("db_account_w")

			It("Should Return No Error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			dsql := "DROP TABLE IF EXISTS ibbwhat.t_user;"
			_, err = dbh.Query(dsql)
			It("Should Return No Error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			csql := `CREATE TABLE ibbwhat.t_user(
                Fuser_id int(11) NOT NULL COMMENT '用户ID',
                Frole_group varchar(255) NOT NULL DEFAULT '' COMMENT '角色',
                Fname varchar(255) NOT NULL DEFAULT '' COMMENT '名称',
                Fpassword varchar(255) NOT NULL DEFAULT 0 COMMENT '密码',
                Fpw_type int(11) NOT NULL DEFAULT 0 COMMENT '密码类型，加密方式',
                Fsalt varchar(255) NOT NULL DEFAULT 0 COMMENT '盐，分量多少创建时候生成',
                Fwx_login varchar(255) NOT NULL DEFAULT 0 COMMENT '绑定微信UNION ID',
                Fqq_login varchar(255) NOT NULL DEFAULT 0 COMMENT '绑定QQ',
                Ftel_login varchar(255) NOT NULL DEFAULT 0 COMMENT '绑定电话号码',
                Fstatus int(11) NOT NULL DEFAULT 1 COMMENT '状态: 1.使用  0.废弃',
                Fcreate_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                Fupdate_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                PRIMARY KEY (Fuser_id)
            )COMMENT='用户表';`
			_, err = dbh.Query(csql)
			It("Should Return No Error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("Test Auto Commit", func() {
			ident := "t_user_insert"
			params := map[string]string{"id": "1", "name": "test", "group": "1000"}
			_, err := s.AutoCommit(context.TODO(), ident, params)

			It("Should Return No Error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("Test Transcation Commit", func() {
			ident := "t_user_insert_transaction"
			gparams := []map[string]string{
				map[string]string{"id": "3", "name": "test3", "group": "1000"},
				map[string]string{"id": "4", "name": "test4", "group": "1001"},
			}

			_, err := s.TransCommit(context.TODO(), ident, gparams)
			It("Should Return No Error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("Test Insert Values In DB", func() {
			ident := "t_user_select_by_uids"
			params := map[string]string{"condition": "1,3,4"}
			res, err := s.AutoCommit(context.TODO(), ident, params)

			It("Should Return 3 Result", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(len(res)).To(Equal(3))
			})
		})

		Context("Test Finish: Drop Table: ", func() {
			dbh, err := s.GetDBHandle("db_account_w")
			It("Should Return No Error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			dsql := "DROP TABLE IF EXISTS ibbwhat.t_user;"
			_, err = dbh.Query(dsql)
			It("Should Return No Error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("Test PostgreSQL Proxy", func() {
		Context("Test Prepare: Create Table: ", func() {
			dbh, err := s.GetDBHandle("pg_account_w")

			It("Should Return No Error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			dsql := "DROP TABLE IF EXISTS t_user;"
			_, err = dbh.Query(dsql)
			It("Should Return No Error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			csql := `CREATE TABLE t_user(
                Fuser_id int NOT NULL ,
                Frole_group varchar NOT NULL DEFAULT '',
                Fname varchar NOT NULL DEFAULT '',
                Fpassword varchar NOT NULL DEFAULT 0,
                Fpw_type int NOT NULL DEFAULT 0,
                Fsalt varchar NOT NULL DEFAULT 0,
                Fwx_login varchar NOT NULL DEFAULT 0,
                Fqq_login varchar NOT NULL DEFAULT 0,
                Ftel_login varchar NOT NULL DEFAULT 0,
                Fstatus int NOT NULL DEFAULT 1,
                Fcreate_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                Fupdate_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                PRIMARY KEY (Fuser_id)
            );`
			_, err = dbh.Query(csql)
			It("Should Return No Error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("Test Auto Commit", func() {
			ident := "pg_user_insert"
			params := map[string]string{"id": "1", "name": "test", "group": "1000"}
			_, err := s.AutoCommit(context.TODO(), ident, params)

			It("Should Return No Error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("Test Transcation Commit", func() {
			ident := "pg_user_insert_transaction"
			gparams := []map[string]string{
				map[string]string{"id": "3", "name": "test3", "group": "1000"},
				map[string]string{"id": "4", "name": "test4", "group": "1001"},
			}

			_, err := s.TransCommit(context.TODO(), ident, gparams)
			It("Should Return No Error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("Test Insert Values In DB", func() {
			ident := "pg_user_select_by_uids"
			params := map[string]string{"condition": "1,3,4"}
			res, err := s.AutoCommit(context.TODO(), ident, params)

			It("Should Return 3 Result", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(len(res)).To(Equal(3))
			})
		})

		Context("Test Finish: Drop Table: ", func() {
			dbh, err := s.GetDBHandle("pg_account_w")
			It("Should Return No Error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			dsql := "DROP TABLE IF EXISTS ibbwhat.t_user;"
			_, err = dbh.Query(dsql)
			It("Should Return No Error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
