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

var mockProject1 = Project{
	ID:              1,
	Name:            "MockProject1",
	Description:     "Project 1 Description",
	Abbreviation:    "P1",
	NumericPriority: float64(1116),
	CreateDate:      "/Date(1553311554000-0500)/",
	ModifyDate:      "/Date(1563297019000-0500)/",
	Effort:          float32(58),
	EffortCompleted: float32(58),
}

var mockProject2 = Project{
	ID:              2,
	Name:            "MockProject2",
	Abbreviation:    "P2",
	NumericPriority: float64(1116),
	Description:     "Project 2 Description",
	CreateDate:      "/Date(1553311554000-0500)/",
	ModifyDate:      "/Date(1563297019000-0500)/",
	Effort:          float32(58),
	EffortCompleted: float32(58),
	// This is not a mistake. It would seem that the API only returns a stub of the process here, missing the description.
	Process: &Process{
		ID:   2,
		Name: "Kanban",
	},
}

var mockProjectResponse = `
{
	"Next": "https://testing.tpondemand.com/api/v1/Projects/?format=json&take=25&skip=25",
	"Items": [
	  {
		"ResourceType": "Project",
		"Id": 1,
		"Name": "MockProject1",
		"Description": "Project 1 Description",
		"StartDate": null,
		"EndDate": null,
		"CreateDate": "/Date(1553311554000-0500)/",
		"ModifyDate": "/Date(1563297019000-0500)/",
		"LastCommentDate": null,
		"Tags": "",
		"NumericPriority": 1116,
		"Effort": 58,
		"EffortCompleted": 58,
		"EffortToDo": 0,
		"IsActive": false,
		"IsProduct": false,
		"Abbreviation": "P1",
		"MailReplyAddress": null,
		"Color": null,
		"Progress": 1,
		"PlannedStartDate": null,
		"PlannedEndDate": null,
		"IsPrivate": null
	  }
	]
}`

var mockProjectResponseSingle = `
{
	"Items": [
	  {
		"ResourceType": "Project",
		"Id": 2,
		"Name": "MockProject2",
		"Description": "Project 2 Description",
		"StartDate": null,
		"EndDate": null,
		"CreateDate": "/Date(1553311554000-0500)/",
		"ModifyDate": "/Date(1563297019000-0500)/",
		"LastCommentDate": null,
		"Tags": "",
		"NumericPriority": 1116,
		"Effort": 58,
		"EffortCompleted": 58,
		"EffortToDo": 0,
		"IsActive": false,
		"IsProduct": false,
		"Abbreviation": "P2",
		"MailReplyAddress": null,
		"Color": null,
		"Progress": 1,
		"PlannedStartDate": null,
		"PlannedEndDate": null,
		"IsPrivate": null,
		"Process": {
		  "ResourceType": "Process",
		  "Id": 2,
		  "Name": "Kanban"
		}
	  }
	]
  }  
`

func TestClient_GetProjects(t *testing.T) {
	client := NewFakeClient()

	tests := []struct {
		name    string
		filters []QueryFilter
		want    []Project
		wantErr bool
	}{
		{
			name:    "simple",
			filters: nil,
			want: []Project{
				mockProject1,
				mockProject2,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.GetProjects(tt.filters...)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}
		})
	}
}

func TestProject_GetProcess(t *testing.T) {
	tests := []struct {
		name    string
		project Project
		want    *Process
		wantErr bool
	}{
		{
			name:    "simple",
			project: mockProject2,
			want:    &mockProcess2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.project.client = NewFakeClient()
			got, err := tt.project.GetProcess()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}
		})
	}
}

func TestProject_NewUserStory(t *testing.T) {
	client := NewFakeClient()
	tests := []struct {
		name        string
		project     Project
		description string
		team        string

		want    UserStory
		wantErr bool
	}{
		{
			name:        "simple",
			project:     mockProject2,
			description: "description",
			want: UserStory{
				client:      client,
				Team:        &mockTeam1,
				Project:     &mockProject2,
				Name:        "simple",
				Description: "description",
			},
			team:    "MockTeam1",
			wantErr: false,
		},
		{
			name:        "team dne",
			project:     mockProject2,
			description: "failure",
			team:        "dne",
			want: UserStory{
				client:      client,
				Team:        &mockTeam1,
				Project:     &mockProject2,
				Name:        "simple",
				Description: "description",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Patch the clients
			tt.project.client = client
			tt.want.Team.client = client
			tt.want.Project.client = client

			got, err := tt.project.NewUserStory(tt.name, tt.description, tt.team)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}
		})
	}
}
