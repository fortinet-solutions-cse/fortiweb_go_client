package fortiwebclient

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
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
func (f *FortiWebClient) CreateVirtualServer(
	name, ipv4Address, ipv6Address, interfaceName string,
	useInterfaceIP, enable bool) error {

	body := map[string]interface{}{
		"name":           name,
		"ipv4Address":    ipv4Address,
		"ipv6Address":    ipv6Address,
		"interface":      interfaceName,
		"useInterfaceIP": useInterfaceIP,
		"enable":         enable,
		"can_delete":     true,
	}

	jsonByte, err := json.Marshal(body)

	if err != nil {
		fmt.Printf("Error in json data: %s\n", err)
		return err
	}

	response, error := f.doPost("api/v1.0/ServerObjects/Server/VirtualServer", string(jsonByte))

	if error != nil || response.StatusCode != 200 {
		fmt.Printf("The HTTP request failed with error %s, %d, %s\n", error, response.StatusCode, response.Status)
		return error
	}

	return nil
}

// SingleOrMultiserverPool is used to define the pool as single server or balanced servers
type SingleOrMultiserverPool string
type ServerPoolType string

const (
	// SingleServer is used when there is single server in the pool
	SingleServer  SingleOrMultiserverPool = "Single Server"
	ServerBalance SingleOrMultiserverPool = "Server Balance"
)

const (
	ReverseProxy          ServerPoolType = "Reverse Proxy"
	OfflineProtection     ServerPoolType = "Offline Protection"
	TrueTransparentProxy  ServerPoolType = "True Transparent Proxy"
	TransparentInspection ServerPoolType = "TransparentInspection"
	WCCP                  ServerPoolType = "WCCP"
)

// CreateServerPool creates a virtual server pool object in FortiWeb
// Simplifies POST operation to external user
func (f *FortiWebClient) CreateServerPool(name string,
	singleOrMultiple SingleOrMultiserverPool,
	poolType ServerPoolType,
	comments string) error {

	body := map[string]interface{}{
		"name": name,
		"dissingleServerOrServerBalance": singleOrMultiple,
		"distype":                        poolType,
		"comments":                       comments,
		"can_delete":                     true,
	}

	jsonByte, err := json.Marshal(body)

	if err != nil {
		fmt.Printf("Error in json data: %s\n", err)
		return err
	}

	response, error := f.doPost("api/v1.0/ServerObjects/Server/ServerPool", string(jsonByte))

	if error != nil || response.StatusCode != 200 {
		fmt.Printf("The HTTP request failed with error %s, %d, %s\n", error, response.StatusCode, response.Status)
		return error
	}

	return nil
}

// CreateHTTPContentRoutingPolicy creates an HTTP Content Routing policy in FortiWeb
// Simplifies POST operation to external user
func (f *FortiWebClient) CreateHTTPContentRoutingPolicy(name, serverPool, matchSeq string) error {

	body := map[string]interface{}{
		"name":       name,
		"serverPool": serverPool,
		"matchSeq":   matchSeq,
		"can_delete": true,
	}

	jsonByte, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("Error in json data: %s\n", err)
		return err
	}

	response, error := f.doPost("api/v1.0/ServerObjects/Server/HTTPContentRoutingPolicy", string(jsonByte))

	if error != nil || response.StatusCode != 200 {
		fmt.Printf("The HTTP request failed with error %s, %d, %s\n", error, response.StatusCode, response.Status)
		return error
	}

	return nil
}

// CreateHTTPContentRouting creates a criteria for matching http content in a policy
// Simplifies POST operation to external user
// Use this json:
// {
// 	"id": "1",
// 	"_id": "1",
// 	"seqId": 1,
// 	"realId": "1",
// 	"matchObject": 2,
// 	"matchExpression": "huy",
// 	"urlCondition": 3,
// 	"concatenate": 2
// }
//}
func (f *FortiWebClient) CreateHTTPContentRouting(HTTPContentRoutingPolicy string, jsonBody string) error {

	url := strings.Join([]string{"api/v1.0/ServerObjects/Server/HTTPContentRoutingPolicy/",
		HTTPContentRoutingPolicy,
		"/HTTPContentRoutingPolicyNewHTTPContentRouting"},
		"")
	response, error := f.doPost(url, jsonBody)

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
