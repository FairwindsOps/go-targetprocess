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

// Team matches up with a targetprocess Team
type Team struct {
	client *Client

	ID              int32         `json:"Id,omitempty"`
	Name            string        `json:",omitempty"`
	Description     string        `json:",omitempty"`
	StartDate       DateTime      `json:",omitempty"`
	EndDate         DateTime      `json:",omitempty"`
	CreateDate      DateTime      `json:",omitempty"`
	ModifyDate      DateTime      `json:",omitempty"`
	NumericPriority float64       `json:",omitempty"`
	CustomFields    []CustomField `json:",omitempty"`
	Abbreviation    string        `json:",omitempty"`
}

// TeamResponse is a representation of the http response for a group of Teams
type TeamResponse struct {
	Items []Team
	Next  string
	Prev  string
}

// GetTeam will return a single team based on its name. If somehow there are teams with the same name,
// this will only return the first one.
func (c *Client) GetTeam(name string) (Team, error) {
	ret := Team{}
	out := TeamResponse{}
	err := c.Get(&out, "Team", nil,
		Where(fmt.Sprintf("Name eq '%s'", name)),
		First(),
	)
	if err != nil {
		return ret, errors.Wrap(err, fmt.Sprintf("error getting team with name '%s'", name))
	}
	ret = out.Items[0]
	ret.client = c
	return ret, nil
}

// NewUserStory will make a UserStory for assigned to the Team that this method is built off of and for the given Project
func (t Team) NewUserStory(name, description, project string) (UserStory, error) {
	us := UserStory{
		client:      t.client,
		Name:        name,
		Description: description,
	}
	t.client.debugLog(fmt.Sprintf("Attempting to Get Project: %s", project))
	p, err := t.client.GetProject(project)
	if err != nil {
		return UserStory{}, err
	}
	us.Project = &p
	us.Team = &t
	return us, nil
}
