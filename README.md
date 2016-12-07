# go-masker [![CircleCI](https://circleci.com/gh/syoya/go-masker.svg?style=svg)](https://circleci.com/gh/syoya/go-masker)

Mask specific fields in JSON.

## Installation

```
go get -u github.com/syoya/go-masker
```

## Usage

```go
m, _ := masker.New(map[string]string{
  "password": "**********",
})
fmt.Printf("%s\n", m.Mask([]byte(`{"email":"foo@example.com","password":"p@ssw0rd"}`)))
    // -> {"email":"foo@example.com","password":"**********"}
```
