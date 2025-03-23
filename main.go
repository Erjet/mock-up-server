package main

import (
	"fmt"
	"net/http"
)

var endPointsConfig EndPointsList

func main() {

	port := serverInit()

	http.HandleFunc("/", RequestReceiverHandler)
	http.HandleFunc("/SendRequest", RequestSenderHandler)

	fmt.Println("Server listening on :8080...")
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
