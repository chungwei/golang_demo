package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	hc        *http.Client
	endpoints []string
)

// 定义全局复用
// 避免每个client定义一个导致 socket: too many open files
// 参考https://www.jianshu.com/p/50ed36e98459
var defaultTransport = http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   4 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	DialTLSContext:         nil,
	TLSClientConfig:        nil,
	TLSHandshakeTimeout:    0,
	DisableKeepAlives:      false,
	DisableCompression:     false,
	MaxIdleConns:           500,
	MaxIdleConnsPerHost:    100,
	MaxConnsPerHost:        100,
	IdleConnTimeout:        time.Duration(10) * time.Second,
	ResponseHeaderTimeout:  0,
	ExpectContinueTimeout:  1 * time.Second,
	TLSNextProto:           nil,
	ProxyConnectHeader:     nil,
	GetProxyConnectHeader:  nil,
	MaxResponseHeaderBytes: 0,
	WriteBufferSize:        0,
	ReadBufferSize:         0,
	ForceAttemptHTTP2:      true,
}

func getEndpoint() string {
	rnd := time.Now().UnixNano() / 1000
	idx := rnd % int64(len(endpoints))
	fmt.Println(rnd, idx)
	return endpoints[idx]
}

func InitHttpClient(cfg *HttpClientConf) {
	if cfg == nil {
		fmt.Println("InitHttpClient fail:invalid cfg")
		os.Exit(1)
	}

	hc = GetHttpClient(cfg.Timeout)
	endpoints = cfg.Addrs
	fmt.Printf("InitHttpClient done, conf=%+v\n", cfg)
}

func GetHttpClient(timeout int) *http.Client {
	return &http.Client{
		Transport:     &defaultTransport,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       time.Duration(timeout) * time.Second, // 整体超时时间
	}
}

func doHttpRequest(ctx context.Context, m string, url string, ps map[string]interface{}) (bs []byte, err error) {
	var request *http.Request
	if m == http.MethodGet {
		request, err = http.NewRequest(http.MethodGet, url, nil)
		if err == nil && ps != nil && len(ps) > 0 {
			q := request.URL.Query()
			for k, v := range ps {
				q.Add(k, fmt.Sprint(v))
			}
			request.URL.RawQuery = q.Encode()
		}
	} else if m == http.MethodPost {
		q := ""
		if ps != nil && len(ps) > 0 {
			for k, v := range ps {
				q = q + fmt.Sprintf("%v=%v&", k, v)
			}
			q = strings.TrimRight(q, "&")
		}
		request, err = http.NewRequest(http.MethodPost, url, strings.NewReader(q))
	}
	if err != nil {
		return nil, errors.New("doHttpRequest Marshal/NewRequest fail")
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := hc.Do(request)
	if err != nil || resp == nil || resp.Body == nil {
		err = fmt.Errorf("doHttpRequest fail,err=%v", err)
		return
	}
	defer resp.Body.Close()

	bs, err = ioutil.ReadAll(resp.Body)
	if err != nil || len(bs) == 0 || resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("doHttpRequest read resp.body fail: %v", err)
	}

	return bs, err
}
