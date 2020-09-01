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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Example of client.Get() for godoc
func ExampleClient_Get() {
	tpClient := NewClient("exampleaccount", "superSecretToken")
	var response = UserResponse{}
	err := tpClient.Get(response,
		"User",
		nil)
	if err != nil {
		fmt.Println("Failed to get users:", err)
		os.Exit(1)
	}
	jsonBytes, _ := json.Marshal(response)
	fmt.Print(string(jsonBytes))
}

// Helpers for fake client
type roundTripFunc func(req *http.Request) *http.Response

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

var mockEmptyResponse = `
{
	"Items": []
}
`
var stringBreaksURL = string(byte(0x7f))

// NewFakeClient returns a client that will respond to very specific requests in order to build out unit tests
// Responses are based on a set of basic JSON objects built from actual TargetProcess queries
func NewFakeClient() *Client {
	c := NewClient("testing", "abc123")
	c.Client = &http.Client{
		Transport: roundTripFunc(func(r *http.Request) *http.Response {
			resp := &http.Response{
				StatusCode: 200,
			}
			path := r.URL.Path
			if r.Method == http.MethodGet {
				switch path {
				case "/api/v1/UserStory/":
					resp.Body = ioutil.NopCloser(strings.NewReader(mockUserStoryResponse))
				case "/api/v1/UserStories/":
					resp.Body = ioutil.NopCloser(strings.NewReader(mockUserStoriesResponse))
				case "/api/v1/Project/":
					if query, ok := r.URL.Query()["where"]; ok {
						if query[0] == "(Name eq 'MockProject2')" {
							resp.Body = ioutil.NopCloser(strings.NewReader(mockProjectResponseSingle))
							return resp
						}
						resp.Body = ioutil.NopCloser(strings.NewReader(mockEmptyResponse))
						return resp
					}
					resp.Body = ioutil.NopCloser(strings.NewReader(mockProjectResponse))
				case "/api/v1/Projects/":
					resp.Body = ioutil.NopCloser(strings.NewReader(mockProjectResponseSingle))
				case "/api/v1/Process/":
					if query, ok := r.URL.Query()["where"]; ok {
						switch query[0] {
						case "(Id eq 1)":
							resp.Body = ioutil.NopCloser(strings.NewReader(mockProcessResponse1))
							return resp
						case "(Id eq 2)":
							resp.Body = ioutil.NopCloser(strings.NewReader(mockProcessResponse2))
							return resp
						default:
							resp.Body = ioutil.NopCloser(strings.NewReader(mockEmptyResponse))
							return resp
						}
					}
					resp.Body = ioutil.NopCloser(strings.NewReader(mockProcessResponse))
				case "/api/v1/Processes/":
					resp.Body = ioutil.NopCloser(strings.NewReader(mockProcessResponse2))
				case "/api/v1/Team/":
					if query, ok := r.URL.Query()["where"]; ok {
						switch query[0] {
						case "(Name eq 'MockTeam1')":
							resp.Body = ioutil.NopCloser(strings.NewReader(mockTeamResponse1))
							return resp
						default:
							resp.Body = ioutil.NopCloser(strings.NewReader(mockEmptyResponse))
							return resp
						}
					}
					resp.Body = ioutil.NopCloser(strings.NewReader(mockTeamResponse))
				case "/api/v1/Teams/":
					resp.Body = ioutil.NopCloser(strings.NewReader(mockTeamResponse2))
				case "/api/v1/Feature/":
					if query, ok := r.URL.Query()["where"]; ok {
						switch query[0] {
						case "(Name eq 'MockFeature1')":
							resp.Body = ioutil.NopCloser(strings.NewReader(mockFeatureResponse1))
							return resp
						case "(Name eq 'MockFeature2')":
							resp.Body = ioutil.NopCloser(strings.NewReader(mockFeatureResponse2))
							return resp
						default:
							resp.Body = ioutil.NopCloser(strings.NewReader(mockEmptyResponse))
							return resp
						}
					}
					resp.Body = ioutil.NopCloser(strings.NewReader(mockFeatureResponse))
				case "/api/v1/Features/":
					resp.Body = ioutil.NopCloser(strings.NewReader(mockFeatureResponse2))
				case "/api/v1/User/":
					if query, ok := r.URL.Query()["where"]; ok {
						switch query[0] {
						case "(Id eq 1)":
							resp.Body = ioutil.NopCloser(strings.NewReader(mockUserResponse1))
							return resp
						case "(Id eq 2)":
							resp.Body = ioutil.NopCloser(strings.NewReader(mockUserResponse2))
							return resp
						default:
							resp.Body = ioutil.NopCloser(strings.NewReader(mockEmptyResponse))
							return resp
						}
					}
					resp.Body = ioutil.NopCloser(strings.NewReader(mockUserResponse))
				case "/api/v1/Users/":
					resp.Body = ioutil.NopCloser(strings.NewReader(mockUserResponse2))
				case "/api/v1/BasicJSON/":
					resp.Body = ioutil.NopCloser(strings.NewReader(`{"one":"two"}`))
				case "/api/v1/BadJSON/":
					resp.Body = ioutil.NopCloser(strings.NewReader(``))
				default:
					resp.Body = ioutil.NopCloser(strings.NewReader(`One or more errors occurred.`))
					resp.StatusCode = 404
				}
			}

			if r.Method == http.MethodPost {
				switch path {
				case "/api/v1/UserStory/":
					resp.Body = ioutil.NopCloser(strings.NewReader(mockUserStoryResponse))
				case "/api/v1/BasicJSON/":
					resp.Body = ioutil.NopCloser(strings.NewReader(`{"one":"two"}`))
				case "/api/v1/BadJSON/":
					resp.Body = ioutil.NopCloser(strings.NewReader(``))
				default:
					resp.Body = ioutil.NopCloser(strings.NewReader(`One or more errors occurred.`))
				}
			}

			return resp
		}),
	}
	return c
}

// Client Unit Tests
func TestGetWithBadUrl(t *testing.T) {
	c := NewClient("", "")
	target := map[string]interface{}{}

	err := c.Get(target, "project", nil)
	assert.EqualError(t, err, "HTTP request failure on https://.tpondemand.com/api/v1//project: Get \"https://.tpondemand.com/api/v1/project/?format=json&resultFormat=json\": dial tcp: lookup .tpondemand.com: no such host")
}

func TestClient_Get(t *testing.T) {
	tests := []struct {
		name       string
		want       interface{}
		entityType string
		values     []string

		// Error handling
		wantErr    bool
		errMessage string
	}{
		{
			name:       "basic json test",
			want:       map[string]string{"one": "two"},
			wantErr:    false,
			entityType: "BasicJSON",
		},
		{
			name:       "JSON decode failure error",
			wantErr:    true,
			errMessage: "JSON decode failed on BadJSON:\n: unexpected end of JSON input",
			entityType: "BadJSON",
		},
		{
			name:       "bad url format in entity type",
			wantErr:    true,
			entityType: stringBreaksURL,
			errMessage: "Error parsing entity type: \u007f: parse \"\\u007f/\": net/url: invalid control character in URL",
		},
		{
			name:       "return 400",
			wantErr:    true,
			errMessage: "HTTP request failure on :\n404: One or more errors occurred.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewFakeClient()
			out := make(map[string]string)
			values := url.Values{
				"test": tt.values,
			}
			err := c.Get(&out, tt.entityType, values)
			if tt.wantErr {
				assert.EqualErrorf(t, err, tt.errMessage, "expected an error")
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, out)
			}
		})
	}
}

func TestClient_Post(t *testing.T) {
	tests := []struct {
		name       string
		want       interface{}
		entityType string
		body       []byte

		// Error handling
		wantErr    bool
		errMessage string
	}{
		{
			name:       "basic json test",
			want:       map[string]string{"one": "two"},
			wantErr:    false,
			entityType: "BasicJSON",
		},
		{
			name:       "JSON decode failure error",
			wantErr:    true,
			errMessage: "JSON decode failed on BadJSON:\n: unexpected end of JSON input",
			entityType: "BadJSON",
		},
		{
			name:       "bad url format in entity type",
			wantErr:    true,
			entityType: stringBreaksURL,
			errMessage: "Error parsing entity type: \u007f: parse \"\\u007f/\": net/url: invalid control character in URL",
		},
		{
			name:       "return 404",
			wantErr:    true,
			errMessage: "JSON decode failed on :\nOne or more errors occurred.: invalid character 'O' looking for beginning of value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewFakeClient()
			out := make(map[string]string)

			err := c.Post(&out, tt.entityType, nil, tt.body)
			if tt.wantErr {
				assert.EqualErrorf(t, err, tt.errMessage, "expected an error")
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, out)
			}
		})
	}
}
