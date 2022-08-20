package main

type RedisConf struct {
	Addr     string `toml:"addr"`
	Password string `toml:"password"`
}

func loadRedisConf() *RedisConf {
	conf := new(RedisConf)
	conf.Password = "abc123"
	conf.Addr = "127.0.0.1:6379"

	return conf
}
