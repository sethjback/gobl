package httpapi

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// APIRequest defines a request
type APIRequest struct {
	Address   string
	Body      []byte
	Signature string
}

// POST sends a post request to the given endpoint
func (r *APIRequest) POST(endpoint string) (*APIResponse, error) {

	url := "http://" + r.Address + endpoint

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(r.Body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-gobl-signature", r.Signature)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var response APIResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return nil, &response
	}

	return &response, nil
}

// GET a reuqets to the agent
func (r *APIRequest) GET(endpoint string) (*APIResponse, error) {
	url := "http://" + r.Address + endpoint

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var response APIResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, &response
	}

	return &response, nil
}
