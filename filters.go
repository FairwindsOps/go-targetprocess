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
)

// QueryFilter accepts a query and returns a query, modifying
// it in some way before sending
type QueryFilter func(r url.Values) (url.Values, error)

// First is a QueryFilter that will only return a single
// Entity. It returns the first one encountered.
//
// **WARNING** if using this with a helper function that supports
// automatic paging, ensure to set page to false (e.g. client.GetUserStories())
func First() QueryFilter {
	return func(values url.Values) (url.Values, error) {
		values.Set("take", "1")
		return values, nil
	}
}

// MaxPerPage is how many results we will allow to return per page. Default enforced by the API is 25.
// If a negative number is passed in it will be converted to a positive number. In the TP API, a negative
// number has the same return as a positive, so we'll
func MaxPerPage(count int) QueryFilter {
	return func(values url.Values) (url.Values, error) {
		if count < 0 {
			count = count * -1
		}
		values.Set("take", strconv.Itoa(count))
		return values, nil
	}
}

// Result is a QueryFilter that represents the `result` parameter
// in a url query. It is used to do custom calculations over the
// entire result set, such as getting the average Effort value
// across multiple items.
//
// It is important to note that this changes the output json and
// therefore you will need to adjust your receiving struct
func Result(query string) QueryFilter {
	return func(values url.Values) (url.Values, error) {
		values.Set("result", fmt.Sprintf("{%s}", query))
		return values, nil
	}
}

// Select is a QueryFilter that represents the `select` parameter
// in a url query. It is used to determine what fields are returned.
// It is important that the struct you are casting your results into
// accept the fields you specify.
func Select(query string) QueryFilter {
	return func(values url.Values) (url.Values, error) {
		values.Set("select", fmt.Sprintf("{%s}", query))
		return values, nil
	}
}

// Where is a QueryFilter that represents the `where` parameter
// in a url query.
func Where(queries ...string) QueryFilter {
	return func(values url.Values) (url.Values, error) {
		for _, query := range queries {
			if _, exists := values["where"]; exists {
				currentWhere := values.Get("where")
				if currentWhere == "" {
					values.Set("where", query)
					continue
				}
				values.Set("where", currentWhere+" and "+query)
			} else {
				values.Set("where", query)
			}
		}
		return values, nil
	}
}
