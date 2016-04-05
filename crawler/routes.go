package main

import (
	"github.com/labstack/echo"
	"github.com/s-aska/Justaway-Ex/crawler/crawler"
)

type (
	response struct {
		Success bool
	}
)

func start(c *echo.Context) error {
	go crawler.Connect(c.Param("id"))

	return c.JSON(200, &response{
		Success: true,
	})
}

func stop(c *echo.Context) error {
	go crawler.Disconnect(c.Param("id"))

	return c.JSON(200, &response{
		Success: true,
	})
}
