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

// EntityState contains metadata for the state of an Entity. Collection of EntityStates
// form Workflow for Entity. For example, Bug has four EntityStates by default: Open, Fixed, Invalid and Done
type EntityState struct {
	ID                int32        `json:"Id,omitempty"`
	Name              string       `json:",omitempty"`
	NumericPriority   float64      `json:",omitempty"`
	ParentEntityState *EntityState `json:",omitempty"`
	Process           *Process     `json:",omitempty"`
	IsInital          bool         `json:",omitempty"`
	IsFinal           bool         `json:",omitempty"`
	IsPlanned         bool         `json:",omitempty"`
	IsCommentRequired bool         `json:",omitempty"`
}

// EntityStateResponse is a representation of the http response for a group of EntityStates
type EntityStateResponse struct {
	Items []EntityState
	Next  string
	Prev  string
}

// GetEntityStates will return all EntityStates
func (c *Client) GetEntityStates(filters ...QueryFilter) ([]EntityState, error) {
	var ret []EntityState
	out := EntityStateResponse{}

	err := c.Get(&out, "EntityState", nil, filters...)
	if err != nil {
		return nil, err
	}
	ret = append(ret, out.Items...)
	for out.Next != "" {
		innerOut := EntityStateResponse{}
		err := c.GetNext(&innerOut, out.Next)
		if err != nil {
			return ret, err
		}
		ret = append(ret, innerOut.Items...)
		out = innerOut
	}
	return ret, nil
}
