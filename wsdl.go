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
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
)

var (
	// ErrClientNil error is returned when SetClient method
	// fails.
	ErrClientNil = errors.New("wsdl: http client is nil")
)

// WSDL client.
type WSDL struct {
	client *http.Client
}

// New returns a new WSDL client.
func New(client *http.Client) (*WSDL, error) {
	if client == nil {
		return nil, ErrClientNil
	}
	wl := &WSDL{
		client: client,
	}

	return wl, nil
}

// Custom method is keep for custom requests.
func (c *WSDL) Custom(r *http.Request, dst interface{}) error {
	resp, err := c.client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = xml.NewDecoder(resp.Body).Decode(dst)
	return err

}

// Soap make a SOAP call to endpoint.
func (c *WSDL) Soap(src, dst interface{}, url, action string) error {
	req, err := NewSoapRequest(src, url, action)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = xml.NewDecoder(resp.Body).Decode(dst)
	return err
}

// NewSoapRequest returns a prepared SOAP request.
func NewSoapRequest(src interface{}, url, action string) (*http.Request, error) {
	// FIXME; make WSDL buffer?
	buf := bytes.NewBuffer([]byte{})
	err := xml.NewEncoder(buf).Encode(src)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "text/xml;charset=UTF-8")
	req.Header.Set("SOAPAction", action)
	return req, nil
}

// NewRawRequest returns a prepared SOAP request.
func NewRawRequest(body io.Reader, url, action string) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "text/xml;charset=UTF-8")
	req.Header.Set("SOAPAction", action)
	return req, nil
}
