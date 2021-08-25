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

// Project matches up with a targetprocess Project
type Project struct {
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
	Effort          float32       `json:",omitempty"`
	EffortCompleted float32       `json:",omitempty"`
	EffortToDo      float32       `json:",omitempty"`
	IsActive        bool          `json:",omitempty"`
	Abbreviation    string        `json:",omitempty"`
	Color           string        `json:",omitempty"`
	Process         *Process      `json:",omitempty"`
}

// ProjectResponse is a representation of the http response for a group of Projects
type ProjectResponse struct {
	Items []Project
	Next  string
	Prev  string
}

// GetProjects will return all projects
func (c *Client) GetProjects(filters ...QueryFilter) ([]Project, error) {
	var ret []Project
	out := ProjectResponse{}

	err := c.Get(&out, "Project", nil, filters...)
	if err != nil {
		return nil, err
	}
	ret = append(ret, out.Items...)
	for out.Next != "" {
		innerOut := ProjectResponse{}
		err := c.GetNext(&innerOut, out.Next)
		if err != nil {
			return ret, err
		}
		ret = append(ret, innerOut.Items...)
		out = innerOut
	}
	return ret, nil
}

// GetProject will return a single project based on its name. If somehow there are projects with the same name,
// this will only return the first one.
func (c *Client) GetProject(name string) (Project, error) {
	ret := Project{}
	out := ProjectResponse{}
	err := c.Get(&out, "Project", nil,
		Where(fmt.Sprintf("Name == '%s'", name)),
		First(),
	)
	if err != nil {
		return Project{}, errors.Wrap(err, fmt.Sprintf("error getting project with name '%s'", name))
	}
	if len(out.Items) < 1 {
		return ret, fmt.Errorf("no items found")
	}
	ret = out.Items[0]
	ret.client = c
	return ret, nil
}

// NewFeature will make a Feature assigned to the Project that this method is built off of and for the given Team
func (p Project) NewFeature(name, description string) (Feature, error) {
	f := Feature{
		client:      p.client,
		Name:        name,
		Description: description,
	}
	f.Project = &p
	return f, nil
}

// NewUserStory will make a UserStory for assigned to the Project that this method is built off of and for the given Team
func (p Project) NewUserStory(name, description, team string) (UserStory, error) {
	us := UserStory{
		client:      p.client,
		Name:        name,
		Description: description,
	}
	p.client.debugLog(fmt.Sprintf("[targetprocess] Attempting to Get Team: %s", team))
	t, err := p.client.GetTeam(team)
	if err != nil {
		return UserStory{}, err
	}
	us.Project = &p
	us.Team = &t
	return us, nil
}

// GetProcess returns the process associated with a project
func (p Project) GetProcess() (*Process, error) {
	processList, err := p.client.GetProcesses(
		Where(fmt.Sprintf("Id == %d", p.Process.ID)),
	)
	if err != nil {
		return nil, err
	}
	if len(processList) != 1 {
		return nil, fmt.Errorf("cannot determine process for project. got list: %v", processList)
	}
	return &processList[0], nil
}

func (p *Project) GetRelease(name string) (Release, error) {
	return p.client.GetRelease(p, name)
}
