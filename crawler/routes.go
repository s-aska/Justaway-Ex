package main

import (
	"github.com/gin-gonic/gin"
)

func start(c *gin.Context) {
	go connect(c.Param("id"))

	c.JSON(200, gin.H{
		"message": "start streaming",
	})
}

func stop(c *gin.Context) {
	go disconnect(c.Param("id"))

	c.JSON(200, gin.H{
		"message": "stop streaming",
	})
}
