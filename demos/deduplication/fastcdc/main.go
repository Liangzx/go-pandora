package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/jotfs/fastcdc-go"
)

func main() {
	opts := fastcdc.Options{
		MinSize:     256 * 1024,
		AverageSize: 1 * 1024 * 1024,
		MaxSize:     4 * 1024 * 1024,
	}
	filename := `E:/test1.bin`
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	//data := make([]byte, 10*1024*1024)
	//rand.Read(data)
	chunker, _ := fastcdc.NewChunker(file, opts)

	for {
		chunk, err := chunker.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("data=%x fingerprint=%x  len=%d offset=%d\n", chunk.Data[:10], chunk.Fingerprint, chunk.Length, chunk.Offset)
	}
}