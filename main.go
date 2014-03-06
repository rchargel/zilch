package main

import (
	"fmt"
	"net/http"
	"os"
	"github.com/rchargel/zilch/zip"
)

func Root(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(response, "Welcome to Zilch!")
}

func main() {
	r := zip.ZipReader{"./resources/us_zip_code_database.csv"}
	ch := make(chan zip.ZipCodeEntry)

	go r.Read(ch)
	for i := range ch {
		i.WriteJson(os.Stdout)
	}

	http.HandleFunc("/", Root)
	port := os.Getenv("PORT")
	fmt.Println("listening on port", port, "...")
	err := http.ListenAndServe(":" + port, nil)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
