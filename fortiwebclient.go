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

// GetStatus returns status of FortiWeb device
func (f *FortiWebClient) GetStatus() string {

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &http.Client{}

	req, err := http.NewRequest("GET", f.URL, nil)

	req.Header.Add("Authorization", encodeBase64(f.Username, f.Password))

	response, error := client.Do(req)

	if error != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		os.Exit(-1)
	}

	fmt.Println(response)

	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	return string(body[:])

}
