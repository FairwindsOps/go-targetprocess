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
)

// DateTime currently does nothing special but will represent the TP DateTime objects
type DateTime string

// Assignments is a generic entity that lists assignments
type Assignments struct {
	Items []Assignment `json:",omitempty"`
}

// Assignment is a generic entity that lists a single assignment
type Assignment struct {
	ID           int32  `json:"Id,omitempty"`
	ResourceType string `json:",omitempty"`
	GeneralUser  *User  `json:",omitempty"`
}

// AssignedUser is used in UserStories and potentially other places. Returns a list of user assignments
type AssignedUser struct {
	Items []UserAssignment `json:",omitempty"`
}

// UserAssignment has its own unique Id and also includes a reference to a user, which also has an Id
type UserAssignment struct {
	User
}

// TeamAssignment has it's own unique Id and also includes a reference to the team, which also has an Id
type TeamAssignment struct {
	ID        int32    `json:"Id,omitempty"`
	StartDate DateTime `json:",omitempty"`
	EndDate   DateTime `json:",omitempty"`
	Team      *Team    `json:",omitempty"`
}

// GenerateURL takes an account name and entityID and returns a URL that should work in a browser
func GenerateURL(account string, entityID int32) string {
	return fmt.Sprintf("https://{%s}.tpondemand.com/entity/%d/RestUI/board.aspx", account, entityID)
}
