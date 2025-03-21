package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

var server_data TestData

func main() {
	jsonFile, err1 := os.Open("data.json")
	if err1 != nil {
		fmt.Println(err1)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	err2 := json.Unmarshal(byteValue, &server_data)
	if err2 != nil {
		fmt.Println(err2)
	}

	http.HandleFunc("/", RequestReceiverHandler)
	http.HandleFunc("/SendRequest", RequestSenderHandler)

	fmt.Println("SOAP server listening on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
