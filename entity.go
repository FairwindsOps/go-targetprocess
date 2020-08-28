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

// EntityState contains metadata for the state of an Entity. Collection of EntityStates
// form Workflow for Entity. For example, Bug has four EntityStates by default: Open, Fixed, Invalid and Done
type EntityState struct {
	client *Client

	ID              int32   `json:"Id,omitempty"`
	Name            string  `json:",omitempty"`
	NumericPriority float64 `json:",omitempty"`
}

// EntityStateResponse is a representation of the http response for a group of EntityStates
type EntityStateResponse struct {
	Items []EntityState
	Next  string
	Prev  string
}

// GetEntityState will return an EntityState object from a name. Returns an error if not found.
func (c *Client) GetEntityState(name string) (EntityState, error) {
	c.debugLog(fmt.Sprintf("attempting to get EntityState: %s", name))
	ret := EntityState{}
	out := EntityStateResponse{}
	err := c.Get(&out, "EntityState", nil,
		Where(fmt.Sprintf("Name eq '%s'", name)),
		First(),
	)
	if err != nil {
		return EntityState{}, errors.Wrap(err, fmt.Sprintf("error getting EntityState with name '%s'", name))
	}
	if len(out.Items) < 1 {
		return ret, fmt.Errorf("no items found")
	}
	ret = out.Items[0]
	ret.client = c
	return ret, nil
}
