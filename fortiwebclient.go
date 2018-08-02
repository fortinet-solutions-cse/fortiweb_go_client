package fortiwebclient

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
		return strings.Join([]string{"Error: The HTTP request failed with error ", error.Error()}, "")
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

type singleServerOrServerBalance int

type serverPoolType int
type loadBalancingAlgorithm int

const (
	_ = iota
	// SingleServer is used when there is single server in the pool
	SingleServer singleServerOrServerBalance = iota
	// ServerBalance is used when there is a cluster of servers
	ServerBalance singleServerOrServerBalance = iota
)

const (
	_ = iota
	// ReverseProxy hides all servers behind FortiWeb
	ReverseProxy serverPoolType = iota
	// OfflineProtection Puts FortiWeb in sniffer mode
	OfflineProtection serverPoolType = iota
	// TrueTransparentProxy Puts FortiWeb as a transparent proxy
	TrueTransparentProxy serverPoolType = iota
	// TransparentInspection FortiWeb inspect traffic asynchronously. It does not modify traffic
	TransparentInspection serverPoolType = iota
	// WCCP Web Cache Communication Protocol: Provides web caching with load balancing and fault tolerance
	WCCP serverPoolType = iota
)

const (
	_ = iota
	// RoundRobin ...
	RoundRobin loadBalancingAlgorithm = iota
	// WeightedRoundRobin ...
	WeightedRoundRobin loadBalancingAlgorithm = iota
	// LeastConnection ...
	LeastConnection loadBalancingAlgorithm = iota
	// URIHash ...
	URIHash loadBalancingAlgorithm = iota
	// FullURIHash ...
	FullURIHash loadBalancingAlgorithm = iota
	// HostHash ...
	HostHash loadBalancingAlgorithm = iota
	// HostDomainHash ...
	HostDomainHash loadBalancingAlgorithm = iota
	// SourceIPHash ...
	SourceIPHash loadBalancingAlgorithm = iota
)

// CreateServerPool creates a virtual server pool object in FortiWeb
// Simplifies POST operation to external user
func (f *FortiWebClient) CreateServerPool(name string,
	singleOrMultiple singleServerOrServerBalance,
	poolType serverPoolType,
	lbAlgorithm loadBalancingAlgorithm,
	comments string) error {

	body := map[string]interface{}{
		"name": name,
		"singleServerOrServerBalance": singleOrMultiple,
		"type":                   poolType,
		"comments":               comments,
		"loadBalancingAlgorithm": lbAlgorithm,
		"can_delete":             true,
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

// CreateServerPoolRule creates a virtual server pool object in FortiWeb
// Simplifies POST operation to external user
func (f *FortiWebClient) CreateServerPoolRule(serverPoolName string,
	ip string,
	port int,
	status int,
	connectionLimit int) error {

	body := map[string]interface{}{
		"ip":            ip,
		"status":        status,
		"port":          port,
		"connectLimit":  connectionLimit,
		"inHeritHCheck": true}

	jsonByte, err := json.Marshal(body)

	if err != nil {
		fmt.Printf("Error in json data: %s\n", err)
		return err
	}

	response, error := f.doPost(
		strings.Join([]string{"api/v1.0/ServerObjects/Server/ServerPool/",
			serverPoolName,
			"/EditServerPoolRule"}, ""),
		string(jsonByte))

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

type matchObject int

const (
	_                                    = iota
	httpHost                 matchObject = iota
	httpURL                  matchObject = iota
	urlParameter             matchObject = iota
	httpReferer              matchObject = iota
	httpCookie               matchObject = iota
	httpHeader               matchObject = iota
	sourceIP                 matchObject = iota
	x509CertificateSubject   matchObject = iota
	x509CertificateExtension matchObject = iota
	httpsSNI                 matchObject = iota
)

type concatenateOperator int

const (
	// AND is used to concatenate conditions in HTTP Content Routing
	AND concatenateOperator = 2
	// OR is used to concatenate conditions in HTTP Content Routing
	OR concatenateOperator = 3
)

// CreateHTTPContentRoutingUsingHost creates a criteria for matching http content in a policy
// Simplifies POST operation to external user
func (f *FortiWebClient) CreateHTTPContentRoutingUsingHost(HTTPContentRoutingPolicy string,
	matchExpression string,
	hostCondition int,
	concatenate concatenateOperator) error {

	body := map[string]interface{}{
		"matchObject":     httpHost,
		"matchExpression": matchExpression,
		"hostCondition":   hostCondition,
		"concatenate":     concatenate,
	}

	jsonByte, err := json.Marshal(body)

	if err != nil {
		fmt.Printf("Error in json data: %s\n", err)
		return err
	}

	url := strings.Join([]string{"api/v1.0/ServerObjects/Server/HTTPContentRoutingPolicy/",
		HTTPContentRoutingPolicy,
		"/HTTPContentRoutingPolicyNewHTTPContentRouting"},
		"")
	response, error := f.doPost(url, string(jsonByte))

	if error != nil || response.StatusCode != 200 {
		fmt.Printf("The HTTP request failed with error %s, %d, %s\n", error, response.StatusCode, response.Status)
		return error
	}

	return nil
}

// CreateHTTPContentRoutingUsingURL creates a criteria for matching http content in a policy
// Simplifies POST operation to external user
func (f *FortiWebClient) CreateHTTPContentRoutingUsingURL(HTTPContentRoutingPolicy string,
	matchExpression string,
	urlCondition int,
	concatenate concatenateOperator) error {

	body := map[string]interface{}{
		"matchObject":     httpURL,
		"matchExpression": matchExpression,
		"urlCondition":    urlCondition,
		"concatenate":     concatenate,
	}

	jsonByte, err := json.Marshal(body)

	if err != nil {
		fmt.Printf("Error in json data: %s\n", err)
		return err
	}

	url := strings.Join([]string{"api/v1.0/ServerObjects/Server/HTTPContentRoutingPolicy/",
		HTTPContentRoutingPolicy,
		"/HTTPContentRoutingPolicyNewHTTPContentRouting"},
		"")
	response, error := f.doPost(url, string(jsonByte))

	if error != nil || response.StatusCode != 200 {
		fmt.Printf("The HTTP request failed with error %s, %d, %s\n", error, response.StatusCode, response.Status)
		return error
	}

	return nil
}

type deploymentMode string

const (
	// HTTPContentRouting sets deployment mode to use http headers for routing
	HTTPContentRouting deploymentMode = "http_content_routing"
	// SingleServerOrServerPool set deployment mode to steer traffic to each node of the pool
	SingleServerOrServerPool deploymentMode = "server_pool"
)

// CreateServerPolicy creates a criteria for matching http content in a policy
// Simplifies POST operation to external user
func (f *FortiWebClient) CreateServerPolicy(name,
	virtualServer,
	protectedHostnames,
	httpService,
	httpsService,
	protectionProfile,
	comments string,
	deployment deploymentMode,
	halfOpenThreshold int,
	clientRealIP,
	synCookie,
	httpRedirectHTTPS,
	monitorMode,
	urlCaseSensitivity bool) error {

	body := map[string]interface{}{
		"name":               name,
		"depInMode":          deployment,
		"disdmode":           "HTTP Content Routing",
		"virtualServer":      virtualServer,
		"HTTPService":        httpService,
		"HTTPSService":       httpsService,
		"clientRealIP":       clientRealIP,
		"halfopenThresh":     halfOpenThreshold,
		"syncookie":          synCookie,
		"hRedirectoHttps":    httpRedirectHTTPS,
		"MonitorMode":        monitorMode,
		"URLCaseSensitivity": urlCaseSensitivity,
		"comments":           comments,
		"enable":             true,
		"status":             "run"}

	if protectedHostnames != "" {
		body["protectedHostnames"] = protectedHostnames
	}

	if protectionProfile != "" {
		body["protectionProfile"] = protectionProfile
	}

	if deployment == HTTPContentRouting {
		body["disdmode"] = "HTTP Content Routing"
	} else {
		body["disdmode"] = "Single Server/Server Pool"
	}

	jsonByte, err := json.Marshal(body)

	if err != nil {
		fmt.Printf("Error in json data: %s\n", err)
		return err
	}

	response, error := f.doPost("api/v1.0/Policy/ServerPolicy/ServerPolicy/", string(jsonByte))

	if error != nil || response.StatusCode != 200 {
		fmt.Printf("The HTTP request failed with error %s, %d, %s\n", error, response.StatusCode, response.Status)
		return error
	}

	return nil
}

//doPost is internal function to apply a generic POST operation to FortiWeb
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
