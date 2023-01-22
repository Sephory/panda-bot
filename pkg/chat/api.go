package chat

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
)

type ApiClient struct {
	httpClient *http.Client
	header     http.Header
	baseUrl    *url.URL
}

func NewApiClient(baseUrl string, header http.Header) (*ApiClient, error) {
	url, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	return &ApiClient{
		httpClient: http.DefaultClient,
		header:     header,
		baseUrl:    url,
	}, nil
}

func (api *ApiClient) Get(path string, query map[string]string) (*http.Response, error) {
	request, err := api.buildRequest("GET", path, query, nil)
	if err != nil {
		return nil, err
	}
	return api.httpClient.Do(request)
}

func (api *ApiClient) Post(path string, query map[string]string, body interface{}) (*http.Response, error) {
	request, err := api.buildRequest("POST", path, query, body)
	if err != nil {
		return nil, err
	}

	return api.httpClient.Do(request)
}

func (api *ApiClient) Delete(path string, query map[string]string) (*http.Response, error) {
	request, err := api.buildRequest("DELETE", path, query, nil)
	if err != nil {
		return nil, err
	}

	return api.httpClient.Do(request)
}

func (api *ApiClient) buildRequest(method string, path string, query map[string]string, body interface{}) (*http.Request, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	url := api.baseUrl.JoinPath(path)
	for key, value := range query {
		url.Query().Add(key, value)
	}
	request, err := http.NewRequest(method, url.String(), bytes.NewBuffer(jsonBody))
	request.Header = api.header
	return request, nil
}
