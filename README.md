# Vera Luup Golang (In Progress)

Use Golang to remotely login and make Luup calls to a Vera™ home controller UI7

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/c0725b4e5c9144a0bb6e128444cd365a)](https://app.codacy.com/gh/scorpionknifes/go-vera-luup?utm_source=github.com&utm_medium=referral&utm_content=scorpionknifes/go-vera-luup&utm_campaign=Badge_Grade)
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
