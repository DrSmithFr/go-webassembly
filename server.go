package main

import (
	"fmt"
	"net/http"
)

const PORT = 8080;

func main() {
	fmt.Println(fmt.Sprintf("Server running on http://localhost:%d/", PORT))
	err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), http.FileServer(http.Dir("public")))

	if err != nil {
		panic(err)
	}
}
