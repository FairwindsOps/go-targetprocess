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
	"fmt"
	"io/ioutil"
	"net/http"
)

type notFoundError interface {
	IsNotFound() bool
}

type permissionDeniedError interface {
	IsPermissionDenied() bool
}

type httpClientError struct {
	msg  string
	code int
}

func makeHTTPClientError(url string, resp *http.Response) error {

	body, _ := ioutil.ReadAll(resp.Body)
	msg := fmt.Sprintf("HTTP request failure on %s:\n%d: %s", url, resp.StatusCode, string(body))

	return &httpClientError{
		msg:  msg,
		code: resp.StatusCode,
	}
}

func (e *httpClientError) Error() string            { return e.msg }
func (e *httpClientError) IsNotFound() bool         { return e.code == 404 }
func (e *httpClientError) IsPermissionDenied() bool { return e.code == 401 }

// IsNotFound takes an error and returns true if the error is exactly a not-found error.
func IsNotFound(err error) bool {
	nf, ok := err.(notFoundError)
	return ok && nf.IsNotFound()
}

// IsPermissionDenied takes an error and returns true if the error is exactly a
// permission-denied error.
func IsPermissionDenied(err error) bool {
	pd, ok := err.(permissionDeniedError)
	return ok && pd.IsPermissionDenied()
}
