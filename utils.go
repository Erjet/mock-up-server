package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func readAndParseRequest(r *http.Request) (RequestData, error) {
	var requestData RequestData

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return requestData, fmt.Errorf("error reading body: %v", err)
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &requestData)
	if err != nil {
		return requestData, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return requestData, nil
}

func sendRESTRequest(requestData RequestData) (ResponseOutputStruct, error) {
	var response ResponseOutputStruct

	client := &http.Client{}
	var req *http.Request
	var err error

	switch requestData.Method {
	case "GET":
		req, err = http.NewRequest("GET", requestData.Url, nil)
	case "POST":
		req, err = http.NewRequest("POST", requestData.Url, bytes.NewBuffer([]byte(requestData.Body)))
	case "PUT":
		req, err = http.NewRequest("PUT", requestData.Url, bytes.NewBuffer([]byte(requestData.Body)))
	default:
		return response, fmt.Errorf("unsupported HTTP method: %s", requestData.Method)
	}

	if err != nil {
		return response, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return response, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, fmt.Errorf("error reading response: %v", err)
	}

	headersJSON, err := json.Marshal(resp.Header)
	if err != nil {
		return response, fmt.Errorf("error marshaling headers to JSON: %v", err)
	}

	response = ResponseOutputStruct{
		OutputCode: resp.StatusCode,
		OutputHead: string(headersJSON),
		OutputBody: string(bodyResp),
	}

	return response, nil
}

func sendSOAPRequest(requestData RequestData) (ResponseOutputStruct, error) {
	var response ResponseOutputStruct

	client := &http.Client{}
	req, err := http.NewRequest("POST", requestData.Url, bytes.NewBuffer([]byte(requestData.Body)))
	if err != nil {
		return response, fmt.Errorf("error creating SOAP request: %v", err)
	}

	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", "")

	resp, err := client.Do(req)
	if err != nil {
		return response, fmt.Errorf("error sending SOAP request: %v", err)
	}
	defer resp.Body.Close()

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, fmt.Errorf("error reading SOAP response: %v", err)
	}

	headersJSON, err := json.Marshal(resp.Header)
	if err != nil {
		return response, fmt.Errorf("error marshaling headers to JSON: %v", err)
	}

	response = ResponseOutputStruct{
		OutputCode: resp.StatusCode,
		OutputHead: string(headersJSON),
		OutputBody: string(bodyResp),
	}

	return response, nil
}

func sendJSONResponse(w http.ResponseWriter, response ResponseOutputStruct) {
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("error marshaling JSON: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
