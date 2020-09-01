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

// User matches up with a targetprocess User
type User struct {
	CustomFields    []CustomField `json:",omitempty"`
	CreateDate      DateTime      `json:",omitempty"`
	ModifyDate      DateTime      `json:",omitempty"`
	DeleteDate      DateTime      `json:",omitempty"`
	Email           string        `json:",omitempty"`
	FirstName       string        `json:",omitempty"`
	GlobalID        string        `json:",omitempty"`
	ID              int32         `json:"Id"`
	IsActive        bool          `json:",omitempty"`
	IsAdministrator bool          `json:",omitempty"`
	LastName        string        `json:",omitempty"`
	Locale          string        `json:",omitempty"`
	Login           string        `json:",omitempty"`
}

// UserResponse is a representation of the http response for a group of Users
type UserResponse struct {
	Items []User
	Next  string
	Prev  string
}

// GetUsers will return all users
func (c *Client) GetUsers(filters ...QueryFilter) ([]User, error) {
	var ret []User
	out := UserResponse{}

	err := c.Get(&out, "User", nil, filters...)
	if err != nil {
		return nil, err
	}
	ret = append(ret, out.Items...)
	for out.Next != "" {
		innerOut := UserResponse{}
		err := c.GetNext(&innerOut, out.Next)
		if err != nil {
			return ret, err
		}
		ret = append(ret, innerOut.Items...)
		out = innerOut
	}
	return ret, nil
}
