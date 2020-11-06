// Copyright 2020 Fairwinds
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License

package targetprocess

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	okResponse = `{
		"items": [],
        "next": "",
        "prev": ""
	}`
)

type genericResponse struct {
	Items []string `json:"items"`
	Next  string   `json:"next"`
	Prev  string   `json:"prev"`
}

func newMockClient(handler http.Handler, account, token string) (*Client, func()) {
	s := httptest.NewTLSServer(handler)

	defaultClient = &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, s.Listener.Addr().String())
			},
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	client, _ := NewClient(account, token)

	return client, s.Close
}

// Example of client.Get() for godoc
func ExampleClient_Get() {
	tpClient, err := NewClient("exampleaccount", "superSecretToken")
	if err != nil {
		fmt.Println("Failed to create tp client:", err)
		os.Exit(1)
	}
	var response = UserResponse{}
	err = tpClient.Get(response,
		"User",
		nil)
	if err != nil {
		fmt.Println("Failed to get users:", err)
		os.Exit(1)
	}
	jsonBytes, _ := json.Marshal(response)
	fmt.Print(string(jsonBytes))
}

func TestGet(t *testing.T) {
	tests := []struct {
		name        string
		account     string
		accessToken string
		entity      string
		path        string
	}{
		{
			name:        "Users GET",
			account:     "example",
			accessToken: "1234abcd",
			entity:      "Users",
		},
	}

	for _, tt := range tests {
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			queryParams := r.URL.Query()
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, fmt.Sprintf("%s.tpondemand.com", tt.account), r.Host)
			assert.Equal(t, fmt.Sprintf("/api/v2/%s/", tt.entity), r.URL.Path)
			assert.Equal(t, tt.accessToken, queryParams.Get("accessToken"))
			assert.Equal(t, "json", queryParams.Get("format"))
			assert.Equal(t, "json", queryParams.Get("resultFormat"))
			assert.Equal(t, "go-targetprocess", r.Header.Get("User-Agent"))
			assert.Equal(t, "application/json", r.Header.Get("Accept"))
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			_, _ = w.Write([]byte(okResponse))
		})
		mockClient, teardown := newMockClient(h, tt.account, tt.accessToken)

		resp := new(genericResponse)
		err := mockClient.Get(resp, tt.entity, nil)
		teardown()
		if err != nil {
			t.Logf("error sending get: %s", err)
			t.Fail()
		}
	}
}

func TestPost(t *testing.T) {
	tests := []struct {
		name        string
		account     string
		accessToken string
		entity      string
		path        string
	}{
		{
			name:        "Users POST",
			account:     "example",
			accessToken: "1234abcd",
			entity:      "UserStory",
		},
	}

	for _, tt := range tests {
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			queryParams := r.URL.Query()
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, fmt.Sprintf("%s.tpondemand.com", tt.account), r.Host)
			assert.Equal(t, fmt.Sprintf("/api/v1/%s/", tt.entity), r.URL.Path)
			assert.Equal(t, tt.accessToken, queryParams.Get("accessToken"))
			assert.Equal(t, "json", queryParams.Get("format"))
			assert.Equal(t, "json", queryParams.Get("resultFormat"))
			assert.Equal(t, "go-targetprocess", r.Header.Get("User-Agent"))
			assert.Equal(t, "application/json", r.Header.Get("Accept"))
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			_, _ = w.Write([]byte(okResponse))
		})
		mockClient, teardown := newMockClient(h, tt.account, tt.accessToken)

		resp := new(genericResponse)
		err := mockClient.Post(resp, tt.entity, nil, nil)
		teardown()
		if err != nil {
			t.Logf("error sending get: %s", err)
			t.Fail()
		}
	}
}
