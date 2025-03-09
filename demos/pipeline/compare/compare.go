package compare

import (
	"log"
	"time"
)

func Sleep() {
	time.Sleep(time.Microsecond)
}

func quareNumber(num int) int {
	Sleep()
	result := num * num
	return result
}

func doubleNumbe(num int) int {
	Sleep()
	result := num * 2
	return result
}

// 顺序处理函数
func processNumbers(count int) []int {
	results := make([]int, count)
	for i := 0; i < count; i++ {
		num := i
		squared := quareNumber(num)
		squared = quareNumber(squared)
		squared = quareNumber(squared)
		squared = quareNumber(squared)
		doubled := doubleNumbe(squared)
		results[i] = doubled
	}
	return results
}

// 不使用piple
func NoPipe() {
	count := 1000

	// 顺序处理
	results := processNumbers(count)
	log.Println(len(results))

}

// 生成整数的函数
func generateNumbersCh(ch chan<- int, count int) {
	for i := 0; i < count; i++ {
		ch <- i
	}
	close(ch)
}

// 平方运算的函数
func squareNumbersCh(in <-chan int, out chan<- int) {
	for num := range in {
		Sleep()
		result := num * num
		out <- result
	}
	close(out)
}

// 结果乘以2的函数
func doubleNumbersCh(in <-chan int, out chan<- int) {
	for num := range in {
		Sleep()
		result := num * 2
		out <- result
	}
	close(out)
}

func WithPipe() {
	count := 1000

	// 定义通道
	numbersCh := make(chan int, 1000)
	squaredCh := make(chan int, 1000)
	squaredCh2 := make(chan int, 1000)
	squaredCh3 := make(chan int, 1000)
	squaredCh4 := make(chan int, 1000)
	doubledCh := make(chan int, 1000)

	// 启动生成整数的goroutine
	go generateNumbersCh(numbersCh, count)

	// 启动平方运算的goroutine
	go squareNumbersCh(numbersCh, squaredCh)

	go squareNumbersCh(squaredCh, squaredCh2)

	go squareNumbersCh(squaredCh2, squaredCh3)

	go squareNumbersCh(squaredCh3, squaredCh4)

	// 启动结果乘以2的goroutine
	go doubleNumbersCh(squaredCh4, doubledCh)
	// 收集并存储最终结果
	results := make([]int, count)
	for i := 0; i < count; i++ {
		results[i] = <-doubledCh
	}
}
