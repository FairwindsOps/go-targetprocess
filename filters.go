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
	"net/url"
	"strconv"
	"strings"
)

// QueryFilter accepts a query and returns a query, modifying
// it in some way before sending
type QueryFilter func(r url.Values) (url.Values, error)

// Where is a QueryFilter that represents the `where` parameter
// in a url query.
func Where(queries ...string) QueryFilter {
	return func(values url.Values) (url.Values, error) {
		for _, query := range queries {
			if _, exists := values["where"]; exists {
				currentWhere := values.Get("where")
				if currentWhere == "" {
					values.Set("where", fmt.Sprintf("(%s)", query))
					continue
				}
				values.Set("where", currentWhere+" and "+fmt.Sprintf("(%s)", query))
			} else {
				values.Set("where", fmt.Sprintf("(%s)", query))
			}
		}
		return values, nil
	}
}

// First is a QueryFilter that will only return a single
// Entity. It returns the first one encountered.
func First() QueryFilter {
	return func(values url.Values) (url.Values, error) {
		values.Set("take", "1")
		return values, nil
	}
}

// Include is used to filter what fields are returned by a given Query
func Include(includes ...string) QueryFilter {
	return func(values url.Values) (url.Values, error) {
		values.Set("include", fmt.Sprintf("[%s]", strings.Join(includes, ",")))
		return values, nil
	}
}

// MaxPerPage is how many results we will allow to return per page. Default enforced by the API is 25.
func MaxPerPage(count int) QueryFilter {
	return func(values url.Values) (url.Values, error) {
		values.Set("take", strconv.Itoa(count))
		return values, nil
	}
}
