<div align="center">

  [![PkgGoDev][godoc-image]][godoc-link][![GitHub release (latest SemVer)][release-image]][release-link] [![GitHub go.mod Go version][version-image]][version-link] [![CircleCI][circleci-image]][circleci-link] [![Code Coverage][codecov-image]][codecov-link] [![Go Report Card][goreport-image]][goreport-link] [![Apache 2.0 license](https://img.shields.io/badge/licence-Apache2-brightgreen.svg)](https://opensource.org/licenses/Apache-2.0)

</div>



[godoc-image]: https://pkg.go.dev/badge/github.com/FairwindsOps/go-targetprocess
[godoc-link]: https://pkg.go.dev/github.com/FairwindsOps/go-targetprocess

[version-image]: https://img.shields.io/github/go-mod/go-version/FairwindsOps/go-targetprocess
[version-link]: https://github.com/FairwindsOps/go-targetprocess

[release-image]: https://img.shields.io/github/v/release/FairwindsOps/go-targetprocess
[release-link]: https://github.com/FairwindsOps/go-targetprocess

[goreport-image]: https://goreportcard.com/badge/github.com/FairwindsOps/go-targetprocess
[goreport-link]: https://goreportcard.com/report/github.com/FairwindsOps/go-targetprocess

[circleci-image]: https://circleci.com/gh/FairwindsOps/go-targetprocess/tree/master.svg?style=svg
[circleci-link]: https://circleci.com/gh/FairwindsOps/go-targetprocess

[codecov-image]: https://codecov.io/gh/FairwindsOps/go-targetprocess/branch/master/graph/badge.svg
[codecov-link]: https://codecov.io/gh/FairwindsOps/go-targetprocess

# Go-Targetprocess

Package targetprocess is a go library to make using the Targetprocess API easier. Some
public types are included to ease in json -> struct unmarshaling.
A lot of inspiration for this package comes from https://github.com/adlio/trello

## Join the Fairwinds Open Source Community

The goal of the Fairwinds Community is to exchange ideas, influence the open source roadmap, and network with fellow Kubernetes users. [Chat with us on Slack](https://join.slack.com/t/fairwindscommunity/shared_invite/zt-e3c6vj4l-3lIH6dvKqzWII5fSSFDi1g) or [join the user group](https://www.fairwinds.com/open-source-software-user-group) to get involved!

## Usage

```go
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	tp "github.com/FairwindsOps/go-targetprocess"
)

func main() {
	logger := logrus.New()
	tpClient, err := tp.NewClient("exampleCompany", "superSecretToken")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	tpClient.Logger = logger

	userStories, err := tpClient.GetUserStories(
		// we set paging to false so we only get the first page of results
		false,
		// The Where() filter function takes in any queries the targetprocess API accepts
		// Read about those here: https://dev.targetprocess.com/docs/sorting-and-filters
		tp.Where("EntityState.Name != 'Done'"),
		tp.Where("EntityState.Name != 'Backlog'"),
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
	tpClient, err := tp.NewClient("exampleCompany", "superSecretToken")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
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
client, err := targetprocess.NewClient(accountName, token)
if err != nil {
    fmt.Println(err)
    os.Exit(1)
}
client.Logger = logger
```

## Contributing

PRs welcome! Check out the [Contributing Guidelines](CONTRIBUTING.md) and
[Code of Conduct](CODE_OF_CONDUCT.md) for more information.


## Other Projects from Fairwinds

Enjoying go-targetprocess? Check out some of our other projects:
* [Polaris](https://github.com/FairwindsOps/Polaris) - Audit, enforce, and build policies for Kubernetes resources, including over 20 built-in checks for best practices
* [Goldilocks](https://github.com/FairwindsOps/Goldilocks) - Right-size your Kubernetes Deployments by compare your memory and CPU settings against actual usage
* [Pluto](https://github.com/FairwindsOps/Pluto) - Detect Kubernetes resources that have been deprecated or removed in future versions
* [Nova](https://github.com/FairwindsOps/Nova) - Check to see if any of your Helm charts have updates available
* [rbac-manager](https://github.com/FairwindsOps/rbac-manager) - Simplify the management of RBAC in your Kubernetes clusters

Or [check out the full list](https://www.fairwinds.com/open-source-software?utm_source=go-targetprocess&utm_medium=go-targetprocess&utm_campaign=go-targetprocess)
