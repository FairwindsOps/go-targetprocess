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
	"testing"

	"github.com/stretchr/testify/assert"
)

var mockFeature1 = Feature{
	ID:   1,
	Name: "Feature1",
	Project: &Project{
		ID:   2,
		Name: "MockProject2",
	},
	NumericPriority: float32(13375),
}

var mockFeature2 = Feature{
	ID:   2,
	Name: "Feature2",
	Project: &Project{
		ID:   2,
		Name: "MockProject2",
	},
	NumericPriority: float32(13375),
}

var mockFeatureResponse = `
{
	"Next": "https://testing.tpondemand.com/api/v1/Features/?format=json&take=25&skip=25",
	"Items": [
	  {
		"ResourceType": "Feature",
		"Id": 1,
		"Name": "Feature1",
		"Description": null,
		"StartDate": null,
		"EndDate": null,
		"LastCommentDate": null,
		"Tags": "",
		"NumericPriority": 13375,
		"Effort": 0,
		"EffortCompleted": 0,
		"EffortToDo": 0,
		"Progress": 0,
		"TimeSpent": 0,
		"TimeRemain": 0,
		"PlannedStartDate": null,
		"PlannedEndDate": null,
		"InitialEstimate": 0,
		"Project": {
		  "ResourceType": "Project",
		  "Id": 2,
		  "Name": "MockProject2"
		}
	  }
	]
  }
`
var mockFeatureResponse1 = `
{
	"Items": [
	  {
		"ResourceType": "Feature",
		"Id": 1,
		"Name": "Feature1",
		"Description": null,
		"StartDate": null,
		"EndDate": null,
		"LastCommentDate": null,
		"Tags": "",
		"NumericPriority": 13375,
		"Effort": 0,
		"EffortCompleted": 0,
		"EffortToDo": 0,
		"Progress": 0,
		"TimeSpent": 0,
		"TimeRemain": 0,
		"PlannedStartDate": null,
		"PlannedEndDate": null,
		"InitialEstimate": 0,
		"Project": {
		  "ResourceType": "Project",
		  "Id": 2,
		  "Name": "MockProject2"
		}
	  }
	]
  }
`

var mockFeatureResponse2 = `
{
	"Items": [
	  {
		"ResourceType": "Feature",
		"Id": 2,
		"Name": "Feature2",
		"Description": null,
		"StartDate": null,
		"EndDate": null,
		"LastCommentDate": null,
		"Tags": "",
		"NumericPriority": 13375,
		"Effort": 0,
		"EffortCompleted": 0,
		"EffortToDo": 0,
		"Progress": 0,
		"TimeSpent": 0,
		"TimeRemain": 0,
		"PlannedStartDate": null,
		"PlannedEndDate": null,
		"InitialEstimate": 0,
		"Project": {
		  "ResourceType": "Project",
		  "Id": 2,
		  "Name": "MockProject2"
		}
	  }
	]
  }  
`

func TestClient_GetFeatures(t *testing.T) {
	client := NewFakeClient()

	tests := []struct {
		name    string
		filters []QueryFilter
		want    []Feature
		wantErr bool
	}{
		{
			name:    "no filter",
			filters: nil,
			want: []Feature{
				mockFeature1,
				mockFeature2,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.GetFeatures(tt.filters...)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}
		})
	}
}

func TestClient_GetFeature(t *testing.T) {
	client := NewFakeClient()

	tests := []struct {
		name    string
		feature string
		want    Feature
		wantErr bool
	}{
		{
			name:    "no error",
			feature: "MockFeature1",
			want:    mockFeature1,
			wantErr: false,
		},
		{
			name:    "not fou nd",
			feature: "dne",
			want:    Feature{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.GetFeature(tt.feature)
			// Patch client
			tt.want.client = client
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}
		})
	}
}
