package main

type RedisConf struct {
	Addr     string `toml:"addr"`
	Password string `toml:"password"`
}

func loadRedisConf() *RedisConf {
	conf := new(RedisConf)
	conf.Password = "X63&%e#8rf$^93fOpe"
	conf.Addr = "8.142.157.45:6379"

	return conf
}
