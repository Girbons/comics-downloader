# Environment Setup

## Go

You must install [go](https://golang.org/doc/install).

## Installing dependencies

Dependencies are handled using [dep](https://github.com/golang/dep), [here](https://golang.github.io/dep/docs/installation.html) is the guide to install it.

Once installed run:

```
dep ensure
```

which will install the project's dependencies

## Linter

This project use [golangci-lint](https://github.com/golangci/golangci-lint), [here](https://github.com/golangci/golangci-lint#install) the guide to install it.

Once installed run:

```
golangci-lint run
```
