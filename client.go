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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	userAgent = "go-targetprocess"
)

var (
	defaultClient *http.Client
)

func init() {
	defaultClient = http.DefaultClient
}

// Client is the API client for Targetprocess. Create this using NewClient.
// This can also be constructed manually but it isn't recommended.
type Client struct {
	// baseURL is the base URL for v1 API requests.
	baseURL *url.URL

	// baseURLReadOnly is the base URL for v2 API requests.
	baseURLReadOnly *url.URL

	// Client is the HTTP client to use for communication.
	Client *http.Client

	// Logger is an optional logging interface for debugging
	Logger logger

	// Token is the user access token to authenticate to the Targetprocess instance
	Token string

	// Timeout is the timeout used in any request made
	Timeout time.Duration

	// UserAgent is the user agent to send with API requests
	UserAgent string

	ctx context.Context
}

type logger interface {
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
}

// NewClient will create a new Targetprocess client
//
// account is the name of your Targetprocess account and will be the first part of your
// Targetprocess url. ex:
//   https://exampleaccount.tpondemand.com
//
// token is your user access token taken from your account settings
// see here: https://dev.targetprocess.com/docs/authentication#token-authentication
func NewClient(account, token string) (*Client, error) {
	c := defaultClient
	c.Timeout = 15 * time.Second
	baseURLString := fmt.Sprintf("https://%s.tpondemand.com/api/v1/", account)
	baseURLReadOnlyString := fmt.Sprintf("https://%s.tpondemand.com/api/v2/", account)
	baseURL, err := url.Parse(baseURLString)
	if err != nil {
		return nil, err
	}
	baseURLReadOnly, err := url.Parse(baseURLReadOnlyString)
	if err != nil {
		return nil, err
	}
	return &Client{
		baseURL:         baseURL,
		baseURLReadOnly: baseURLReadOnly,
		Client:          c,
		Token:           token,
		UserAgent:       userAgent,
		ctx:             context.Background(),
	}, nil
}

// WithContext takes a context.Context, sets it as context on the client and returns
// a Client pointer.
func (c *Client) WithContext(ctx context.Context) {
	c.ctx = ctx
}

// Get is a generic HTTP GET call to the targetprocess api passing in the type of entity and any query filters
func (c *Client) Get(out interface{}, entityType string, values url.Values, filters ...QueryFilter) error {
	rel, err := url.Parse(entityType + "/")
	if err != nil {
		return errors.Wrapf(err, "Error parsing entity type: %s", entityType)
	}
	u := c.baseURLReadOnly.ResolveReference(rel)

	if values == nil {
		values = url.Values{}
	}

	for _, filter := range filters {
		values, err = filter(values)
		if err != nil {
			return errors.Wrap(err, "Error running query filter")
		}
	}
	values = c.defaultParams(values)

	c.debugLog("[targetprocess] GET %s%s?%s", c.baseURLReadOnly, entityType, values.Encode())
	fullURL := fmt.Sprintf("%s?%s", u.String(), values.Encode())
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return errors.Wrapf(err, "Invalid GET request: %s/%s", c.baseURLReadOnly, entityType)
	}
	return c.do(out, req, entityType)
}

// GetNext is a helper method to get the next page of results from a query.
func (c *Client) GetNext(out interface{}, nextURL string) error {
	prevFull, err := url.Parse(nextURL)
	if err != nil {
		return errors.Wrapf(err, "Invalid Next URL: %s", nextURL)
	}

	// The v1 and v2 API return differently formatted Next urls, so we need to be sure the Path ends with "/"
	if !strings.HasSuffix(prevFull.Path, "/") {
		prevFull.Path += "/"
	}

	splitPath := strings.Split(prevFull.EscapedPath(), "/")
	entityType := splitPath[len(splitPath)-2]
	entityURLType, err := url.Parse(entityType + "/")
	if err != nil {
		return errors.Wrapf(err, "Invalid Next URL Entity Type: %s", entityURLType)
	}

	return c.Get(out, entityType, prevFull.Query())
}

// Post is for both creating and updating objects in TargetProcess
func (c *Client) Post(out interface{}, entityType string, values url.Values, body []byte) error {
	rel, err := url.Parse(entityType + "/")
	if err != nil {
		return errors.Wrapf(err, "Error parsing entity type: %s", entityType)
	}
	u := c.baseURL.ResolveReference(rel)

	if values == nil {
		values = url.Values{}
	}
	values = c.defaultParams(values)

	c.debugLog("[targetprocess] POST %s/%s?%s", c.baseURL, entityType, values.Encode())
	fullURL := fmt.Sprintf("%s?%s", u.String(), values.Encode())

	req, err := http.NewRequest("POST", fullURL, bytes.NewBuffer(body))
	if err != nil {
		return errors.Wrapf(err, "Invalid POST request: %s/%s", c.baseURL, entityType)
	}
	return c.do(out, req, entityType)
}

func (c *Client) do(out interface{}, req *http.Request, urlPath string) error {
	noParameterURL := fmt.Sprintf("%s://%s%s", req.URL.Scheme, req.URL.Host, req.URL.Path)

	// Set the headers that will be required for every request
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	if c.UserAgent != "" {
		req.Header.Add("User-Agent", c.UserAgent)
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "HTTP request failure on %s", noParameterURL)
	}

	// Empty the body and close it to reuse the Transport
	defer func() {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return makeHTTPClientError(urlPath, resp)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, "HTTP Read error on response for %s", urlPath)
	}
	c.debugLog(fmt.Sprintf("[targetprocess] raw response: %s", string(b)))
	err = json.Unmarshal(b, out)
	if err != nil {
		return errors.Wrapf(err, "JSON decode failed on %s:\n%s", urlPath, string(b))
	}
	return nil
}

func (c *Client) defaultParams(v url.Values) url.Values {
	if c.Token != "" {
		v.Add("accessToken", c.Token)
	}
	v.Set("format", "json")
	v.Set("resultFormat", "json")
	return v
}

func (c *Client) debugLog(format string, args ...interface{}) {
	if c.Logger != nil {
		c.Logger.Debugf(format, args...)
	}
}

func (c *Client) infoLog(format string, args ...interface{}) { // nolint:golint,unused
	if c.Logger != nil {
		c.Logger.Infof(format, args...)
	}
}
