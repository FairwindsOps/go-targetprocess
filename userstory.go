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

// UserStory matches up with a targetprocess UserStory
type UserStory struct {
	client *Client

	ID                  int32          `json:"Id,omitempty"`
	Name                string         `json:",omitempty"`
	Description         string         `json:",omitempty"`
	StartDate           DateTime       `json:",omitempty"`
	EndDate             DateTime       `json:",omitempty"`
	CreateDate          DateTime       `json:",omitempty"`
	ModifyDate          DateTime       `json:",omitempty"`
	NumericPriority     float64        `json:",omitempty"`
	CustomFields        []CustomField  `json:",omitempty"`
	Effort              float32        `json:",omitempty"`
	EffortCompleted     float32        `json:",omitempty"`
	EffortToDo          float32        `json:",omitempty"`
	Project             *Project       `json:",omitempty"`
	Progress            float32        `json:",omitempty"`
	TimeSpent           float32        `json:",omitempty"`
	TimeRemain          float32        `json:",omitempty"`
	LastStateChangeDate DateTime       `json:",omitempty"`
	InitialEstimate     float32        `json:",omitempty"`
	Team                *Team          `json:",omitempty"`
	EntityState         *EntityState   `json:",omitempty"`
	AssignedTeams       *AssignedTeams `json:",omitempty"`
}

// UserStoryResponse is a representation of the http response for a group of UserStories
type UserStoryResponse struct {
	Items []UserStory
	Next  string
	Prev  string
}

// NewUserStory creates a new UserStory with the required fields of
// name, description, and project
func NewUserStory(c *Client, name, description, project string) (UserStory, error) {
	us := UserStory{
		client:      c,
		Name:        name,
		Description: description,
	}
	p, err := c.GetProject(project)
	if err != nil {
		return UserStory{}, err
	}
	us.Project = &p
	return us, nil
}

// NewUserStoryForTeam is mostly the same as NewUserStory but assigns it to a team
func NewUserStoryForTeam(c *Client, name, description, project, team string) (UserStory, error) {
	us := UserStory{
		client:      c,
		Name:        name,
		Description: description,
	}
	c.debugLog(fmt.Sprintf("Attempting to Get Project: %s", project))
	p, err := c.GetProject(project)
	if err != nil {
		return UserStory{}, err
	}
	c.debugLog(fmt.Sprintf("Attempting to Get Team: %s", team))
	t, err := c.GetTeam(team)
	if err != nil {
		return UserStory{}, err
	}
	us.Project = &p
	us.Team = &t
	return us, nil
}

// GetUserStories will return all user stories
//
// Use with caution if you have a lot and are not setting the MaxPerPage to a high number
// as it could cause a lot of requests to the API and may take a long time.
// If you know you have a lot you may want to include the QueryFilter MaxPerPage
func (c *Client) GetUserStories(filters ...QueryFilter) ([]UserStory, error) {
	var ret []UserStory
	out := UserStoryResponse{}

	err := c.Get(&out, "UserStory", nil, filters...)
	if err != nil {
		return nil, err
	}
	ret = append(ret, out.Items...)
	for out.Next != "" {
		innerOut := UserStoryResponse{}
		err := c.GetNext(&innerOut, out.Next)
		if err != nil {
			return ret, err
		}
		ret = append(ret, innerOut.Items...)
		out = innerOut
	}
	return ret, nil
}

// Create takes a UserStory struct and crafts a POST to make it so in TP
// it returns the ID of the UserStory created
func (us UserStory) Create() (int32, error) {
	client := us.client
	resp := &struct {
		ID int32 `json:"Id"`
	}{}
	body, err := json.Marshal(us)
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("error marshaling POST body for UserStory %s", us.Name))
	}

	client.debugLog(fmt.Sprintf("Attempting to POST UserStory: %+v", us))
	err = client.Post(resp, "UserStory", nil, body)
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("error POSTing UserStory %s", us.Name))
	}
	client.debugLog("Successfully POSTed UserStory")
	client.debugLog(fmt.Sprintf("UserStory created. ID: %d", resp.ID))
	return resp.ID, nil
}
