// Package wsdl contains WSDL client.
package wsdl

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type T struct {
	Purpose       string
	Input         Input
	ExpectedError error
}

type Input struct {
	Client   *http.Client
	Endpoint string
	Action   string
	Username string
	Password string
}

func TestSoap(t *testing.T) {
	tsBasic := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, tmplResponse)
	}))
	defer tsBasic.Close()

	table := []T{
		T{
			Input: Input{
				Client: &http.Client{
					Timeout: time.Second,
				},
				Endpoint: tsBasic.URL,
				Action:   "",
			},
			ExpectedError: nil,
		},
	}

	for i := range table {
		x := table[i]

		wclient, err := New(x.Input.Client)
		assert.Nil(t, err)
		if err != nil {
			t.Errorf("err : [%s]", err)
		}

		src := &SampleRequest{
			XSI: "Something",
		}

		var res *SampleResponse
		err = wclient.Soap(src, &res, x.Input.Endpoint, x.Input.Action)
		if err != x.ExpectedError {
			t.Errorf("expected [%s] actual [%s]", x.ExpectedError, err)
		}

		log.Printf("soap res [%#v]", res)
	}
}

// TestCustom demonstrates a custom request (basic auth, middlewares, etc.)
func TestCustom(t *testing.T) {
	tsBasic := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, p, _ := r.BasicAuth()
		if u != "admin" && p != "123456" {
			t.Logf("basic auth header not found")
			fmt.Fprintln(w, "unauthorized")
			return
		}

		fmt.Fprintln(w, tmplResponse)
	}))
	defer tsBasic.Close()

	table := []T{
		T{
			Purpose: "Custom method with basic auth",
			Input: Input{
				Client: &http.Client{
					Timeout: time.Second,
				},
				Endpoint: tsBasic.URL,
				Action:   "",
				Username: "admin",
				Password: "123456",
			},
			ExpectedError: nil,
		},
		T{
			Purpose: "Custom method with Failing auth",
			Input: Input{
				Client: &http.Client{
					Timeout: time.Second,
				},
				Endpoint: tsBasic.URL,
				Action:   "",
				Username: "",
				Password: "",
			},
			ExpectedError: io.EOF,
		},
	}

	for i := range table {
		x := table[i]

		wclient, err := New(x.Input.Client)
		if err != nil {
			t.Errorf("err : [%s]", err)
		}

		src := &SampleRequest{
			Envelope: "xx",
		}

		req, err := NewSoapRequest(src, x.Input.Endpoint, x.Input.Action)
		if err != nil {
			t.Errorf("err : [%s]", err)
		}
		req.SetBasicAuth(x.Input.Username, x.Input.Password)

		var res *SampleResponse
		err = wclient.Custom(req, &res)
		if err != x.ExpectedError {
			t.Errorf("expected [%s] actual [%s]", x.ExpectedError, err)
		}

		log.Printf("soap res [%#v]", res)
	}
}

func TestRaw(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, p, _ := r.BasicAuth()
		if u != "admin" && p != "123456" {
			t.Logf("basic auth header not found")
			fmt.Fprintln(w, "unauthorized")
			return
		}

		fmt.Fprintln(w, tmplResponse)
	}))
	defer ts.Close()

	wclient, err := New(http.DefaultClient)
	if err != nil {
		t.Errorf("err : [%s]", err)
	}

	b := []byte(tmplRequest)
	buf := bytes.NewBuffer(b)

	req, err := NewRawRequest(buf, ts.URL, "")
	if err != nil {
		t.Errorf("err : [%s]", err)
	}
	req.SetBasicAuth("admin", "123456")

	var res *SampleResponse
	err = wclient.Custom(req, &res)
	if err != nil {
		t.Errorf("err [%s]", err)
	}

	log.Printf("soap res [%#v]", res)

}

type SampleRequest struct {
	XMLName  xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Envelope string   `xml:"SOAP-ENV,attr"`
	XSI      string   `xml:"xmlns xsi,attr"`
	XSD      string   `xml:"xmlns xsd,attr"`
	Soap     string   `xml:"xmlns soap,attr"`
}

type SampleResponse struct {
	XMLName  xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Envelope string   `xml:"SOAP-ENV,attr"`
	XSI      string   `xml:"xmlns xsi,attr"`
	XSD      string   `xml:"xmlns xsd,attr"`
	Soap     string   `xml:"xmlns soap,attr"`
	Body     ResBody
}

type ResBody struct {
	XMLName                  string `xml:"http://namespaces.snowboard-info.com m"`
	EndorsingBoarderResponse string `xml:"GetEndorsingBoarderResponse,attr"`
	Endorsing                string `xml:"endorsingBoarder"`
}

const (
	tmplRequest = `<SOAP-ENV:Envelope
  xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/"
  SOAP-ENV:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
  <SOAP-ENV:Body>
    <m:GetEndorsingBoarder xmlns:m="http://namespaces.snowboard-info.com">
      <manufacturer>K2</manufacturer>
      <model>Fatbob</model>
    </m:GetEndorsingBoarder>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

	tmplResponse = `<SOAP-ENV:Envelope
  xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/"
  SOAP-ENV:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
  <SOAP-ENV:Body>
    <m:GetEndorsingBoarderResponse xmlns:m="http://namespaces.snowboard-info.com">
      <endorsingBoarder>Chris Englesmann</endorsingBoarder>
    </m:GetEndorsingBoarderResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`
)
