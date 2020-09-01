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

var mockProcess1 = Process{
	ID:   1,
	Name: "Onboarding",
}

var mockProcess2 = Process{
	ID:          2,
	Name:        "Kanban",
	Description: "Kanban Board",
}

var mockProcessResponse = `
{
	"Next": "https://testing.tpondemand.com/api/v1/Processes/?take=25&skip=25",
	"Items": [
	  {
		"ResourceType": "Process",
		"Id": 1,
		"Name": "Onboarding",
		"IsDefault": false,
		"Description": null
	  }
	]
  }
  `

var mockProcessResponse2 = `
{
	"Items": [
	  {
		"ResourceType": "Process",
		"Id": 2,
		"Name": "Kanban",
		"IsDefault": false,
		"Description": "Kanban Board"
	  }
	]
  }
  `
var mockProcessResponse1 = `
  {
	  "Items": [
		{
		  "ResourceType": "Process",
		  "Id": 1,
		  "Name": "Onboarding",
		  "IsDefault": false,
		  "Description": null
		}
	  ]
	}
	`

func TestClient_GetProcesses(t *testing.T) {
	client := NewFakeClient()
	tests := []struct {
		name    string
		filters []QueryFilter
		want    []Process
		wantErr bool
	}{
		{
			name:    "simple",
			filters: nil,
			want: []Process{
				mockProcess1,
				mockProcess2,
			},
			wantErr: false,
		},
		{
			name: "filter",
			filters: []QueryFilter{
				Where("Id eq 2"),
			},
			want: []Process{
				mockProcess2,
			},
			wantErr: false,
		},
		{
			name: "empty",
			filters: []QueryFilter{
				Where("Id eq 4"),
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.GetProcesses(tt.filters...)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}
		})
	}
}
