package main

import (
	"net/http"
)

func main() {
	http.ListenAndServe(":8888", http.FileServer(http.Dir(".")))
}
