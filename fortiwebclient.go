package fortiwebclient

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// FortiWebClient keeps connection data to FortiWeb
type FortiWebClient struct {
	URL      string
	Username string
	Password string
}

func encodeBase64(username string, password string) string {
	stringToEncode := strings.Join([]string{username, ":", password}, "")
	encoded := base64.StdEncoding.EncodeToString([]byte(stringToEncode))
	return encoded

}

func init() {

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}

// GetStatus returns status of FortiWeb device
func (f *FortiWebClient) GetStatus() string {

	client := &http.Client{}

	req, err := http.NewRequest("GET", strings.Join([]string{f.URL, "api/v1.0/System/Status/Status"}, ""), nil)
	req.Header.Add("Authorization", encodeBase64(f.Username, f.Password))
	response, error := client.Do(req)

	if error != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		os.Exit(-1)
	}

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	return string(body[:])

}

// CreateVirtualServer creates a virtual server object in FortiWeb
// Simplifies POST operation to external user
func (f *FortiWebClient) CreateVirtualServer(jsonBody string) error {

	response, error := f.doPost("api/v1.0/ServerObjects/Server/VirtualServer", jsonBody)

	if error != nil || response.StatusCode != 200 {
		fmt.Printf("The HTTP request failed with error %s, %d, %s\n", error, response.StatusCode, response.Status)
		return error
	}

	return nil
}

func (f *FortiWebClient) doPost(path string, jsonBody string) (*http.Response, error) {

	client := &http.Client{}

	req, error := http.NewRequest("POST",
		strings.Join([]string{f.URL, path}, ""),
		strings.NewReader(jsonBody))
	if error != nil {
		fmt.Printf("The HTTP request failed with error %s\n", error)
		return &http.Response{}, error
	}
	req.Header.Add("Authorization", encodeBase64(f.Username, f.Password))
	return client.Do(req)

}
