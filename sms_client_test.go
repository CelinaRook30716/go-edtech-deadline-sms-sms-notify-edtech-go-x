package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestSendUsesBodyField(t *testing.T) {
	client := &SMSClient{Key: "test", HTTPClient: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		var payload map[string]string
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatal(err)
		}
		if len(payload) != 2 || payload["to"] != "+15550000000" || payload["body"] != "deadline reminder" {
			t.Fatalf("unexpected sms.send payload: %#v", payload)
		}
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(`{"ok":true,"data":{"message_id":"test-id"}}`)), Header: make(http.Header)}, nil
	})}}
	result, err := client.Send("+15550000000", "deadline reminder", "test-key")
	if err != nil || result.MessageID != "test-id" {
		t.Fatalf("Send() = %#v, %v", result, err)
	}
}
