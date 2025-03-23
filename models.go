package main

import (
	"encoding/xml"
)

type ResponseOutputStruct struct {
	OutputCode int    `json:"output_code"`
	OutputHead string `json:"output_head"`
	OutputBody string `json:"output_body"`
}

type input_data struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    SOAPBody `xml:"Body"`
}

type EndPointParams struct {
	Url        string                 `json:"domen"`
	Method     string                 `json:"method"`
	InputBody  map[string]interface{} `json:"input_body"`
	OutputCode int                    `json:"output_code"`
	OutputHead string                 `json:"output_head"`
	OutputBody string                 `json:"output_body"`
}

type RequestData struct {
	Protocol string `json:"protocol"`
	Body     string `json:"input_body"`
	Method   string `json:"method"`
	Url      string `json:"domen"`
}

type EndPointsList struct {
	List []EndPointParams `json:"data"`
}

type SOAPEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    SOAPBody `xml:"Body"`
}

type SOAPBody struct {
	XMLName xml.Name     `xml:"Body"`
	Request HelloRequest `xml:"HelloRequest,omitempty"`
}

type HelloRequest struct {
	XMLName xml.Name `xml:"HelloRequest"`
	Name    string   `xml:"Name"`
}

type SOAPResponse struct {
	XMLName xml.Name `xml:"HelloResponse"`
	Message string   `xml:"Message"`
}

type ServerSettings struct {
	Port int `json:"Port"`
}
