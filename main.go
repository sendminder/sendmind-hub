package main

import (
	"fmt"
	"net/http"
)

func main() {
	// "/" 경로에 대한 요청 처리
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World Go!")
	})

	// 서버 시작 (포트 8080)
	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
