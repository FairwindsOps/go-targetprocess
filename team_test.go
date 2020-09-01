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

var mockTeam1 = Team{
	ID:              1,
	Name:            "MockTeam1",
	Abbreviation:    "T1",
	CreateDate:      "/Date(1536962113000-0500)/",
	ModifyDate:      "/Date(1593723490000-0500)/",
	NumericPriority: 12,
}

var mockTeam2 = Team{
	ID:              2,
	Name:            "MockTeam2",
	Abbreviation:    "T2",
	CreateDate:      "/Date(1536962113000-0500)/",
	ModifyDate:      "/Date(1593723490000-0500)/",
	NumericPriority: 12,
}

var mockTeamResponse = `
{
	"Next": "https://testing.tpondemand.com/api/v1/Teams/?take=25&skip=25",
	"Items": [
	  {
		"ResourceType": "Team",
		"Id": 1,
		"Name": "MockTeam1",
		"Description": null,
		"StartDate": null,
		"EndDate": null,
		"CreateDate": "/Date(1536962113000-0500)/",
		"ModifyDate": "/Date(1593723490000-0500)/",
		"LastCommentDate": null,
		"Tags": "",
		"NumericPriority": 12,
		"Icon": null,
		"EmojiIcon": null,
		"IsActive": true,
		"Abbreviation": "T1"
	  }
	]
  }  
`
var mockTeamResponse1 = `
{
	"Items": [
	  {
		"ResourceType": "Team",
		"Id": 1,
		"Name": "MockTeam1",
		"Description": null,
		"StartDate": null,
		"EndDate": null,
		"CreateDate": "/Date(1536962113000-0500)/",
		"ModifyDate": "/Date(1593723490000-0500)/",
		"LastCommentDate": null,
		"Tags": "",
		"NumericPriority": 12,
		"Icon": null,
		"EmojiIcon": null,
		"IsActive": true,
		"Abbreviation": "T1"
	  }
	]
  }  
`

var mockTeamResponse2 = `
{
	"Items": [
	  {
		"ResourceType": "Team",
		"Id": 2,
		"Name": "MockTeam2",
		"Description": null,
		"StartDate": null,
		"EndDate": null,
		"CreateDate": "/Date(1536962113000-0500)/",
		"ModifyDate": "/Date(1593723490000-0500)/",
		"LastCommentDate": null,
		"Tags": "",
		"NumericPriority": 12,
		"Icon": null,
		"EmojiIcon": null,
		"IsActive": true,
		"Abbreviation": "T2"
	  }
	]
  }  
`

func TestClient_GetTeam(t *testing.T) {
	client := NewFakeClient()

	tests := []struct {
		name    string
		team    string
		want    Team
		wantErr bool
	}{
		{
			name:    "team 1",
			team:    "MockTeam1",
			want:    mockTeam1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.GetTeam(tt.team)
			// Patch the client
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

func TestClient_GetTeams(t *testing.T) {
	client := NewFakeClient()
	tests := []struct {
		name    string
		filters []QueryFilter
		want    []Team
		wantErr bool
	}{
		{
			name:    "no filter",
			filters: nil,
			want: []Team{
				mockTeam1,
				mockTeam2,
			},
			wantErr: false,
		},
		{
			name: "filter",
			filters: []QueryFilter{
				Where("Name eq 'MockTeam1'"),
			},
			want: []Team{
				mockTeam1,
			},
			wantErr: false,
		},
		{
			name: "no results",
			filters: []QueryFilter{
				Where("Name eq 'dne'"),
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.GetTeams(tt.filters...)

			// Patch clients
			for i := range got {
				got[i].client = client
			}
			for i := range tt.want {
				tt.want[i].client = client
			}

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}
		})
	}
}
