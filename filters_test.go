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
	"os"
)

func ExampleWhere() {
	tpClient := NewClient("accountName", "superSecretToken")
	userStories, err := tpClient.GetUserStories(
		Where("EntityState.Name eq 'Done'"),
		Where("Team.Name eq 'Team-1'"),
	)
	if err != nil {
		fmt.Println("Error getting UserStories:", err)
		os.Exit(1)
	}
	fmt.Printf("%+v\n", userStories)
}

func ExampleInclude() {
	tpClient := NewClient("accountName", "superSecretToken")
	userStories, err := tpClient.GetUserStories(
		Include("Name", "Description"),
	)
	if err != nil {
		fmt.Println("Error getting UserStories:", err)
		os.Exit(1)
	}
	fmt.Printf("%+v\n", userStories)
}

func ExampleMaxPerPage() {
	tpClient := NewClient("accountName", "superSecretToken")
	userStories, err := tpClient.GetUserStories(
		MaxPerPage(200),
	)
	if err != nil {
		fmt.Println("Error getting UserStories:", err)
		os.Exit(1)
	}
	fmt.Printf("%+v\n", userStories)
}

func ExampleFirst() {
	tpClient := NewClient("accountName", "superSecretToken")
	userStories, err := tpClient.GetUserStories(
		Where("EntityState.Name eq 'Done'"),
		First(),
	)
	if err != nil {
		fmt.Println("Error getting UserStories:", err)
		os.Exit(1)
	}
	fmt.Printf("%+v\n", userStories)
}
