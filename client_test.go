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
	"encoding/json"
	"fmt"
	"os"
)

// Example of client.Get() for godoc
func ExampleClient_Get() {
	tpClient, err := NewClient("exampleaccount", "superSecretToken")
	if err != nil {
		fmt.Println("Failed to create tp client:", err)
		os.Exit(1)
	}
	var response = UserResponse{}
	err = tpClient.Get(response,
		"User",
		nil)
	if err != nil {
		fmt.Println("Failed to get users:", err)
		os.Exit(1)
	}
	jsonBytes, _ := json.Marshal(response)
	fmt.Print(string(jsonBytes))
}
