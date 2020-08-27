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

// Feature matches up with a targetprocess Feature
type Feature struct {
	client *Client

	ID               int32   `json:"Id,omitempty"`
	Name             string  `json:",omitempty"`
	Effort           float32 `json:",omitempty"`
	UserStoriesCount float32 `json:"UserStories-Count",omitempty`
	Project          Project `json:",omitempty"`
}

// FeatureResponse is a representation of the http response for a group of Features
type FeatureResponse struct {
	Items []Feature
	Next  string
	Prev  string
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

// GetFeature will return a single feature based on its name. If somehow there are projects with the same name,
// this will only return the first one.
func (c *Client) GetFeature(name string) (Feature, error) {
	ret := Feature{}
	out := FeatureResponse{}
	err := c.Get(&out, "Feature", nil,
		Where(fmt.Sprintf("Name eq '%s'", name)),
		First(),
	)
	if err != nil {
		return Feature{}, errors.Wrap(err, fmt.Sprintf("error getting feature with name '%s'", name))
	}
	ret = out.Items[0]
	ret.client = c
	return ret, nil
}