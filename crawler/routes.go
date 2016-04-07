package main

import (
	"github.com/labstack/echo"
	"github.com/s-aska/Justaway-Ex/crawler/crawler"
)

type (
	responseConnect struct {
		Success bool
	}
	responseCount struct {
		Count int
	}
)

func start(c echo.Context) error {
	go crawler.Connect(c.Param("id"))

	return c.JSON(200, &responseConnect{
		Success: true,
	})
}

func stop(c echo.Context) error {
	go crawler.Disconnect(c.Param("id"))

	return c.JSON(200, &responseConnect{
		Success: true,
	})
}

func status(c echo.Context) error {
	count := crawler.Count()

	return c.JSON(200, &responseCount{
		Count: count,
	})
}

func startup(c echo.Context) error {
	crawler.ConnectAll()

	return c.JSON(200, &responseConnect{
		Success: true,
	})
}

func shutdown(c echo.Context) error {
	crawler.DisconnectAll()

	return c.JSON(200, &responseConnect{
		Success: true,
	})
}
