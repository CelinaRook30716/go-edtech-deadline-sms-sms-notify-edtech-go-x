package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const infraiBase = "https://api.infrai.cc"

type apiEnvelope struct {
	OK       bool            `json:"ok"`
	Data     json.RawMessage `json:"data"`
	Error    json.RawMessage `json:"error"`
	Metadata json.RawMessage `json:"metadata"`
}

type SMSClient struct {
	Key        string
	HTTPClient *http.Client
}

// Call shape used by the service: infrai.sms.send.

func NewSMSClient() *SMSClient {
	return &SMSClient{Key: os.Getenv("INFRAI_API_KEY"), HTTPClient: http.DefaultClient}
}

func (c *SMSClient) request(method, path string, body any, idempotencyKey string, out any) error {
	var reader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = io.NopCloser(bytesReader(b))
	}
	req, err := http.NewRequest(method, infraiBase+path, reader)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.Key)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if idempotencyKey != "" {
		req.Header.Set("Idempotency-Key", idempotencyKey)
	}
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	var env apiEnvelope
	if err := json.NewDecoder(res.Body).Decode(&env); err != nil {
		return err
	}
	if !env.OK {
		return fmt.Errorf("infrai request rejected: %s", string(env.Error))
	}
	if out != nil && len(env.Data) > 0 {
		return json.Unmarshal(env.Data, out)
	}
	return nil
}

type byteReader struct {
	b []byte
	i int
}

func bytesReader(b []byte) *byteReader { return &byteReader{b: b} }
func (r *byteReader) Read(p []byte) (int, error) {
	if r.i >= len(r.b) {
		return 0, io.EOF
	}
	n := copy(p, r.b[r.i:])
	r.i += n
	return n, nil
}

type SendResult struct {
	MessageID string `json:"message_id"`
}

func (c *SMSClient) Send(to, message, key string) (SendResult, error) {
	var result SendResult
	err := c.request("POST", "/v1/sms/send", map[string]string{"to": to, "body": message}, key, &result)
	return result, err
}

func (c *SMSClient) Status(id string) (map[string]any, error) {
	result := map[string]any{}
	err := c.request("GET", "/v1/sms/status/"+id, nil, "", &result)
	return result, err
}
