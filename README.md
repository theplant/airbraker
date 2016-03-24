# Airbraker

[![Build Status](https://semaphoreci.com/api/v1/projects/3cff44c2-9a48-4efe-a768-2d469a0b9074/744940/badge.svg)](https://semaphoreci.com/theplant/airbraker)

Generic airbraker that integrated with [airbrake](https://airbrake.io/) service for [Gin-backed](https://gin-gonic.github.io/gin/) web application

# Quick Start

## Airbrake configuration

Set up environment variables.

```sh
$ export AIRBRAKE_PROJECT_ID="your-project-id"
$ export AIRBRAKE_TOKEN="your-token"
$ export AIRBRAKE_ENV="your-app-env"
```

## Set up Gin middleware

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/theplant/airbraker"
)

func main() {
    r := gin.Default()

    // Set up recover middleware
    r.Use(airbraker.Recover())

    r.GET("/panic", func(c *gin.Context) {
        panic(errors.New("unexpected error"))
    })
    r.Run() // listen and server on 0.0.0.0:8080
}
```

After you start with `go run main.go`. You'll see:

```sh
Logging errors to Airbrake '<AIRBRAKE_ENV>' env on project <AIRBRAKE_PROJECT_ID>
```

Then it works. :)
