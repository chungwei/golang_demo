package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	InitRedisClient(loadRedisConf())
	defer GetRedisClient().Close()

	bizHandler()
}

func bizHandler() {
	ctx := context.Background()

	// set
	str, err := GetRedisClient().Set(ctx, "k1", `{"orderId":112233}`, 10*time.Second).Result()
	fmt.Println("set res string:", str, err)

	// get
	str, err = GetRedisClient().Get(ctx, "k1").Result()
	fmt.Println("get res string:", str, err)

	// del
	i, err := GetRedisClient().Del(ctx, "k1").Result()
	fmt.Println("del res int64:", i, err)

	// setnx
	b, err := GetRedisClient().SetNX(ctx, "k1", 1, 2*time.Minute).Result()
	fmt.Println("setnx res bool:", b, err)

	// setbit
	i, err = GetRedisClient().SetBit(ctx, "k2", 1000, 1).Result()
	fmt.Println("setbit res int64:", i, err)

	// getbit
	i, err = GetRedisClient().GetBit(ctx, "k2", 1000).Result()
	fmt.Println("getbit res int64:", i, err)

	// decrby
	i, err = GetRedisClient().DecrBy(ctx, "k3", 3).Result()
	fmt.Println("decrby res int64:", i, err)

	// incrby
	i, err = GetRedisClient().IncrBy(ctx, "k3", 2).Result()
	fmt.Println("incrby res int64:", i, err)

	// lpush
	i, err = GetRedisClient().LPush(ctx, "k4", `{"id": 1, "age": 23}`).Result()
	fmt.Println("lpush res int64:", i, err)

	// lrange
	strs, err := GetRedisClient().LRange(ctx, "k4", 0, 10).Result()
	fmt.Println("lrange res []string:", strs, err)

}
