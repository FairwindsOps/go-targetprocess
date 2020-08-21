# Contributing

Issues, whether bugs, tasks, or feature requests are essential for improving this project. We believe it should be as easy as possible to contribute changes that get things working in your environment. There are a few guidelines that we need contributors to follow so that we can keep on top of things.

## Code of Conduct

This project adheres to a [code of conduct](CODE_OF_CONDUCT.md). Please review this document before contributing to this project.

## Sign the CLA

Before you can contribute, you will need to sign the [Contributor License Agreement](https://cla-assistant.io/fairwindsops/insights-plugins).

## Project Structure

Package targetprocess is a go library to make using the Targetprocess API easier.

## Getting Started

We label issues with the ["good first issue" tag](https://github.com/FairwindsOps/go-targetprocess/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22) if we believe they'll be a good starting point for new contributors. If you're interested in working on an issue, please start a conversation on that issue, and we can help answer any questions as they come up. Another good place to start would be adding regression tests to existing plugins.

## Setting Up Your Development Environment

### Prerequisites

* A fully functioning golang 1.14 environment

### Setup

When doing local development it helps to have a separate project on your machine to test new features or bugfixes. To accomplish this, in your copy of the project, set up the go.mod with a `replace` directive:
```
module github.com/example/foo

replace github.com/FairwindsOps/go-targetprocess => /Users/example/Projects/go-targetprocess

require (
	github.com/FairwindsOps/go-targetprocess v0.0.0
)
```
On the `replace` directive, modify the path after the `=>` to the location on your workstation where you made the copy of go-targetprocess.

## Running Tests

Running `make` will lint go files and run tests

## Creating a New Issue

If you've encountered an issue that is not already reported, please create a [new issue](https://github.com/FairwindsOps/go-targetprocess/issues), choose `Bug Report`, `Feature Request` or `Misc.` and follow the instructions in the template.


## Creating a Pull Request

Each new pull request should:

- Reference any related issues
- Pass existing tests and linting
- Contain a clear indication of if they're ready for review or a work in progress
- Be up to date and/or rebased on the master branch

## Creating a new release

Push a tag with the new version.
