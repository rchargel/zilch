package main

import (
	"os"
	"github.com/rchargel/zilch/zip"
)

func main() {
	server := zip.NewZipCodeController()

	port := os.Getenv("PORT")
	server.Start(port)
}
