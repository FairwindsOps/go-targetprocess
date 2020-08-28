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

	"github.com/pkg/errors"
)

// Process contains metadata for the state of an Process. Collection of Processs
// form Process for Entity. For example, Bug has four Processs by default: Open, Fixed, Invalid and Done
type Process struct {
	client *Client

	ID          int32  `json:"Id,omitempty"`
	Name        string `json:",omitempty"`
	Description string `json:",omitempty"`
}

// ProcessResponse is a representation of the http response for a group of Processs
type ProcessResponse struct {
	Items []Process
	Next  string
	Prev  string
}

// GetProcess will return an Process object from a name. Returns an error if not found.
func (c *Client) GetProcess(name string) (Process, error) {
	c.debugLog(fmt.Sprintf("attempting to get Process: %s", name))
	ret := Process{}
	out := ProcessResponse{}
	err := c.Get(&out, "Process", nil,
		Where(fmt.Sprintf("Name eq '%s'", name)),
		First(),
	)
	if err != nil {
		return Process{}, errors.Wrap(err, fmt.Sprintf("error getting Process with name '%s'", name))
	}
	if len(out.Items) < 1 {
		return ret, fmt.Errorf("no items found")
	}
	data, _ := json.Marshal(out.Items[0])
	c.debugLog(fmt.Sprintf("gotProcess: %s", string(data)))
	ret = out.Items[0]
	ret.client = c
	return ret, nil
}

// GetProcesses will return all Processs
func (c *Client) GetProcesses(filters ...QueryFilter) ([]Process, error) {
	var ret []Process
	out := ProcessResponse{}

	err := c.Get(&out, "Process", nil, filters...)
	if err != nil {
		return nil, err
	}
	ret = append(ret, out.Items...)
	for out.Next != "" {
		innerOut := ProcessResponse{}
		err := c.GetNext(&innerOut, out.Next)
		if err != nil {
			return ret, err
		}
		ret = append(ret, innerOut.Items...)
		out = innerOut
	}
	return ret, nil
}
