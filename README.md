[![GitHub Release](https://img.shields.io/github/v/release/astronetes/go-sdk)](https://github.com/astronetes/go-sdk/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/astronetes/go-sdk.svg)](https://pkg.go.dev/github.com/astronetes/go-sdk)
[![go.mod](https://img.shields.io/github/go-mod/go-version/astronetes/go-sdk)](go.mod)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://img.shields.io/github/license/astronetes/go-sdk)
[![Build Status](https://img.shields.io/github/actions/workflow/status/astronetes/go-sdk/build.yml?branch=main)](https://github.com/astronetes/go-sdk/actions?query=workflow%3ABuild+branch%3Amain)
[![CodeQL](https://github.com/astronetes/go-sdk/actions/workflows/codeql.yml/badge.svg?branch=main)](https://github.com/astronetes/go-sdk/actions/workflows/codeql.yml)

# Go SDK for Astronetes

This module contains a collection of Go utilities that could be mainly used by Astronetes developments. The main goal
of this module is to provide a set of handful and reusable API's to be used when building Kubernetes operators.

## History and proeject status

This module is still `in active development` and the API is still subject to breaking changes.

Most of the provided functionalities are used by the operators that belongs to the Astronetes ecosystem.

## Installation

Use go get to retrieve the SDK to add it to your GOPATH workspace, or project's Go module dependencies.

```bash
go get -u github.com/astronetes/go-sdk
```

To update the SDK use go get -u to retrieve the latest version of the SDK.

```bash
go get -u github.com/astronetes/go-sdk
```

You could specify a concrete version of this module as It's shown on the below. Replace x.y.z by the desired version.

```bash
module github.com/<org>/<repository>
require ( 
  github.com/astronetes/go-sdk vX.Y.Z
)
```

## Overview of SDK's Packages

The SDK is composed of @TODO components,

* `astronetes`: SDK Core, provides common shared types such as Config, Logger, and utilities to make working with API 
  parameters easier.
* `cloud`: SDK cloud, provides API to interact with Cloud service providers
* `k8s`: Set of interfaces to interact with k8 and other belonging tools to it.  

## Getting started

### Examples

A rich and growing set of examples of usage of this module can be found in folder `eamples`.


### Contributing

See the [contributing](https://github.com/astronetes/go-sdk/blob/main/CONTRIBUTING.md) documentation.


