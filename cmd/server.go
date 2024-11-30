package main

import (
	"fmt"
	"log"
	"net/http"

	// "sync"
	api "main/api"
)

func main() {
	router := api.NewRouter()

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server is running on port 8080")
}
