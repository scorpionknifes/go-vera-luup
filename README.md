# Vera Luup Golang (In Progress)

Use Golang to remotely login and make Luup calls to a Vera™ home controller UI7

[![Go Report Card](https://goreportcard.com/badge/github.com/scorpionknifes/go-vera-luup)](https://goreportcard.com/report/github.com/scorpionknifes/go-vera-luup) [![GoDoc](https://godoc.org/github.com/gogolfing/cbus?status.svg)](https://godoc.org/github.com/scorpionknifes/go-vera-luup)

## Features

-   Remote access using Vera account
-   Polling for device updates
-   Switching Power Status
-   Switching Lock Door Status

## How to use

`go get github.com/scorpionknifes/go-vera-luup`

```go
import vera "github.com/scorpionknifes/go-vera-luup"
```

-   examples shown in examples/main.go
-   config variables in .env
-   use `Controller` and `Vera` struct

## TODO

-   Unit tests
-   Luup live_energy_usage
-   Luup variableset

## Luup Information

-   [Luup Introduction](http://wiki.micasaverde.com/index.php/Luup_Intro)
-   [Luup Requests](http://wiki.micasaverde.com/index.php/Luup_Requests)
