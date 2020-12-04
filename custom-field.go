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

	"github.com/pkg/errors"
)

// CustomField are user defined in TargetProcess and are arbitrarily set for many different
// resource types
type CustomField struct {
	client *Client

	ID    int32       `json:"Id,omitempty"`
	Type  string      `json:",omitempty"`
	Name  string      `json:",omitempty"`
	Value interface{} `json:",omitempty"`
}

// CustomFieldResponse is a representation of the http response for a group of CustomField
type CustomFieldResponse struct {
	Items []CustomField
	Next  string
	Prev  string
}

// GetCustomField will return a CustomField object from a name. Returns an error if not found.
func (c *Client) GetCustomField(name string) (CustomField, error) {
	c.debugLog(fmt.Sprintf("[targetprocess] attempting to get CustomField: %s", name))
	ret := CustomField{}
	out := CustomFieldResponse{}
	err := c.Get(&out, "CustomField", nil,
		Where(fmt.Sprintf("Name == '%s'", name)),
		First(),
	)
	if err != nil {
		return CustomField{}, errors.Wrap(err, fmt.Sprintf("error getting CustomField with name '%s'", name))
	}
	if len(out.Items) < 1 {
		return ret, fmt.Errorf("no items found")
	}
	ret = out.Items[0]
	ret.client = c
	return ret, nil
}
