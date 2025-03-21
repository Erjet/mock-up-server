package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
)

func RequestSenderHandler(w http.ResponseWriter, r *http.Request) {
	requestData, err := readAndParseRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var Response ResponseOutputStruct

	if requestData.Protocol == "REST" {
		Response, err = sendRESTRequest(requestData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if requestData.Protocol == "SOAP" {
		Response, err = sendSOAPRequest(requestData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	sendJSONResponse(w, Response)
}

func RequestReceiverHandler(res http.ResponseWriter, req *http.Request) {
	requestType := req.URL.Query().Get("type")

	if requestType == "gui" {
		file, err := os.ReadFile("index.html")
		if err != nil {
			http.Error(res, fmt.Sprintf("Error reading HTML file: %v", err), http.StatusInternalServerError)
			return
		}
		res.Write(file)
	} else {
		contentType := req.Header.Get("Content-Type")

		switch contentType {
		case "application/json":
			handleRESTRequest(res, req)
		case "application/xml":
			handleSOAPRequest(res, req)
		default:
			http.Error(res, "Unsuported Content-Type", http.StatusNotImplemented)
			return
		}
	}
}

func handleRESTRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	fmt.Println("Request Body:", string(body))

	for i := 0; i < len(server_data.Data); i++ {

		responseJSON, err := json.Marshal(server_data.Data[i].OutputBody)
		if err != nil {
			http.Error(w, "Failed to generate JSON response", http.StatusInternalServerError)
			return
		}

		if server_data.Data[i].OutputHead == "" {
			w.Header().Set("Content-Type", "application/json")
		} else {
			w.Header().Set("Content-Type", server_data.Data[i].OutputHead)
		}

		w.WriteHeader(server_data.Data[i].OutputCode)
		w.Write(responseJSON)
		return
	}
}

func handleSOAPRequest(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < len(server_data.Data); i++ {
		if r.Method != http.MethodPost {
			continue
		}

		body, err := io.ReadAll(r.Body)
		fmt.Println(string(body))
		if err != nil {
			continue
		}
		defer r.Body.Close()

		respMessage := SOAPResponse{
			Message: fmt.Sprintf(server_data.Data[i].OutputBody),
		}

		responseXML, err := xml.MarshalIndent(SOAPEnvelope{
			Body: SOAPBody{
				Request: HelloRequest{
					Name: respMessage.Message,
				},
			},
		}, "", "  ")

		if err != nil {
			http.Error(w, "Failed to generate response", http.StatusInternalServerError)
			continue
		}

		if server_data.Data[i].OutputHead == "" {
			w.Header().Set("Content-Type", "text/xml; charset=utf-8")
		} else {
			w.Header().Set("Content-Type", server_data.Data[i].OutputHead)
		}

		w.WriteHeader(server_data.Data[i].OutputCode)
		w.Write(responseXML)
	}
}
