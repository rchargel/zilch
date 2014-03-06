package main

import (
	"fmt"
	"net/http"
	"os"
)

func Root(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(response, "Welcome to Zilch!")
}

func main() {
	http.HandleFunc("/", Root)
	port := os.Getenv("PORT")
	fmt.Println("listening on port", port, "...")
	err := http.ListenAndServe(":" + port, nil)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
