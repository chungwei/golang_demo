package main

type MysqlConf struct {
	Dsn             string `toml:"dsn"`
	Retry           int    `toml:"retry"`
	PoolMaxIdleConn int    `toml:"pool_max_idle_conn"`
	PoolMaxOpenConn int    `toml:"pool_max_open_conn"`
}

func loadMysqlConf() *MysqlConf {
	conf := new(MysqlConf)
	conf.Dsn = "bulleap_rw:xxxxxxxx@tcp(203.1.1.1:3306)/dbname?charset=utf8&timeout=2000ms"
	conf.Retry = 1
	conf.PoolMaxIdleConn = 3
	conf.PoolMaxOpenConn = 10

	return conf
}
