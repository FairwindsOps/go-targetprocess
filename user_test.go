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

var mockUser1 = User{
	ID:        1,
	FirstName: "Mock",
	LastName:  "User",
	Email:     "mockuser@example.com",
	Login:     "mockuser",
	IsActive:  true,
}

var mockUser2 = User{
	ID:        2,
	FirstName: "Mock",
	LastName:  "User",
	Email:     "mockuser2@example.com",
	Login:     "mockuser2",
	IsActive:  true,
}

var mockUserResponse = `
{
	"Next": "https://testing.tpondemand.com/api/v1/Users/?format=json&take=25&skip=25",
	"Items": [
	  {
		"ResourceType": "User",
		"Kind": "User",
		"Id": 1,
		"FirstName": "Mock",
		"LastName": "User",
		"Email": "mockuser@example.com",
		"Login": "mockuser",
		"DeleteDate": null,
		"IsActive": true,
		"IsAdministrator": false,
		"Locale": null,
		"WeeklyAvailableHours": 40,
		"CurrentAllocation": 100,
		"CurrentAvailableHours": 0,
		"AvailableFutureAllocation": 100,
		"AvailableFutureHours": 40,
		"IsObserver": false,
		"IsContributor": false,
		"Skills": null,
		"ActiveDirectoryName": null,
		"RichEditor": "Markdown"
	  }
	]
}
`

var mockUserResponse1 = `
{
	"Items": [
	  {
		"ResourceType": "User",
		"Kind": "User",
		"Id": 1,
		"FirstName": "Mock",
		"LastName": "User",
		"Email": "mockuser@example.com",
		"Login": "mockuser",
		"DeleteDate": null,
		"IsActive": true,
		"IsAdministrator": false,
		"Locale": null,
		"WeeklyAvailableHours": 40,
		"CurrentAllocation": 100,
		"CurrentAvailableHours": 0,
		"AvailableFutureAllocation": 100,
		"AvailableFutureHours": 40,
		"IsObserver": false,
		"IsContributor": false,
		"Skills": null,
		"ActiveDirectoryName": null,
		"RichEditor": "Markdown"
	  }
	]
}
`

var mockUserResponse2 = `
{
	"Items": [
	  {
		"ResourceType": "User",
		"Kind": "User",
		"Id": 2,
		"FirstName": "Mock",
		"LastName": "User",
		"Email": "mockuser2@example.com",
		"Login": "mockuser2",
		"DeleteDate": null,
		"IsActive": true,
		"IsAdministrator": false,
		"Locale": null,
		"WeeklyAvailableHours": 40,
		"CurrentAllocation": 100,
		"CurrentAvailableHours": 0,
		"AvailableFutureAllocation": 100,
		"AvailableFutureHours": 40,
		"IsObserver": false,
		"IsContributor": false,
		"Skills": null,
		"ActiveDirectoryName": null,
		"RichEditor": "Markdown"
	  }
	]
}
`

func TestClient_GetUsers(t *testing.T) {
	client := NewFakeClient()

	tests := []struct {
		name    string
		filters []QueryFilter
		want    []User
		wantErr bool
	}{
		{
			name:    "no filter",
			filters: nil,
			want: []User{
				mockUser1,
				mockUser2,
			},
			wantErr: false,
		},
		{
			name: "filter 1",
			filters: []QueryFilter{
				Where("Id eq 1"),
			},
			want: []User{
				mockUser1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.GetUsers(tt.filters...)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}
		})
	}
}
