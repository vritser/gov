package main

import (
	"fmt"

	gov "gov/lib"
)

func main() {
	r := gov.New()

	r.Use(func(c *gov.Context) {
		fmt.Println(c.Method(), c.Path())

		c.Set("usr", map[string]interface{}{
			"id":   10,
			"name": "vritser",
			"age":  30,
		})
	})

	r.Get("/usr", func(c *gov.Context) {
		c.Json(c.GetStringMap("usr"))
	})

	r.Get("/test", func(c *gov.Context) {
		fmt.Fprint(c.Response, "hellow from router")
	})

	r.Get("/test", func(c *gov.Context) {
		c.String("Welcome to gov web framework!")
	})

	r.Get("/", func(c *gov.Context) {
		c.Status(200)
	})

	r.Post("/", func(c *gov.Context) {
		fmt.Fprint(c.Response, "post result")
	})

	r.Get("/list", func(c *gov.Context) {
		xs := []int{1, 2, 3, 4, 5}
		c.Json(xs)
	})

	r.Get("/book/:id", func(c *gov.Context) {
		c.Json(map[string]interface{}{
			"id": c.Param("id"),
		})
	})

	r.Get("/book/:id/test", func(c *gov.Context) {
		c.Json(map[string]interface{}{
			"id":   c.Param("id"),
			"path": interface{}(c.Path()),
		})
	})

	for _, rt := range r.Routes() {
		fmt.Println(rt.Method, rt.Path)
	}

	r.Run(":9000")
}
