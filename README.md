# Furdarius\subswitch
[![Build Status](https://travis-ci.org/Furdarius/subswitch.svg?branch=master)](https://travis-ci.org/Furdarius/subswitch) [![Coverage Status](https://coveralls.io/repos/github/Furdarius/subswitch/badge.svg?branch=master)](https://coveralls.io/github/Furdarius/subswitch?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/Furdarius/subswitch)](https://goreportcard.com/report/github.com/Furdarius/subswitch)


Want to serve static from `static.domain.io` without nginx or apache? You wanna separate backend api on `api.domain.io`? Just do it with this library!

It works like simple middleware, so you can inject it in your app incredibly fast.

## Installation

With a properly configured Go toolchain:
```sh
go get github.com/furdarius/subswitch
```

## Usage

```go
// .. router init ..

// Let's serve static from domain.io
// and REST api from api.domain.io
ss := subswitch.New(http.FileServer(http.Dir("/tmp")), map[string]http.Handler{
    "api": router,
})

// Use it like http.Handler
http.ListenAndServe(":8080", ss)
```
