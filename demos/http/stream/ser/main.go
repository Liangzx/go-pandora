package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func streamHandler(w http.ResponseWriter, r *http.Request) {
	// 设置 Content-Type 和 Transfer-Encoding header
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Transfer-Encoding", "chunked")

	// 打开文件
	file, err := os.Open("largefile.bin")
	if err != nil {
		log.Printf("Failed to open file: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// 创建一个 bufio.Reader 来读取文件
	bufferSize := 8192 // 8KB 缓冲区大小
	buffer := make([]byte, bufferSize)

	// 循环读取文件内容并写入 http.ResponseWriter
	for {
		n, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading file: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			break
		}

		_, writeErr := w.Write(buffer[:n])
		if writeErr != nil {
			log.Printf("Error writing to client: %v", writeErr)
			return
		}
	}

	// 确保所有数据都被写入到客户端
	flusher, ok := w.(http.Flusher)
	if ok {
		flusher.Flush()
	}
}

func main() {
	http.HandleFunc("/stream", streamHandler)
	log.Println("Listening on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("ListenAndServe: %s", err)
	}
}
