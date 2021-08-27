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

// Release matches up with a targetprocess Release
type Release struct {
	client *Client

	ID               int32         `json:"Id,omitempty"`
	Name             string        `json:",omitempty"`
	Effort           float32       `json:",omitempty"`
	Description      string        `json:",omitempty"`
	NumericPriority  float32       `json:",omitempty"`
	CustomFields     []CustomField `json:",omitempty"`
	
	Project          *Project      `json:",omitempty"`
}

// ReleaseResponse is a representation of the http response for a group of Releases
type ReleaseResponse struct {
	Items []Release
	Next  string
	Prev  string
}

// NewRelease creates a Release struct with the required fields of
// name, description, and project.
func NewRelease(c *Client, name, description, project string) (Release, error) {
	r := Release{
		client:      c,
		Name:        name,
		Description: description,
	}
	c.debugLog(fmt.Sprintf("[targetprocess] Attempting to Get Project: %s", project))
	p, err := c.GetProject(project)
	if err != nil {
		return Release{}, err
	}
	r.Project = &p
	return r, nil
}

// GetReleases will return all releases
func (c *Client) GetReleases(filters ...QueryFilter) ([]Release, error) {
	var ret []Release
	out := ReleaseResponse{}

	err := c.Get(&out, "Release", nil, filters...)
	if err != nil {
		return nil, err
	}
	ret = append(ret, out.Items...)
	for out.Next != "" {
		innerOut := ReleaseResponse{}
		err := c.GetNext(&innerOut, out.Next)
		if err != nil {
			return ret, err
		}
		ret = append(ret, innerOut.Items...)
		out = innerOut
	}
	return ret, nil
}

// GetRelease will return a single release based on its name. If somehow there are releases with the same name,
// this will only return the first one.
func (c *Client) GetRelease(projectName string, name string) (Release, error) {
	c.debugLog(fmt.Sprintf("[targetprocess] Attempting to Get Project: %s", projectName))
	p, err := c.GetProject(projectName)
	if err != nil {
		return Release{}, errors.Wrap(err, fmt.Sprintf("error getting release with name '%s' for project '%s'", name, projectName))
	}

	return c.getRelease(&p, name)
}

func (c *Client) getRelease(p *Project, name string) (Release, error) {
	c.debugLog(fmt.Sprintf("[targetprocess] Attempting to Get Release: %s, for Project: %s", name, p.Name))
	ret := Release{}
	out := ReleaseResponse{}
	err := c.Get(&out, "Release", nil,
		Where(fmt.Sprintf("Name == '%s'", name), fmt.Sprintf("Project.Id == %d", p.ID)),
		First(),
	)
	if err != nil {
		return Release{}, errors.Wrap(err, fmt.Sprintf("error getting release with name '%s' for project '%s'", name, p.Name))
	}
	if len(out.Items) < 1 {
		return ret, fmt.Errorf("no items found")
	}
	ret = out.Items[0]
	ret.client = c
	return ret, nil
}

// NewUserStory will make a UserStory with the Release that this method is built off of
func (r Release) NewUserStory(name, description, project string) (UserStory, error) {
	us, err := NewUserStory(r.client, name, description, project)
	if err != nil {
		return UserStory{}, nil
	}
	us.Release = &r
	return us, nil
}

// Create takes a Release struct and crafts a POST request to TP and sends it
// it returns the ID of the Release created as well as a link to the entity
// on the Target Process frontend
func (r Release) Create() (int32, string, error) {
	client := r.client
	resp := &struct {
		ID int32 `json:"Id"`
	}{}
	body, err := json.Marshal(r)
	if err != nil {
		return 0, "", errors.Wrap(err, fmt.Sprintf("error marshaling POST body for Release %s", r.Name))
	}

	client.debugLog(fmt.Sprintf("Attempting to POST Release: %+v", r))
	err = client.Post(resp, "Release", nil, body)
	if err != nil {
		return 0, "", errors.Wrap(err, fmt.Sprintf("error POSTing Release %s", r.Name))
	}
	client.debugLog("[targetprocess] Successfully POSTed Release")
	client.debugLog(fmt.Sprintf("[targetprocess] Release created. ID: %d", resp.ID))
	link := GenerateURL(client.account, resp.ID)
	return resp.ID, link, nil
}
