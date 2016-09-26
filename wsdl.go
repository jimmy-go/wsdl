// Package wsdl contains WSDL client.
//
// MIT License
//
// Copyright (c) 2016 Angel Del Castillo
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package wsdl

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

// WSDL struct defines _
type WSDL struct {
	client *http.Client
}

// New returns a new WSDL client.
func New(timeout time.Time) (*WSDL, error) {
	ws := &WSDL{
		client: &http.Client{
			Timeout: timeout,
		},
	}

	return w, nil
}

// Do makes a request.
func (c *WSDL) Do(body io.Reader, url, action string, dst io.Writer) error {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "text/xml; charset=utf-8")
	req.Header.Add("SOAPAction", action)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(dst, resp.Body)
	return err
}

// Soap make a SOAP call to endpoint.
func (c *WSDL) Soap(envelope *Envelope) error {

}

// Envelope struct for SOAP envelope.
type Envelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	XSI     string   `xml:"xmlns xsi,attr"`
	XSD     string   `xml:"xmlns xsd,attr"`
	Soap    string   `xml:"xmlns soap,attr"`
	Body    Body
}

// NewEnvelope returns a new Envelope type.
func NewEnvelope() *Envelope {
	se := &Envelope{
		XSI:  "http://www.w3.org/2001/XMLSchema-instance",
		XSD:  "http://www.w3.org/2001/XMLSchema",
		Soap: "http://schemas.xmlsoap.org/soap/envelope/",
	}
	return se
}

// Body struct for Envelope content.
type Body struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
	Content string   `xml:",innerxml"`
}

// Fault struct for SOAP fault.
type Fault struct {
	XMLName     xml.Name `xml:"Fault"`
	FaultCode   string   `xml:"faultcode"`
	FaultString string   `xml:"faultstring"`
	Detail      string   `xml:"detail"`
}
