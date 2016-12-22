# go-masker [![CircleCI](https://circleci.com/gh/syoya/go-masker.svg?style=svg)](https://circleci.com/gh/syoya/go-masker)

Masks fields in JSON.

## Features

- Masks specified fields.
- Masks with specified string.
- Masks at deep position.

## Use Cases

- Logging a request body includes secret fields in API server.
- Logging a DB record includes secret fields.

## Installation

```
go get -u github.com/syoya/go-masker
```

## Usage

See [example](./example/main.go)
