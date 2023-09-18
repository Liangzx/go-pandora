package main

import (
	"crypto/sha256"
	"fmt"
	"github.com/restic/chunker"
	"io"
	"log"
	"os"
)

func main() {
	// create a chunker
	filename := `E:/test1.bin`
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	chk := chunker.New(file, chunker.Pol(0x3DA3358B4DC173))
	// reuse this buffer
	buf := make([]byte, 8*1024*1024)

	for i := 0; i < 5; i++ {
		chunk, err := chk.Next(buf)
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("offset=%d len=%d hash=%02x\n",chunk.Start, chunk.Length, sha256.Sum256(chunk.Data))
	}
}
