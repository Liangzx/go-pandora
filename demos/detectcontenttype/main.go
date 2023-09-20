package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	// 打开一个 MP3 文件
	filePathName := os.Args[1]
	file, err := os.Open(filePathName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// 读取文件的前 512 个字节
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 调用 http.DetectContentType 方法判断文件类型
	// 实际上，如果字节数超过 512，该函数也只会使用前 512 个字节
	contentType := http.DetectContentType(buffer[:n])
	fmt.Println(contentType) // 输出 audio/mpeg
}

// https://www.cnblogs.com/enjong/articles/10741244.html
