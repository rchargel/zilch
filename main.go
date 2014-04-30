package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/rchargel/zilch/zilch"
	"io"
	"os"
	"runtime"
	"strings"
)

type ConsoleWriter struct{}

func (w ConsoleWriter) Write(p []byte) (n int, err error) {
	fmt.Print(string(p))
	return len(p), nil
}

func main() {
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)
	fmt.Printf("Running on %v CPU cores.\n", cpus)

	var country, file string

	flag.StringVar(&country, "c", "", "The country to create CSV for")
	flag.StringVar(&file, "f", "", "Location of the Zip Codes file to parse")
	flag.Parse()

	if len(country) == 0 || len(file) == 0 {
		port := os.Getenv("PORT")
		zilch.StartServer("resources", port)
	} else {
		file, err := os.Open(file)
		defer file.Close()

		if err != nil {
			panic(err)
		}

		reader := bufio.NewReader(file)
		zipCodes := make([]string, 0, 100)
		for {
			if line, _, err := reader.ReadLine(); err == io.EOF {
				break
			} else if err != nil {
				panic(err)
			} else {
				zipCodes = append(zipCodes, strings.TrimSpace(string(line)))
			}
		}
		resolver := zilch.Resolver{strings.ToLower(country), zipCodes}
		resolver.OutputCSV(ConsoleWriter{})
	}
}
