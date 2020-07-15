# Go-Targetprocess

Package targetprocess is a go library to make using the Targetprocess API easier. Some
public types are included to ease in json -> struct unmarshaling.
A lot of inspiration for this package comes from https://github.com/adlio/trello

## Usage

```go
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	tp "github.com/fairwindsops/go-targetprocess"
)

func main() {
	logger := logrus.New()
	tpClient := tp.NewClient("exampleCompany", "superSecretToken")
	tpClient.Logger = logger

	userStories, err := tpClient.GetUserStories(
		// The Where() filter function takes in any queries the targetprocess API accepts
        // Read about those here: https://dev.targetprocess.com/docs/sorting-and-filters
		tp.Where("EntityState.Name ne 'Done'"),
		tp.Where("EntityState.Name ne 'Backlog'"),
        // Simlar to Where(), the Include() function will limit the
        // response to a given list of fields
		tp.Include("Team", "Name", "ModifyDate"),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	jsonBytes, _ := json.Marshal(userStories)
	fmt.Print(string(jsonBytes))
}
```

## Custom structs for queries

go-targetprocess includes some built-in structs that can be used for Users, Projects, Teams, and UserStories. You don't
have to use those though and can use the generic `Get()` method with a custom struct as the output for a response to be
JSON decoded into. Filtering functions (`Where()`, `Include()`, etc.) can be used in `Get()` just like they can in
any of the helper functions.

Ex:

```go
func main() {
	out := struct {
		Next  string
		Prev  string
		Items []interface{}
	}{}
	tpClient := tp.NewClient("exampleCompany", "superSecretToken")
	err := tpClient.Get(&out, "Users", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	jsonBytes, _ := json.Marshal(out)
	fmt.Print(string(jsonBytes))
}
```

## Debug Logging

This idea was taken directly from the https://github.com/adlio/trello package. To add a debug logger,
do the following:

```go
logger := logrus.New()
// Also supports logrus.InfoLevel but that is default if you leave out the SetLevel method
logger.SetLevel(logrus.DebugLevel)
client := targetprocess.NewClient(accountName, token)
client.Logger = logger
```
