package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func main() {

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(gin.Recovery())

	r.GET("/", func(c *gin.Context) {
		c.String(200, "This is my website!")
	})

	r.GET("/hello/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(200, "Hello %s", name)
	})

	r.POST("/jhello", func(c *gin.Context) {
		var req HelloRequest
		err := c.ShouldBind(&req)
		if err != nil {
			c.JSON(400, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		resp := HelloResponse{
			Status:  "ok",
			Time:    time.Now(),
			Message: fmt.Sprintf("Hello %s", req.GetFullName()),
		}

		c.JSON(200, resp)
	})

	httpServer := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	httpServer.ListenAndServe()
}

type HelloResponse struct {
	Status  string    `json:"status"`
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
}

type HelloRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (hr HelloRequest) GetFullName() string {
	return hr.FirstName + " " + hr.LastName
}
