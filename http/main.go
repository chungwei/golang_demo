package main

import (
	"context"
	"fmt"
)

func main() {
	InitHttpClient(loadHttpClientConf())

	bizHandler()
}

func bizHandler() {
	ctx := context.Background()

	// get: 无入参
	endpoint := getEndpoint() + "/ping"
	r, err := doHttpRequest(ctx, "GET", endpoint, nil)
	fmt.Println(endpoint, string(r), err)

	// get: 有入参
	endpoint = getEndpoint() + "/ping"
	r, err = doHttpRequest(ctx, "GET", endpoint, map[string]interface{}{"p1": 123})
	fmt.Println(endpoint, string(r), err)

	// post: 无入参
	endpoint = getEndpoint() + "/version"
	r, err = doHttpRequest(ctx, "POST", endpoint, nil)
	fmt.Println(endpoint, string(r), err)

	// post: 有入参
	endpoint = getEndpoint() + "/version"
	r, err = doHttpRequest(ctx, "POST", endpoint, map[string]interface{}{"p1": 123})
	fmt.Println(endpoint, string(r), err)

}
