package httpclient

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/kosha/dna-center/pkg/logger"
	"io/ioutil"
	"net/http"
)

var token string

type AuthToken struct {
	Token string `json:"Token,omitempty"`
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func makeHttpBasicAuthReq(username, password string, method, url string, body interface{}, log logger.Logger, isSecure bool) ([]byte, int) {

	var req *http.Request
	if body != nil {
		jsonReq, _ := json.Marshal(body)
		req, _ = http.NewRequest(method, url, bytes.NewBuffer(jsonReq))
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}

	req.Header.Set("Authorization", "Basic "+basicAuth(username, password))

	req.Header.Set("Accept-Encoding", "identity")

	client := &http.Client{}

	if !isSecure {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client.Transport = tr
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Error(err)
		return nil, 500
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}
	return bodyBytes, resp.StatusCode
}

func makeHttpApiKeyReq(apiKeyHeaderName, apiKey string, method, url string, body interface{}, log logger.Logger, isSecure bool) ([]byte, int) {

	var req *http.Request
	if body != nil {
		jsonReq, _ := json.Marshal(body)
		req, _ = http.NewRequest(method, url, bytes.NewBuffer(jsonReq))
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}

	if apiKeyHeaderName != "" {
		req.Header.Set(apiKeyHeaderName, apiKey)
	} else {
		// if there is no accompanying header name, assume it is the Authorization header that needs to be sent
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	req.Header.Set("Accept-Encoding", "identity")

	client := &http.Client{}

	if !isSecure {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client.Transport = tr
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Error(err)
		return nil, 500
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}
	return bodyBytes, resp.StatusCode
}

func MakeHttpCall(headers map[string]string, username, password, method, serverUrl, url string, body interface{}, log logger.Logger, isSecure bool) (interface{}, int, error) {

	var response interface{}
	var payloadRes []byte

	var statusCode int

	if token != "" {
		payloadRes, statusCode = makeHttpApiKeyReq("X-Auth-Token", token, method, url, body, log, isSecure)
		if string(payloadRes) == "" {
			return nil, statusCode, fmt.Errorf("nil")
		}
		// Convert response body to target struct
		err := json.Unmarshal(payloadRes, &response)
		if err != nil {
			log.Error("Unable to parse response as json")
			log.Error(err)
			return nil, 500, err
		}
		if statusCode == 200 && response != nil {
			return response, statusCode, nil
		}
	}
	// token is not generated, or is invalid so get new token
	token, _ = getToken(username, password, serverUrl, log, isSecure)
	if token == "" {
		return nil, 500, fmt.Errorf("error generating token")
	}

	var newResponse interface{}
	payloadResponse, statusCode := makeHttpApiKeyReq("X-Auth-Token", token, method, url, body, log, isSecure)
	if string(payloadResponse) == "" {
		return nil, statusCode, fmt.Errorf("nil")
	}
	// Convert response body to target struct
	err := json.Unmarshal(payloadResponse, &newResponse)
	if err != nil {
		log.Error("Unable to parse response as json")
		log.Error(err)
		return nil, 500, err
	}

	return newResponse, statusCode, nil
}

func getToken(username, password, serverUrl string, log logger.Logger, isSecure bool) (string, int) {

	var tokenResponse AuthToken

	url := serverUrl + "/dna/system/api/v1/auth/token"

	res, _ := makeHttpBasicAuthReq(username, password, "POST", url, nil, log, isSecure)
	if string(res) == "" {
		return "", 500
	}
	// Convert response body to target struct
	err := json.Unmarshal(res, &tokenResponse)
	if err != nil {
		log.Error("Unable to parse auth token response as json")
		log.Error(err)
		return "", 500
	}
	return tokenResponse.Token, 200
}
