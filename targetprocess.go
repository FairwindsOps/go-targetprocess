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

// Package targetprocess is a go library to make using the Targetprocess API easier. Some
// public types are included to ease in json -> struct unmarshaling.
// A lot of inspiration for this package comes from https://github.com/adlio/trello
//
// Example Usage:
//  func main() {
//    logger := logrus.New()
//    tpClient := tp.NewClient("exampleCompany", "superSecretToken")
//    tpClient.Logger = logger
//    userStories, err := tpClient.GetUserStories(
//    	// The Where() filter function takes in any queries the targetprocess API accepts
//    	// Read about those here: https://dev.targetprocess.com/docs/sorting-and-filters
//    	tp.Where("EntityState.Name != 'Done'"),
//    	tp.Where("EntityState.Name != 'Backlog'"),
//    	// Simlar to Where(), the Include() function will limit the
//    	// response to a given list of fields
//    	tp.Include("Team", "Name", "ModifyDate"),
//    )
//    if err != nil {
//    	fmt.Println(err)
//    	os.Exit(1)
//    }
//    jsonBytes, _ := json.Marshal(userStories)
//    fmt.Print(string(jsonBytes))
//  }
//
// go-targetprocess includes some built-in structs that can be used for Users, Projects, Teams, and UserStories. You don't
// have to use those though and can use the generic `Get()` method with a custom struct as the output for a response to be
// JSON decoded into. Filtering functions (`Where()`, `Include()`, etc.) can be used in `Get()` just like they can in
// any of the helper functions.
//
// Example:
//  func main() {
//    out := struct {
//    	Next  string
//    	Prev  string
//    	Items []interface{}
//    }{}
//    tpClient := tp.NewClient("exampleCompany", "superSecretToken")
//    err := tpClient.Get(&out, "Users", nil)
//    if err != nil {
//    	fmt.Println(err)
//    	os.Exit(1)
//    }
//    jsonBytes, _ := json.Marshal(out)
//    fmt.Print(string(jsonBytes))
//  }
//
// Debug Logging:
// This idea was taken directly from the https://github.com/adlio/trello package. To add a debug logger,
// do the following:
//  logger := logrus.New()
//  // Also supports logrus.InfoLevel but that is default if you leave out the SetLevel method
//  logger.SetLevel(logrus.DebugLevel)
//  client := targetprocess.NewClient(accountName, token)
//  client.Logger = logger
//
package targetprocess
