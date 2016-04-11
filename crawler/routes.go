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

type Router struct {
	crawler *crawler.Crawler
}

func (r *Router) start(c echo.Context) error {
	go r.crawler.Connect(c.Param("id"))

	return c.JSON(200, &responseConnect{
		Success: true,
	})
}

func (r *Router) stop(c echo.Context) error {
	go r.crawler.Disconnect(c.Param("id"))

	return c.JSON(200, &responseConnect{
		Success: true,
	})
}

func (r *Router) status(c echo.Context) error {
	count := r.crawler.Count()

	return c.JSON(200, &responseCount{
		Count: count,
	})
}

func (r *Router) startup(c echo.Context) error {
	r.crawler.ConnectAll()

	return c.JSON(200, &responseConnect{
		Success: true,
	})
}

func (r *Router) shutdown(c echo.Context) error {
	r.crawler.DisconnectAll()

	return c.JSON(200, &responseConnect{
		Success: true,
	})
}
