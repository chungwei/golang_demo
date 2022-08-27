package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func main() {
	http.HandleFunc("/hello1", Hello1Handler)
	http.HandleFunc("/hello2", wrapper(Hello2Handler))

	err := http.ListenAndServe(":18080", nil)
	fmt.Println(err)
}

/**
统计API耗时 实现思路1（侵入业务逻辑代码）

测试：
curl http://127.0.0.1:18080/hello1?a=1
*/
func Hello1Handler(rw http.ResponseWriter, r *http.Request) {
	st := time.Now()
	defer func(r []byte) {
		if rcv := recover(); rcv != nil {
			rw.Write(r)
			fmt.Println("http_request_out, latency =", time.Since(st).Milliseconds())
		}
	}(nil)

	r.ParseForm()
	fmt.Println("http_request_in, uri =", r.URL.Path, " param =", r.Form)
	rw.Write([]byte("hello world"))
	time.Sleep(10 * time.Millisecond)
	fmt.Println("http_request_out, latency =", time.Since(st).Milliseconds())
}

/**
统计API耗时 实现思路2（使用中间件）

测试：
curl http://127.0.0.1:18080/hello2?a=1
*/
func Hello2Handler(ctx context.Context, params url.Values) []byte {
	return []byte("hello world")
}

type httpHandler func(ctx context.Context, params url.Values) []byte

func wrapper(handle httpHandler) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		st := time.Now()

		defer func(r []byte) {
			if rcv := recover(); rcv != nil {
				rw.Write(r)
				fmt.Println("http_request_out, latency =", time.Since(st).Milliseconds())
			}
		}(nil)

		r.ParseForm()
		fmt.Println("http_request_in, uri =", r.URL.Path, " param =", r.Form)
		resp := handle(context.Background(), r.Form)
		rw.Write(resp)
		time.Sleep(10 * time.Millisecond)
		fmt.Println("http_request_out, latency =", time.Since(st).Milliseconds())
	}
}
