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

// Process contains metadata for the state of an Process. Collection of Processs
// form Process for Entity. For example, Bug has four Processs by default: Open, Fixed, Invalid and Done
type Process struct {
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
