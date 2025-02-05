package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {

}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("ArtiSpace is running on port 8080...")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Error occured when starting server: %v\n", err)
	}
}
