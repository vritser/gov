# Gov Web Framework

Gov is a web framework written in Golang.

## Getting Started

After installing Go and setting up your [GOPATH](http://golang.org/doc/code.html#GOPATH), create your first `.go` file. We'll call it `server.go`.

~~~ go
package main

import (
	"github.com/vritser/gov"
)

func main() {
	r := gov.New()
	r.Get("/", func(c *gov.Context) {
		c.String("Hello world!")
	})
	r.Run()
}
~~~

Then install the Gov package (**go 1.1** or greater is required):
~~~
go get github.com/vritser/gov
~~~

Then run your server:
~~~
go run server.go
~~~

You will now have a Gov webserver running on `localhost:9000`.


