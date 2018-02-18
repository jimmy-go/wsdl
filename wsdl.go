// Package wsdl contains WSDL client.
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
	if err := xml.NewDecoder(resp.Body).Decode(dst); err != nil {
		return err
	}
	if err := resp.Body.Close(); err != nil {
		return err
	}
	return nil

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
	err = xml.NewDecoder(resp.Body).Decode(dst)
	if err != nil {
		return err
	}
	if err := resp.Body.Close(); err != nil {
		return err
	}
	return nil
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
