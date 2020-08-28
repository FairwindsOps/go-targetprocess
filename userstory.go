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
	Feature             *Feature       `json:",omitempty"`
}

// UserStoryList is a list of user stories. Can be used to create multiple stories at once
type UserStoryList struct {
	client  *Client
	Stories []UserStory `json:"Stories"`
}

// UserStoryResponse is a representation of the http response for a group of UserStories
type UserStoryResponse struct {
	Items []UserStory
	Next  string
	Prev  string
}

// NewUserStory creates a new UserStory with the required fields of
// name, description, and project.
// If more fields are required, use the Add<Field> method of UserStory to add them
func NewUserStory(c *Client, name, description, project string) (UserStory, error) {
	us := UserStory{
		client:      c,
		Name:        name,
		Description: description,
	}
	err := us.SetProject(project)
	if err != nil {
		return UserStory{}, err
	}

	return us, nil
}

// SetProject sets the Project field for a user story
func (us *UserStory) SetProject(project string) error {
	us.client.debugLog(fmt.Sprintf("Attempting to Get Team: %s", project))
	p, err := us.client.GetProject(project)
	if err != nil {
		return err
	}
	us.Project = &p
	return nil
}

// SetTeam sets the Team field for a user story
func (us *UserStory) SetTeam(team string) error {
	us.client.debugLog(fmt.Sprintf("Attempting to Get Team: %s", team))
	t, err := us.client.GetTeam(team)
	if err != nil {
		return err
	}
	us.Team = &t
	return nil
}

// SetFeature sets the Feature field for a user story
func (us *UserStory) SetFeature(feature string) error {
	f, err := us.client.GetFeature(feature)
	if err != nil {
		return err
	}
	us.Feature = &f
	return nil
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

// NewUserStoryList returns a UserStoryList from a list of user stories.
// Used for batch POSTing of UserStories
func (c *Client) NewUserStoryList(list []UserStory) *UserStoryList {
	ret := &UserStoryList{
		client:  c,
		Stories: list,
	}
	return ret
}

// Create posts a list of user stories to create them
func (usl UserStoryList) Create() ([]int32, error) {
	client := usl.client
	resp := &UserStoryResponse{}
	body, err := json.Marshal(usl.Stories)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("error marshaling POST body for UserStoryList %v", usl))
	}
	client.debugLog(fmt.Sprintf("Attempting to POST UserStory: %+v", usl))
	err = client.Post(resp, "UserStories/bulk", nil, body)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("error POSTing UserStoryList %v", usl))
	}
	client.debugLog("Successfully POSTed UserStoryList")

	ret := []int32{}
	for _, story := range resp.Items {
		ret = append(ret, story.ID)
	}
	client.debugLog(fmt.Sprintf("User stories created with IDs: %v", ret))
	return ret, nil
}
