package main

import (
	"context"
	"fmt"
	"time"

	cache "github.com/Code-Hex/go-generics-cache"
)

// Refs:
// 	Go 1.18 泛型的一些技巧和困扰 https://zhuanlan.zhihu.com/p/438252333
//	Go 1.18 泛型全面讲解：一篇讲清泛型的全部 https://juejin.cn/post/7080938405449695268

// 泛型函数
func Add[T int | float64](a, b T) T {
	return a + b
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// use simple cache algorithm without options.
	c := cache.NewContext[string, int](ctx)
	c.Set("a", 1)
	gota, aok := c.Get("a")
	gotb, bok := c.Get("b")
	fmt.Println(gota, aok) // 1 true
	fmt.Println(gotb, bok) // 0 false

	// Create a cache for Number constraint. key as string, value as int.
	nc := cache.NewNumber[string, int]()
	nc.Set("age", 26, cache.WithExpiration(time.Hour))

	incremented := nc.Increment("age", 1)
	fmt.Println(incremented) // 27

	decremented := nc.Decrement("age", 1)
	fmt.Println(decremented) // 26
}
