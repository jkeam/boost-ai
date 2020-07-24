package boostai

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Client to communicate with Boost AI
type Client struct {
	BaseURL string
}

type Response struct {
	Status string
	Body   string
	Header http.Header
}

// NewClient - Constructor for Boost AI Client
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
	}
}

func (client *Client) StartConversation() {
}

// Private

func (client *Client) post(path string, data []byte) *Response {
	url := fmt.Sprintf("%s%s", client.BaseURL, path)

	// var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println("response Body:", string(body))
	// log.Debug(fmt.Sprintf("Base url: %s", baseURL))
	return &Response{
		Status: resp.Status,
		body:   string(body),
		Header: resp.Header,
	}
}
