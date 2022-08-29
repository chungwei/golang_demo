package main

import (
	"fmt"
	"sync"
	"time"
)

var m1 map[int]int
var m2 map[int]int
var m3 sync.Map
var lock *sync.RWMutex

func main() {
	goodCaseRWMutex()
	fmt.Println("finish goodCaseRWMutex")

	goodCaseSyncMap()
	fmt.Println("finish goodCaseSyncMap")

	badCase()
	fmt.Println("finish badCase")
}

func badCase() {
	m1 = make(map[int]int)
	for i := 0; i < 20; i++ {
		go writeBadCase(i, i)
		go readBadCase(i)
	}
	time.Sleep(10 * time.Second)
}

func readBadCase(key int) {
	fmt.Println(m1[key])
}
func writeBadCase(key int, value int) {
	m1[key] = value
}

func goodCaseRWMutex() {
	m2 = make(map[int]int)
	lock = new(sync.RWMutex)
	for i := 0; i < 20; i++ {
		go writeGoodCase(i, i)
		go readGoodCase(i)
	}
	time.Sleep(3 * time.Second)
}

func goodCaseSyncMap() {
	for i := 0; i < 20; i++ {
		go m3.Store(i, i)
		go fmt.Println(m3.Load(i))
	}
	time.Sleep(3 * time.Second)
}

func readGoodCase(key int) {
	lock.RLock()
	fmt.Println(m2[key])
	lock.RUnlock()
}

func writeGoodCase(key int, value int) {
	lock.Lock()
	m2[key] = value
	lock.Unlock()
}
