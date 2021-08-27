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

// Feature matches up with a targetprocess Feature
type Feature struct {
	client *Client

	ID               int32         `json:"Id,omitempty"`
	Name             string        `json:",omitempty"`
	Effort           float32       `json:",omitempty"`
	UserStoriesCount int64         `json:"UserStories-Count,omitempty"`
	Project          *Project      `json:",omitempty"`
	Description      string        `json:",omitempty"`
	NumericPriority  float32       `json:",omitempty"`
	CustomFields     []CustomField `json:",omitempty"`
}

// FeatureResponse is a representation of the http response for a group of Features
type FeatureResponse struct {
	Items []Feature
	Next  string
	Prev  string
}

// NewFeature creates a Feature struct with the required fields of
// name, description, and project.
func NewFeature(c *Client, name, description, project string) (Feature, error) {
	f := Feature{
		client:      c,
		Name:        name,
		Description: description,
	}
	c.debugLog(fmt.Sprintf("[targetprocess] Attempting to Get Project: %s", project))
	p, err := c.GetProject(project)
	if err != nil {
		return Feature{}, err
	}
	f.Project = &p
	return f, nil
}

// GetFeatures will return all features
func (c *Client) GetFeatures(filters ...QueryFilter) ([]Feature, error) {
	var ret []Feature
	out := FeatureResponse{}

	err := c.Get(&out, "Feature", nil, filters...)
	if err != nil {
		return nil, err
	}
	ret = append(ret, out.Items...)
	for out.Next != "" {
		innerOut := FeatureResponse{}
		err := c.GetNext(&innerOut, out.Next)
		if err != nil {
			return ret, err
		}
		ret = append(ret, innerOut.Items...)
		out = innerOut
	}
	return ret, nil
}

// GetFeature will return a single feature based on its name. If somehow there are features with the same name,
// this will only return the first one.
func (c *Client) GetFeature(name string) (Feature, error) {
	c.debugLog(fmt.Sprintf("[targetprocess] Attempting to Get Feature: %s", name))
	ret := Feature{}
	out := FeatureResponse{}
	err := c.Get(&out, "Feature", nil,
		Where(fmt.Sprintf("Name == '%s'", name)),
		First(),
	)
	if err != nil {
		return Feature{}, errors.Wrap(err, fmt.Sprintf("error getting feature with name '%s'", name))
	}
	if len(out.Items) < 1 {
		return ret, fmt.Errorf("no items found")
	}
	ret = out.Items[0]
	ret.client = c
	return ret, nil
}

// NewUserStory will make a UserStory with the Feature that this method is built off of
func (f Feature) NewUserStory(name, description, project string) (UserStory, error) {
	us, err := NewUserStory(f.client, name, description, project)
	if err != nil {
		return UserStory{}, nil
	}
	us.Feature = &f
	return us, nil
}

// Create takes a Feature struct and crafts a POST request to TP and sends it
// it returns the ID of the Feature created as well as a link to the entity
// on the Target Process frontend
func (f Feature) Create() (int32, string, error) {
	client := f.client
	resp := &struct {
		ID int32 `json:"Id"`
	}{}
	body, err := json.Marshal(f)
	if err != nil {
		return 0, "", errors.Wrap(err, fmt.Sprintf("error marshaling POST body for Feature %s", f.Name))
	}

	client.debugLog(fmt.Sprintf("Attempting to POST Feature: %+v", f))
	err = client.Post(resp, "Feature", nil, body)
	if err != nil {
		return 0, "", errors.Wrap(err, fmt.Sprintf("error POSTing Feature %s", f.Name))
	}
	client.debugLog("[targetprocess] Successfully POSTed Feature")
	client.debugLog(fmt.Sprintf("[targetprocess] Feature created. ID: %d", resp.ID))
	link := GenerateURL(client.account, resp.ID)
	return resp.ID, link, nil
}
