/**
*ref:
*	https://juejin.cn/post/6997361079906992165
*	https://zhuanlan.zhihu.com/p/296409942
 */

package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os/exec"
	"sync"
)

func shell1() {
	// 第一种：只执行命令，不获取结果
	cmd := exec.Command("bash", "-c", "ls -l /var/log/")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func shell2() {
	cmd := exec.Command("bash", "-c", "ping www.baidu.com")
	// 第二种：执行命令，并获取结果
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("combined out:\n%s\n", string(out))
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", string(out))
}

func shell3() {
	// 第三种：执行命令，并区分stdout 和 stderr
}


func main() {
	//shell1()
	//shell2()

}

// 普通用法（一次性获取所有输出）
// 这里为了简化，我省去了stderr和其他信息

func Command1(cmd string) error {
	c := exec.Command("bash", "-c", cmd)
	// 此处是windows版本
	// c := exec.Command("cmd", "/C", cmd)
	output, err := c.CombinedOutput()
	fmt.Println(string(output))
	return err
}

// 实时显示

func Command2(cmd string) error {
	//c := exec.Command("cmd", "/C", cmd) 	// windows
	c := exec.Command("bash", "-c", cmd)  // mac or linux
	stdout, err := c.StdoutPipe()
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		reader := bufio.NewReader(stdout)
		for {
			readString, err := reader.ReadString('\n')
			if err != nil || err == io.EOF {
				return
			}
			fmt.Print(readString)
		}
	}()
	err = c.Start()
	wg.Wait()
	return err
}


// 输出和错误

func read(ctx context.Context, wg *sync.WaitGroup, std io.ReadCloser) {
	reader := bufio.NewReader(std)
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			readString, err := reader.ReadString('\n')
			if err != nil || err == io.EOF {
				return
			}
			fmt.Print(readString)
		}
	}
}

func Command3(ctx context.Context, cmd string) error {
	//c := exec.CommandContext(ctx, "cmd", "/C", cmd) // windows
	c := exec.CommandContext(ctx, "bash", "-c", cmd) // mac linux
	stdout, err := c.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := c.StderrPipe()
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	// 因为有2个任务, 一个需要读取stderr 另一个需要读取stdout
	wg.Add(2)
	go read(ctx, &wg, stderr)
	go read(ctx, &wg, stdout)
	// 这里一定要用start,而不是run 详情请看下面的图
	err = c.Start()
	// 等待任务结束
	wg.Wait()
	return err
}

