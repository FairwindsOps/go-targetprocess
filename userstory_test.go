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

var mockUserStory1 = UserStory{
	ID:              12345,
	Name:            "Some user story name",
	Description:     "<!--markdown-->Some user story description",
	StartDate:       "/Date(1599065504000-0500)/",
	EndDate:         "",
	CreateDate:      "/Date(1599065469000-0500)/",
	ModifyDate:      "/Date(1599065505000-0500)/",
	NumericPriority: 13340.313369934804,
	CustomFields: []CustomField{
		{
			Name:  "Estimated Effort",
			Type:  "Number",
			Value: float64(0),
		},
	},
	Project: &Project{
		ID:   1,
		Name: "MockProject1",
	},
	LastStateChangeDate: "/Date(1599065504000-0500)/",
	Team: &Team{
		ID:   1,
		Name: "MockTeam1",
	},
}

var mockUserStory2 = UserStory{
	ID:              12346,
	Name:            "Some user story name",
	Description:     "<!--markdown-->Some user story description",
	StartDate:       "/Date(1599065504000-0500)/",
	EndDate:         "",
	CreateDate:      "/Date(1599065469000-0500)/",
	ModifyDate:      "/Date(1599065505000-0500)/",
	NumericPriority: 13340.313369934804,
	CustomFields: []CustomField{
		{
			Name: "Work Category",
			Type: "DropDown",
		},
	},
	Project: &Project{
		ID:   2,
		Name: "MockProject2",
		Process: &Process{
			ID: 2,
		},
	},
	LastStateChangeDate: "/Date(1599065504000-0500)/",
	Team: &Team{
		ID:   1,
		Name: "MockTeam1",
	},
}

var mockUserStoryResponse = `
{
	"Next": "https://testing.tpondemand.com/api/v1/UserStories/?take=25&skip=25",
	"Items": [
	  {
		"ResourceType": "UserStory",
		"Id": 12345,
		"Name": "Some user story name",
		"Description": "<!--markdown-->Some user story description",
		"StartDate": "/Date(1599065504000-0500)/",
		"EndDate": null,
		"CreateDate": "/Date(1599065469000-0500)/",
		"ModifyDate": "/Date(1599065505000-0500)/",
		"LastCommentDate": null,
		"Tags": "",
		"NumericPriority": 13340.313369934804,
		"Effort": 0,
		"EffortCompleted": 0,
		"EffortToDo": 0,
		"Progress": 0,
		"TimeSpent": 0,
		"TimeRemain": 0,
		"LastStateChangeDate": "/Date(1599065504000-0500)/",
		"PlannedStartDate": null,
		"PlannedEndDate": null,
		"InitialEstimate": 0,
		"Units": "pt",
		"EntityType": {
		  "ResourceType": "EntityType",
		  "Id": 4,
		  "Name": "UserStory"
		},
		"Project": {
		  "ResourceType": "Project",
		  "Id": 1,
		  "Name": "MockProject1"
		},
		"LastCommentedUser": null,
		"LinkedTestPlan": null,
		"Milestone": null,
		"Release": null,
		"Iteration": null,
		"TeamIteration": null,
		"Team": {
		  "ResourceType": "Team",
		  "Id": 1,
		  "Name": "MockTeam1"
		},
		"Priority": {
		  "ResourceType": "Priority",
		  "Id": 5,
		  "Name": "Nice To Have",
		  "Importance": 5
		},
		"Feature": null,
		"Build": null,
		"CustomFields": [
		  {
			"Name": "Estimated Effort",
			"Type": "Number",
			"Value": 0
		  }
		]
	  }
	]
  }  
`
var mockUserStoriesResponse = `
{
	"Items": [
	  {
		"ResourceType": "UserStory",
		"Id": 12346,
		"Name": "Some user story name",
		"Description": "<!--markdown-->Some user story description",
		"StartDate": "/Date(1599065504000-0500)/",
		"EndDate": null,
		"CreateDate": "/Date(1599065469000-0500)/",
		"ModifyDate": "/Date(1599065505000-0500)/",
		"LastCommentDate": null,
		"Tags": "",
		"NumericPriority": 13340.313369934804,
		"Effort": 0,
		"EffortCompleted": 0,
		"EffortToDo": 0,
		"Progress": 0,
		"TimeSpent": 0,
		"TimeRemain": 0,
		"LastStateChangeDate": "/Date(1599065504000-0500)/",
		"PlannedStartDate": null,
		"PlannedEndDate": null,
		"InitialEstimate": 0,
		"Units": "pt",
		"EntityType": {
		  "ResourceType": "EntityType",
		  "Id": 4,
		  "Name": "UserStory"
		},
		"Project": {
		  "ResourceType": "Project",
		  "Id": 2,
		  "Name": "MockProject2",
		  "Process": {
			"ResourceType": "Process",
			"Id": 2
		  }
		},
		"LastCommentedUser": null,
		"LinkedTestPlan": null,
		"Milestone": null,
		"Release": null,
		"Iteration": null,
		"TeamIteration": null,
		"Team": {
		  "ResourceType": "Team",
		  "Id": 1,
		  "Name": "MockTeam1"
		},
		"CustomFields": [
		  {
			"Name": "Work Category",
			"Type": "DropDown",
			"Value": null
		  }
		]
	  }
	]
  }  
`

func TestClient_GetUserStories(t *testing.T) {

	tests := []struct {
		name    string
		filters []QueryFilter
		want    []UserStory
		wantErr bool
	}{
		{
			name:    "simple",
			filters: nil,
			want: []UserStory{
				mockUserStory1,
				mockUserStory2,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewFakeClient()
			got, err := c.GetUserStories(tt.filters...)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}
		})
	}
}

func TestNewUserStory(t *testing.T) {
	client := NewFakeClient()
	type args struct {
		c           *Client
		name        string
		description string
		project     string
	}
	tests := []struct {
		name    string
		args    args
		want    UserStory
		wantErr bool
	}{
		{
			name: "simple",
			args: args{
				c:           client,
				name:        "test story",
				description: "test description",
				project:     "MockProject2",
			},
			want: UserStory{
				client:      client,
				Name:        "test story",
				Description: "test description",
				Project:     &mockProject2,
			},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				c:           client,
				name:        "test story",
				description: "test description",
				project:     "none",
			},
			want: UserStory{
				client:      client,
				Name:        "test story",
				Description: "test description",
				Project:     &mockProject1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUserStory(tt.args.c, tt.args.name, tt.args.description, tt.args.project)
			// We have to patch the mockProject with our client to get this to work.
			tt.want.Project.client = client
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}
		})
	}
}
