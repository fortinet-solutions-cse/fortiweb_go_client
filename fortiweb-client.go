package fortiwebclient

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type FortiWebClient struct {
	Url      string
	Username string
	Password string
}

func (f *FortiWebClient) getStatus() string {

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://192.168.122.40:90/api/v1.0/System/Status/Status", nil)

	req.Header.Add("Authorization", "YWRtaW46")
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
