package boostai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Client to communicate with Boost AI
type Client struct {
	BaseURL string
}

// Response response of http call
type Response struct {
	Status string
	Body   []byte
	Header http.Header
}

// Message - a message in a conversation
type Message struct {
	ConversationID string
}

// MessageResponse response of a message call
type MessageResponse struct {
	Conversation conversation `json:"conversation"`
	Response     response     `json:"response"`
}

// NewClient - Constructor for Boost AI Client
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
	}
}

// GetMessageTexts - Get all the messages as a list of strings
func (messageResponse *MessageResponse) GetMessageTexts() []string {
	elements := messageResponse.Response.Elements
	texts := make([]string, len(elements))
	for i, element := range elements {
		texts[i] = element.Payload.Text
	}
	return texts
}

// GetMessageText - Get all the messages as a single message
func (messageResponse *MessageResponse) GetMessageText() string {
	return strings.Join(messageResponse.GetMessageTexts(), " ")
}

// StartConversation - start the conversation
func (client *Client) StartConversation() (*MessageResponse, error) {
	return client.sendCommand(`{"command":"START", "clean": true}`)
}

// SendMessage - Send message to the bot
func (client *Client) SendMessage(message string, conversationID string) (*MessageResponse, error) {
	data := fmt.Sprintf(`{"command":"POST", "clean": true, "type": "text", "conversation_id": "%s", "value": "%s"}`, conversationID, message)
	return client.sendCommand(data)
}

// Private

type conversation struct {
	ID string `json:"id"`
}
type payloadElement struct {
	Text string `json:"text"`
}
type element struct {
	Payload payloadElement `json:"payload"`
}
type response struct {
	AvatarURL string    `json:"avatar_url"`
	Elements  []element `json:"elements"`
}

func (client *Client) sendCommand(data string) (*MessageResponse, error) {
	response := &MessageResponse{}
	err := deserialize(client.post("/api/chat/v2", []byte(data)), response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func deserialize(resp *Response, target interface{}) error {
	return json.Unmarshal(resp.Body, target)
}

func (client *Client) post(path string, data []byte) *Response {
	url := fmt.Sprintf("%s%s", client.BaseURL, path)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{Timeout: 10 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return &Response{
		Status: resp.Status,
		Body:   body,
		Header: resp.Header,
	}
}
