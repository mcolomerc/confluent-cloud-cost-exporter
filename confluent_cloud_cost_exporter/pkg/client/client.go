package client

import (
	"bytes"
	"log"

	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/mcolomerc/confluent_cost_exporter/config"
)

type HttpClient struct {
	Client *http.Client
	Bearer string
}

type Response struct {
	StatusCode int
	Body       []byte
}

type ResponseData struct {
	ApiVersion string        `json:"api_version"`
	Data       []interface{} `json:"data"`
}

func NewHttpClient(credentials config.Credentials) *HttpClient {
	user := credentials.Key + ":" + credentials.Secret
	bearer := b64.StdEncoding.EncodeToString([]byte(user))
	return &HttpClient{
		Bearer: bearer,
		Client: &http.Client{},
	}
}
func (c *HttpClient) Get(urlReq string) (interface{}, error) {
	log.Println("GET URL: ", urlReq)
	req, err := http.NewRequest("GET", urlReq, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Basic "+c.Bearer)
	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var response *map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	return response, nil

}
func (c *HttpClient) GetData(urlReq string) (*ResponseData, error) {
	log.Println("GET URL: ", urlReq)
	req, err := http.NewRequest("GET", urlReq, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Basic "+c.Bearer)
	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var response *ResponseData
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	return response, nil

}

func (c *HttpClient) Post(url string, body interface{}) (*ResponseData, error) {
	log.Println("POST URL: ", url)
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	log.Println("JSON BODY: ", string(jsonBody))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Basic "+c.Bearer)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response *ResponseData
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
