package main

import (
	"os"
	"github.com/rchargel/zilch/zip"
)

func main() {
	port := os.Getenv("PORT")
	zip.StartServer(port)
}
