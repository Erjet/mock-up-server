package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
)

func RequestSenderHandler(w http.ResponseWriter, r *http.Request) {
	requestList, err := readAndParseRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var Response ResponseOutputStruct

	if requestList.Protocol == "REST" {
		Response, err = sendRESTRequest(requestList)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if requestList.Protocol == "SOAP" {
		Response, err = sendSOAPRequest(requestList)
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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var httpJSON map[string]interface{}
	err = json.Unmarshal(body, &httpJSON)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	for i := 0; i < len(endPointsConfig.List); i++ {

		if endPointsConfig.List[i].Url != r.URL.Path {
			continue
		}

		if endPointsConfig.List[i].Method != r.Method {
			continue
		}

		if !reflect.DeepEqual(endPointsConfig.List[i].InputBody, httpJSON) {
			continue
		}

		responseJSON, err := json.Marshal(endPointsConfig.List[i].OutputBody)
		if err != nil {
			http.Error(w, "Failed to generate JSON response", http.StatusInternalServerError)
			return
		}

		if endPointsConfig.List[i].OutputHead == "" {
			w.Header().Set("Content-Type", "application/json")
		} else {
			w.Header().Set("Content-Type", endPointsConfig.List[i].OutputHead)
		}

		w.WriteHeader(endPointsConfig.List[i].OutputCode)
		w.Write(responseJSON)
		return
	}

	http.Error(w, "Not Found", http.StatusNotFound)
}

func handleSOAPRequest(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < len(endPointsConfig.List); i++ {
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
			Message: fmt.Sprintf(endPointsConfig.List[i].OutputBody),
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

		if endPointsConfig.List[i].OutputHead == "" {
			w.Header().Set("Content-Type", "text/xml; charset=utf-8")
		} else {
			w.Header().Set("Content-Type", endPointsConfig.List[i].OutputHead)
		}

		w.WriteHeader(endPointsConfig.List[i].OutputCode)
		w.Write(responseXML)
	}
	http.Error(w, "Not Found", http.StatusNotFound)
}
