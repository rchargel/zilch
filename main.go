package main

import (
	"fmt"
	"github.com/rchargel/zilch/zilch"
	"os"
	"runtime"
)

func main() {
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)
	fmt.Printf("Running on %v CPU cores,\n", cpus)
	port := os.Getenv("PORT")
	zilch.StartServer("resources", port)
}
