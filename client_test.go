package boostai

import (
	"errors"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetMessageTexts(t *testing.T) {
	element1 := &element{Payload: &payloadElement{Text: "element1"}}
	element2 := &element{Payload: &payloadElement{Text: "element2"}}
	elements := []element{*element1, *element2}

	response := &response{
		AvatarURL: "http://avatar.com",
		Elements:  elements,
	}
	conversation := &conversation{ID: "1"}
	messageResponse := &MessageResponse{
		Conversation: *conversation,
		Response:     *response,
	}

	messageTexts := messageResponse.GetMessageTexts()
	assert.Equal(t, len(messageTexts), 2, "Should be two elements")
	assert.Equal(t, "element1", messageTexts[0], "First element should match")
	assert.Equal(t, "element2", messageTexts[1], "Second element should match")
}

func TestGetMessageText(t *testing.T) {
	element1 := &element{Payload: &payloadElement{Text: "element1"}}
	element2 := &element{Payload: &payloadElement{Text: "element2"}}
	elements := []element{*element1, *element2}

	response := &response{
		AvatarURL: "http://avatar.com/person.png",
		Elements:  elements,
	}
	conversation := &conversation{ID: "1"}
	messageResponse := &MessageResponse{
		Conversation: *conversation,
		Response:     *response,
	}

	messageText := messageResponse.GetMessageText()
	assert.Equal(t, "element1 element2", messageText, "Should match message text")
}

func TestStartConversation(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	client := NewClient("https://boost.ai")
	httpmock.RegisterResponder("POST", "https://boost.ai/api/chat/v2",
		httpmock.NewStringResponder(200, `{"conversation": {"id": "1"}, "response": {"avatar_url": "http://avatar.com/person.png", "elements": [{"payload": {"text": "Hi there"}}]}}`))

	messageResponse, err := client.StartConversation()
	if err != nil {
		assert.Fail(t, "Should not throw an error", err)
	}

	assert.Equal(t, 1, httpmock.GetTotalCallCount(), "Should only have been called once")
	assert.Equal(t, "1", messageResponse.Conversation.ID, "Should match conversation id")
	assert.Equal(t, "http://avatar.com/person.png", messageResponse.Response.AvatarURL, "Should match avatar url")
	assert.Equal(t, 1, len(messageResponse.Response.Elements), "Should be 1 element")
	assert.Equal(t, "Hi there", messageResponse.Response.Elements[0].Payload.Text, "Should match element text")
}

func TestStartConversationWithFilters(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	client := NewClient("https://boost.ai")
	httpmock.RegisterResponder("POST", "https://boost.ai/api/chat/v2",
		httpmock.NewStringResponder(200, `{"conversation": {"id": "1"}, "response": {"avatar_url": "http://avatar.com/person.png", "elements": [{"payload": {"text": "Hi there"}}]}}`))

	messageResponse, err := client.StartConversationWithFilters([]string{"filter"})
	if err != nil {
		assert.Fail(t, "Should not throw an error", err)
	}

	assert.Equal(t, 1, httpmock.GetTotalCallCount(), "Should only have been called once")
	assert.Equal(t, "1", messageResponse.Conversation.ID, "Should match conversation id")
	assert.Equal(t, "http://avatar.com/person.png", messageResponse.Response.AvatarURL, "Should match avatar url")
	assert.Equal(t, 1, len(messageResponse.Response.Elements), "Should be 1 element")
	assert.Equal(t, "Hi there", messageResponse.Response.Elements[0].Payload.Text, "Should match element text")
}

func TestSendMessage(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	client := NewClient("https://boost.ai")
	httpmock.RegisterResponder("POST", "https://boost.ai/api/chat/v2",
		httpmock.NewStringResponder(200, `{"conversation": {"id": "1"}, "response": {"avatar_url": "http://avatar.com/person.png", "elements": [{"payload": {"text": "Hi there"}}]}}`))

	messageResponse, err := client.SendMessage("Hi", "1")
	if err != nil {
		assert.Fail(t, "Should not throw an error", err)
	}

	assert.Equal(t, 1, httpmock.GetTotalCallCount(), "Should only have been called once")
	assert.Equal(t, "1", messageResponse.Conversation.ID, "Should match conversation id")
	assert.Equal(t, "http://avatar.com/person.png", messageResponse.Response.AvatarURL, "Should match avatar url")
	assert.Equal(t, 1, len(messageResponse.Response.Elements), "Should be 1 element")
	assert.Equal(t, "Hi there", messageResponse.Response.Elements[0].Payload.Text, "Should match element text")
}

func TestSendMessageFromPhone(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	client := NewClient("https://boost.ai")
	httpmock.RegisterResponder("POST", "https://boost.ai/api/chat/v2",
		httpmock.NewStringResponder(200, `{"conversation": {"id": "1"}, "response": {"avatar_url": "http://avatar.com/person.png", "elements": [{"payload": {"text": "Hi there"}}]}}`))

	messageResponse, err := client.SendMessageFromPhone("Hi", "1", "+15551112222")
	if err != nil {
		assert.Fail(t, "Should not throw an error", err)
	}

	assert.Equal(t, 1, httpmock.GetTotalCallCount(), "Should only have been called once")
	assert.Equal(t, "1", messageResponse.Conversation.ID, "Should match conversation id")
	assert.Equal(t, "http://avatar.com/person.png", messageResponse.Response.AvatarURL, "Should match avatar url")
	assert.Equal(t, 1, len(messageResponse.Response.Elements), "Should be 1 element")
	assert.Equal(t, "Hi there", messageResponse.Response.Elements[0].Payload.Text, "Should match element text")
}

func TestSendMessageFailParse(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	client := NewClient("https://boost.ai")
	httpmock.RegisterResponder("POST", "https://boost.ai/api/chat/v2", httpmock.NewStringResponder(200, ""))

	_, err := client.SendMessage("Hi", "1")
	if err == nil {
		assert.Fail(t, "Should throw an error", err)
	}

	assert.Equal(t, 1, httpmock.GetTotalCallCount(), "Should only have been called once")
}

func TestSendMessageErrorDuringRequest(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	client := NewClient("https://boost.ai")
	httpmock.RegisterResponder("POST", "https://boost.ai/api/chat/v2",
		func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("Request failure")
		},
	)

	_, err := client.SendMessage("Hi", "1")
	if err == nil {
		assert.Fail(t, "Should throw an error", err)
	}

	assert.Equal(t, 1, httpmock.GetTotalCallCount(), "Should only have been called once")
}
