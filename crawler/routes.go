package main

import (
	"github.com/labstack/echo"
)

type (
	response struct {
		Message string
	}
)

func start(c *echo.Context) error {
	go connect(c.Query("id"))

	return c.JSON(200, &response{
		Message: "start",
	})
}

func stop(c *echo.Context) error {
	go disconnect(c.Query("id"))

	return c.JSON(200, &response{
		Message: "stop",
	})
}
