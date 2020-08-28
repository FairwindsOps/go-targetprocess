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

// Workflow contains metadata for the state of a Workflow.
type Workflow struct {
	ID      int32    `json:"Id,omitempty"`
	Name    string   `json:",omitempty"`
	Process *Process `json:",omitempty"`
}

// WorkflowResponse is a representation of the http response for a group of Workflows
type WorkflowResponse struct {
	Items []Workflow
	Next  string
	Prev  string
}

// GetWorkflows will return all Workflows
func (c *Client) GetWorkflows(filters ...QueryFilter) ([]Workflow, error) {
	var ret []Workflow
	out := WorkflowResponse{}

	err := c.Get(&out, "Workflow", nil, filters...)
	if err != nil {
		return nil, err
	}
	ret = append(ret, out.Items...)
	for out.Next != "" {
		innerOut := WorkflowResponse{}
		err := c.GetNext(&innerOut, out.Next)
		if err != nil {
			return ret, err
		}
		ret = append(ret, innerOut.Items...)
		out = innerOut
	}
	return ret, nil
}
