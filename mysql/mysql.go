package main

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var mysqlConn = new(sqlx.DB)

func InitMysql(cfg *MysqlConf) {
	//定义数据库对象

	//根据数据源dsn和mysql驱动, 创建数据库对象
	conn, err := sqlx.Open("mysql", cfg.Dsn)
	if err != nil {
		fmt.Printf(`init mysql fail: %+v`, err)
		os.Exit(1)
	}
	conn.SetMaxIdleConns(cfg.PoolMaxIdleConn)
	conn.SetMaxOpenConns(cfg.PoolMaxOpenConn)
	err = conn.Ping()
	if err != nil {
		fmt.Printf(`ping mysql fail: %+v||dsn=%s`, err, cfg.Dsn)
		os.Exit(1)
	}

	mysqlConn = conn
}

func GetMysqlConn() *sqlx.DB {
	return mysqlConn
}

func ReleaseMysqlConn() {
	if mysqlConn != nil {
		mysqlConn.Close()
	}
}

func ReleaseStmt(stmt *sqlx.NamedStmt) {
	if stmt != nil {
		stmt.Close()
	}
}
