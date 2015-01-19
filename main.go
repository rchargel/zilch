package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/rchargel/zilch/zilch"
)

type consoleWriter struct{}

func (w consoleWriter) Write(p []byte) (n int, err error) {
	fmt.Print(string(p))
	return len(p), nil
}

func main() {
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)
	fmt.Printf("Running on %v CPU cores.\n", cpus)

	var country, file, appKey, outputFile string

	flag.StringVar(&country, "c", "", "The country to create CSV for")
	flag.StringVar(&file, "f", "", "Location of the Zip Codes file to parse")
	flag.StringVar(&appKey, "k", "", "The application key")
	flag.StringVar(&outputFile, "o", "", "The output file")
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

		var writer io.Writer
		writer = consoleWriter{}
		if len(outputFile) > 0 {
			file, err := os.Create(outputFile)
			defer file.Close()
			if err == nil {
				writer = file
			}
		}
		resolver := zilch.Resolver{CountryCode: strings.ToLower(country), AppKey: appKey, ZipCodes: zipCodes}
		resolver.OutputCSV(writer)
	}
}
