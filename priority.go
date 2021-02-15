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

// Priority defines importance in targetprocess and can be assigned to multiple Assignable EntityTypes
type Priority struct {
	client *Client

	ID         int32  `json:"Id,omitempty"`
	Importance int32  `json:",omitempty"`
	Name       string `json:",omitempty"`
	IsDefault  bool   `json:",omitempty"`
}

// PriorityResponse is a representation of the http response for a group of Priority objects
type PriorityResponse struct {
	Items []Priority
	Next  string
	Prev  string
}

// GetPriority will return one Priority object by matching the name as well as the EntityType that it's assigned to (ex. UserStory)
func (c *Client) GetPriority(name, entityType string) (Priority, error) {
	c.debugLog(fmt.Sprintf("[targetprocess] attempting to get Priority: %s, for EntityType: %s", name, entityType))
	ret := Priority{}
	out := PriorityResponse{}
	err := c.Get(&out, "Priority", nil,
		Where(fmt.Sprintf("Name == '%s'", name)),
		Where(fmt.Sprintf("EntityType.Name == '%s'", entityType)),
		First(),
	)
	if err != nil {
		return Priority{}, errors.Wrap(err, fmt.Sprintf("error getting Priority with name '%s'", name))
	}
	if len(out.Items) < 1 {
		return ret, fmt.Errorf("no Priority found with the name: %s", name)
	}
	ret = out.Items[0]
	ret.client = c
	return ret, nil
}

// SetPriority assigns a priority to a UserStory by first finding the proper Priority in the TargetProcess API and then
// assigning it to the UserStory object
func (us *UserStory) SetPriority(priorityName string) error {
	priority, err := us.client.GetPriority(priorityName, "UserStory")
	if err != nil {
		return err
	}
	us.Priority = &priority
	return nil
}
