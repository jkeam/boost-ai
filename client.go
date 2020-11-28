package boostai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Client to communicate with Boost AI
type Client struct {
	BaseURL        string
	TimeoutSeconds int
}

// Response response of http call
type Response struct {
	Status string
	Body   []byte
	Header http.Header
}

// MessageResponse response of a message call
type MessageResponse struct {
	Conversation conversation `json:"conversation"`
	Response     response     `json:"response"`
}

// NewClient - Constructor for Boost AI Client
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL:        baseURL,
		TimeoutSeconds: 10,
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
	startCommand := newStartCommand()
	data, err := json.Marshal(startCommand)
	if err != nil {
		return nil, err
	}
	return client.sendCommand(string(data))
}

// StartConversationWithFilters - start the conversation, passing in the filter values
func (client *Client) StartConversationWithFilters(filterValues []string) (*MessageResponse, error) {
	startMessage := newStartCommand()
	startMessage.FilterValues = filterValues
	data, err := json.Marshal(startMessage)
	if err != nil {
		return nil, err
	}
	return client.sendCommand(string(data))
}

// SendMessage - Send message to the bot
func (client *Client) SendMessage(message string, conversationID string) (*MessageResponse, error) {
	postMessage := newPostCommand(message, conversationID)
	data, err := json.Marshal(postMessage)
	if err != nil {
		return nil, err
	}
	return client.sendCommand(string(data))
}

// SendMessageFromPhone - Send message to the bot
func (client *Client) SendMessageFromPhone(
	message string,
	conversationID string,
	fromPhone string,
) (*MessageResponse, error) {
	postMessage := newPostCommand(message, conversationID)
	postMessage.CustomMessagePayload = &customMessagePayload{ClientPhone: fromPhone}
	data, err := json.Marshal(postMessage)
	if err != nil {
		return nil, err
	}
	return client.sendCommand(string(data))
}

// Private

type conversation struct {
	ID string `json:"id"`
}

type payloadElement struct {
	Text string `json:"text"`
}

type element struct {
	Payload *payloadElement `json:"payload"`
}

type response struct {
	AvatarURL string    `json:"avatar_url"`
	Elements  []element `json:"elements"`
}

type customMessagePayload struct {
	ClientPhone string `json:"client_phone"`
}

type postCommand struct {
	Command              string                `json:"command"`
	Clean                bool                  `json:"clean"`
	Type                 string                `json:"type"`
	ConversationID       string                `json:"conversation_id"`
	Value                string                `json:"value"`
	CustomMessagePayload *customMessagePayload `json:"custom_payload,omitempty"`
}

type startCommand struct {
	Command      string   `json:"command"`
	Clean        bool     `json:"clean"`
	FilterValues []string `json:"filter_values,omitempty"`
}

func newStartCommand() *startCommand {
	return &startCommand{
		Command: "START",
		Clean:   true,
	}
}

func newPostCommand(message string, conversationID string) *postCommand {
	return &postCommand{
		Command:        "POST",
		Clean:          true,
		Type:           "text",
		ConversationID: conversationID,
		Value:          message,
	}
}

func (client *Client) sendCommand(data string) (*MessageResponse, error) {
	response, postError := client.post("/api/chat/v2", []byte(data))
	if postError != nil {
		return nil, postError
	}

	return response, nil
}

func (client *Client) post(path string, data []byte) (*MessageResponse, error) {
	url := fmt.Sprintf("%s%s", client.BaseURL, path)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{Timeout: time.Duration(client.TimeoutSeconds) * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response MessageResponse
	readErr := json.NewDecoder(resp.Body).Decode(&response)
	if readErr != nil {
		return nil, readErr
	}
	return &response, nil
}
