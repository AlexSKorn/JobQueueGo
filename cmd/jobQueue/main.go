package main

import (
	"fmt"
	"net/http"

	"github.com/AlexSKorn/goJobQueue/internal/routes"
)

func main() {
	jq := &routes.JobQueue{}
	router := routes.NewRouter(jq)

	port := 8080
	addr := fmt.Sprintf(":%d", port)

	fmt.Printf("Server listening on http://localhost:%s\n", addr)
	err := http.ListenAndServe(addr, router)

	if err != nil {
		panic(err)
	}
}
