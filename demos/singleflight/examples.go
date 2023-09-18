package main

import (
	"fmt"
	"golang.org/x/sync/singleflight"
	"sync"
	"sync/atomic"
	"time"
)

var count int32

func main() {
	total := 1000
	sg := &singleflight.Group{}

	var wg sync.WaitGroup
	wg.Add(total)
	key := "key"
	for i := 0; i < total; i++ {
		go func() {
			defer wg.Done()
			sg.Do(key, func() (interface{}, error) {
				res, err := getData(key)
				return res, err
			})
			// getData(key)
		}()
	}

	wg.Wait()
	fmt.Printf("total num is %v\n", count)
}

func getData(key string) (interface{}, error) {
	atomic.AddInt32(&count, 1)
	time.Sleep(10 * time.Millisecond)
	return "result", nil
}

