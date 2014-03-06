package zip

import (
	"os"
	"io"
	"fmt"
	"encoding/csv"
)

type ZipReader struct {
	Path string
} 

func (r ZipReader) Read (ch chan []string) {
	fmt.Println("Reading file:",r.Path)
	file, err := os.Open(r.Path)
	if err != nil {
		fmt.Println("Error:",err)
		close(ch)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:",err)
			break
		}
		ch <- record
	}
	close(ch)
}
