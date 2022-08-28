package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	var count = 100 // 并发数
	for i := 0; i < count; i++ {
		wg.Add(1) // 计数器+1
		go func(wg *sync.WaitGroup) {
			defer wg.Done() // 计数器-1
			doSth()
		}(&wg) // 必须传指针 否则会死锁
	}
	wg.Wait() // 等待计数器变为0
}

func doSth() {
	fmt.Println("do sth:", time.Now().UnixNano())
}
