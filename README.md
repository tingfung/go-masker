# go-masker

Masks fields in JSON.

## Features

- Masks specified fields(case insensitive).
- Masks with specified string.
- Masks at deep position.

## Use Cases

- Logging a request body includes secret fields in API server.
- Logging a DB record includes secret fields.

## Installation

```
go get -u github.com/tingfung/go-masker
```

## Usage

See [example](./example/main.go)

@forked from [syoya/go-masker](https://github.com/syoya/go-masker) with tiny modify

