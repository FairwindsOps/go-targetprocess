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
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func matchedURLValues(expected, got url.Values, key string) (bool, string, string) {
	return expected.Get(key) == got.Get(key), expected.Get(key), got.Get(key)
}

func ExampleFirst() {
	tpClient := NewClient("accountName", "superSecretToken")
	userStories, err := tpClient.GetUserStories(
		false,
		Where("EntityState.Name == 'Done'"),
		First(),
	)
	if err != nil {
		fmt.Println("Error getting UserStories:", err)
		os.Exit(1)
	}
	fmt.Printf("%+v\n", userStories)
}

func TestFirst(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
		want    url.Values
	}{
		{
			name:    "valid",
			wantErr: false,
			want:    url.Values{"take": []string{"1"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vals := url.Values{}
			first := First()
			got, err := first(vals)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				match, w, g := matchedURLValues(tt.want, got, "take")
				assert.True(t, match)
				assert.EqualValues(t, w, g)
			}
		})
	}
}

func ExampleMaxPerPage() {
	tpClient := NewClient("accountName", "superSecretToken")
	userStories, err := tpClient.GetUserStories(
		false,
		MaxPerPage(200),
	)
	if err != nil {
		fmt.Println("Error getting UserStories:", err)
		os.Exit(1)
	}
	fmt.Printf("%+v\n", userStories)
}

func TestMaxPerPage(t *testing.T) {
	tests := []struct {
		name    string
		count   int
		wantErr bool
		want    url.Values
	}{
		{
			name:    "valid",
			count:   100,
			wantErr: false,
			want:    url.Values{"take": []string{"100"}},
		},
		{
			name:    "negative conversion",
			count:   -100,
			wantErr: false,
			want:    url.Values{"take": []string{"100"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vals := url.Values{}
			mpp := MaxPerPage(tt.count)
			got, err := mpp(vals)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				match, w, g := matchedURLValues(tt.want, got, "take")
				assert.True(t, match)
				assert.EqualValues(t, w, g)
			}
		})
	}
}

func ExampleResult() {
	tpClient := NewClient("accountName", "superSecretToken")
	response := struct {
		EffortSum float64 `json:",omitempty"`
	}{}
	err := tpClient.Get(&response,
		"UserStories",
		nil,
		Result("effortSum:sum(effort)"),
	)
	if err != nil {
		fmt.Println("Error getting UserStories:", err)
		os.Exit(1)
	}
	fmt.Printf("%+v\n", response)
}

func TestResult(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		wantErr bool
		want    url.Values
	}{
		{
			name:    "valid",
			query:   "effortSum:sum(effort)",
			wantErr: false,
			want:    url.Values{"result": []string{"{effortSum:sum(effort)}"}},
		},
		{
			name:    "empty",
			query:   "",
			wantErr: false,
			want:    url.Values{"result": []string{"{}"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vals := url.Values{}
			result := Result(tt.query)
			got, err := result(vals)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				match, w, g := matchedURLValues(tt.want, got, "result")
				assert.True(t, match)
				assert.EqualValues(t, w, g)
			}
		})
	}
}

func ExampleSelect() {
	tpClient := NewClient("accountName", "superSecretToken")
	userStories, err := tpClient.GetUserStories(
		false,
		Select("name,id"),
	)
	if err != nil {
		fmt.Println("Error getting UserStories:", err)
		os.Exit(1)
	}
	fmt.Printf("%+v\n", userStories)
}

func TestSelect(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		wantErr bool
		want    url.Values
	}{
		{
			name:    "valid",
			query:   "id,name,assignedUser.Where(login=='jane@example.com'),responsibleTeam:{responsibleTeam.id,responsibleTeam.team},entityState",
			wantErr: false,
			want:    url.Values{"select": []string{"{id,name,assignedUser.Where(login=='jane@example.com'),responsibleTeam:{responsibleTeam.id,responsibleTeam.team},entityState}"}},
		},
		{
			name:    "empty",
			query:   "",
			wantErr: false,
			want:    url.Values{"select": []string{"{}"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vals := url.Values{}
			sel := Select(tt.query)
			got, err := sel(vals)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				match, w, g := matchedURLValues(tt.want, got, "select")
				assert.True(t, match)
				assert.EqualValues(t, w, g)
			}
		})
	}
}

func ExampleWhere() {
	tpClient := NewClient("accountName", "superSecretToken")
	userStories, err := tpClient.GetUserStories(
		false,
		Where("EntityState.Name == 'Done'"),
		Where("Team.Name == 'Team-1'"),
	)
	if err != nil {
		fmt.Println("Error getting UserStories:", err)
		os.Exit(1)
	}
	fmt.Printf("%+v\n", userStories)
}

func TestWhere(t *testing.T) {
	tests := []struct {
		name    string
		query   []string
		wantErr bool
		want    url.Values
	}{
		{
			name:    "valid",
			query:   []string{"EntityState.Name == 'Done'"},
			wantErr: false,
			want:    url.Values{"where": []string{"EntityState.Name == 'Done'"}},
		},
		{
			name:    "multiple",
			query:   []string{"EntityState.Name == 'Done'", "Team.Name == 'Administrators'"},
			wantErr: false,
			want:    url.Values{"where": []string{"EntityState.Name == 'Done' and Team.Name == 'Administrators'"}},
		},
		{
			name:    "empty",
			query:   []string{""},
			wantErr: false,
			want:    url.Values{"where": []string{""}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vals := url.Values{}
			where := Where(tt.query...)
			got, err := where(vals)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				match, w, g := matchedURLValues(tt.want, got, "where")
				assert.True(t, match)
				assert.EqualValues(t, w, g)
			}
		})
	}
}
