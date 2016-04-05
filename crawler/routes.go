package main

import (
	"github.com/labstack/echo"
	"./crawler"
)

type (
	response struct {
		Success bool
	}
)

func start(c *echo.Context) error {
	go crawler.Connect(c.Query("id"))

	return c.JSON(200, &response{
		Success: true,
	})
}

func stop(c *echo.Context) error {
	go crawler.Disconnect(c.Query("id"))

	return c.JSON(200, &response{
		Success: true,
	})
}
