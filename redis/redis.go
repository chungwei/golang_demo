package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

func InitRedisClient(cfg *RedisConf) {

	// 相关参数说明参看
	// https://blog.csdn.net/pengpengzhou/article/details/105385666
	opt := redis.Options{
		Network: "tcp", // 默认
		Addr:    cfg.Addr,
		Dialer:  nil,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			//fmt.Println(fmt.Sprintf("conn:%+v", cn))
			return nil
		},
		Username:           "",
		Password:           cfg.Password,
		DB:                 0,
		MaxRetries:         0, // 默认不重试
		MinRetryBackoff:    -1,
		MaxRetryBackoff:    -1,
		DialTimeout:        5 * time.Second, // 建立连接的超时时间
		ReadTimeout:        1 * time.Second, // 读超时
		WriteTimeout:       1 * time.Second, // 写超时
		PoolFIFO:           false,
		PoolSize:           100, // 连接池大小
		MinIdleConns:       5,   // 闲置的连接池
		MaxConnAge:         0,
		PoolTimeout:        2 * time.Second, //当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒
		IdleTimeout:        2 * time.Minute,
		IdleCheckFrequency: 60 * time.Second,
		TLSConfig:          nil,
		Limiter:            nil,
	}
	redisClient = redis.NewClient(&opt)
	res, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("InitRedisClient.Ping fail:" + err.Error())
		os.Exit(1)
	}
	fmt.Println(fmt.Sprintf("InitRedisClient.Ping done:%v", res))
}

func GetRedisClient() *redis.Client {
	if redisClient == nil {
		InitRedisClient(loadRedisConf())
	}
	return redisClient
}
