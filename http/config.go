package main

var conf *HttpClientConf

type HttpClientConf struct {
	Addrs     []string `toml:"addrs"`
	Timeout   int      `toml:"timeout_sec"`
	AppKey    string   `toml:"app_key"`
	AppSecret string   `toml:"app_secret"`
}

func InitHttpClientConf() *HttpClientConf {
	if conf != nil {
		return conf
	}
	conf = &HttpClientConf{
		Addrs:     []string{"http://8.142.1.1:9072", "http://8.142.1.2:9072"},
		Timeout:   1,
		AppKey:    "app key",
		AppSecret: "app secret",
	}

	return conf
}
