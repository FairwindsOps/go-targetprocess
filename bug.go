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

// Bug matches up with a targetprocess Bug
type Bug struct {
	client *Client

	ID                  int32           `json:"Id,omitempty"`
	Name                string          `json:",omitempty"`
	Description         string          `json:",omitempty"`
	StartDate           DateTime        `json:",omitempty"`
	EndDate             DateTime        `json:",omitempty"`
	CreateDate          DateTime        `json:",omitempty"`
	ModifyDate          DateTime        `json:",omitempty"`
	NumericPriority     float64         `json:",omitempty"`
	CustomFields        []CustomField   `json:",omitempty"`
	Effort              float32         `json:",omitempty"`
	EffortCompleted     float32         `json:",omitempty"`
	EffortToDo          float32         `json:",omitempty"`
	Project             *Project        `json:",omitempty"`
	Progress            float32         `json:",omitempty"`
	TimeSpent           float32         `json:",omitempty"`
	TimeRemain          float32         `json:",omitempty"`
	LastStateChangeDate DateTime        `json:",omitempty"`
	InitialEstimate     float32         `json:",omitempty"`
	Assignments         *Assignments    `json:",omitempty"`
	ResponsibleTeam     *TeamAssignment `json:",omitempty"`
	Team                *Team           `json:",omitempty"`
	Priority            *Priority       `json:",omitempty"`
	EntityState         *EntityState    `json:",omitempty"`
	AssignedUser        *AssignedUser   `json:",omitempty"`
	Feature             *Feature        `json:",omitempty"`
	Release             *Release        `json:",omitempty"`
}

// BugList is a list of bugs. Can be used to create multiple stories at once
type BugList struct {
	client  *Client
	Stories []Bug `json:"Bugs"`
}

// BugResponse is a representation of the http response for a group of Bugs
type BugResponse struct {
	Items []Bug
	Next  string
	Prev  string
}

// NewBug creates a new Bug with the required fields of
// name, description, and project.
// If more fields are required, use the Add<Field> method of Bug to add them
func NewBug(c *Client, name, description, project string) (Bug, error) {
	us := Bug{
		client:      c,
		Name:        name,
		Description: description,
	}
	err := us.SetProject(project)
	if err != nil {
		return Bug{}, err
	}

	return us, nil
}

// SetProject sets the Project field for a bug
func (us *Bug) SetProject(project string) error {
	us.client.debugLog(fmt.Sprintf("[targetprocess] Attempting to Get Project: %s", project))
	p, err := us.client.GetProject(project)
	if err != nil {
		return err
	}
	us.Project = &p
	return nil
}

// SetTeam sets the Team field for a bug
func (us *Bug) SetTeam(team string) error {
	us.client.debugLog(fmt.Sprintf("[targetprocess] Attempting to Get Team: %s", team))
	t, err := us.client.GetTeam(team)
	if err != nil {
		return err
	}
	us.Team = &t
	return nil
}

// SetFeature sets the Feature field for a bug
func (us *Bug) SetFeature(feature string) error {
	us.client.debugLog(fmt.Sprintf("[targetprocess] Attempting to Get Feature: %s", feature))
	f, err := us.client.GetFeature(feature)
	if err != nil {
		return err
	}
	us.Feature = &f
	return nil
}

// SetAssignedUserID assigns the Bug to a User based on their ID number
func (us *Bug) SetAssignedUserID(userID int32) {
	u := User{
		ID: userID,
	}
	au := Assignment{
		GeneralUser: &u,
	}
	assignments := Assignments{Items: []Assignment{au}}
	us.Assignments = &assignments
}

// GetBugs will return all bugs
//
// Use with caution if you have a lot and are not setting the MaxPerPage to a high number
// as it could cause a lot of requests to the API and may take a long time.
// If you know you have a lot you may want to include the QueryFilter MaxPerPage
//
func (c *Client) GetBugs(page bool, filters ...QueryFilter) ([]Bug, error) {
	var ret []Bug
	out := BugResponse{}

	err := c.Get(&out, "Bugs", nil, filters...)
	if err != nil {
		return nil, err
	}
	ret = append(ret, out.Items...)
	if page {
		for out.Next != "" {
			innerOut := BugResponse{}
			err := c.GetNext(&innerOut, out.Next)
			if err != nil {
				return ret, err
			}
			ret = append(ret, innerOut.Items...)
			out = innerOut
		}
	}
	return ret, nil
}

// Create takes a Bug struct and crafts a POST to make it so in TP
// it returns the ID of the Bug created as well as a link to the entity
// on the Target Process frontend
func (us Bug) Create() (int32, string, error) {
	client := us.client
	resp := &struct {
		ID int32 `json:"Id"`
	}{}
	body, err := json.Marshal(us)
	if err != nil {
		return 0, "", errors.Wrap(err, fmt.Sprintf("error marshaling POST body for Bug %s", us.Name))
	}

	client.debugLog(fmt.Sprintf("Attempting to POST Bug: %+v", us))
	err = client.Post(resp, "Bug", nil, body)
	if err != nil {
		return 0, "", errors.Wrap(err, fmt.Sprintf("error POSTing Bug %s", us.Name))
	}
	client.debugLog("[targetprocess] Successfully POSTed Bug")
	client.debugLog(fmt.Sprintf("[targetprocess] Bug created. ID: %d", resp.ID))
	link := GenerateURL(client.account, resp.ID)
	return resp.ID, link, nil
}

// NewBugList returns a BugList from a list of bugs.
// Used for batch POSTing of Bugs
func (c *Client) NewBugList(list []Bug) *BugList {
	return &BugList{
		client:  c,
		Stories: list,
	}
}

// Create posts a list of bugs to create them
// returns a list of entity IDs along with a list of links to them
func (usl BugList) Create() ([]int32, []string, error) {
	client := usl.client
	resp := &BugResponse{}
	body, err := json.Marshal(usl.Stories)
	if err != nil {
		return nil, nil, errors.Wrap(err, fmt.Sprintf("error marshaling POST body for BugList %v", usl))
	}
	client.debugLog(fmt.Sprintf("[targetprocess] Attempting to POST Bug: %+v", usl))
	err = client.Post(resp, "Bugs/bulk", nil, body)
	if err != nil {
		return nil, nil, errors.Wrap(err, fmt.Sprintf("error POSTing BugList %v", usl))
	}
	client.debugLog("[targetprocess] Successfully POSTed BugList")

	var (
		ret   []int32
		links []string
	)

	for _, bug := range resp.Items {
		ret = append(ret, bug.ID)
		links = append(links, GenerateURL(client.account, bug.ID))
	}
	client.debugLog(fmt.Sprintf("[targetprocess] Bugs created with IDs: %v", ret))
	return ret, links, nil
}
