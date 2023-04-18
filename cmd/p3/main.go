package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var Store *Storage

func main() {
	Store = new(Storage)

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

		resp := GeneralResponse{
			Status:  "ok",
			Time:    time.Now(),
			Message: fmt.Sprintf("Hello %s", req.GetFullName()),
		}

		c.JSON(200, resp)
	})

	group := r.Group("/v1")
	{
		group.POST("/rooms", AddNewRoom)
		group.GET("/rooms", GetRooms)

	}

	httpServer := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	httpServer.ListenAndServe()
}

func GetRooms(c *gin.Context) {
	c.JSON(http.StatusCreated, GeneralResponse{
		Status: "ok",
		Time:   time.Now(),
		Data:   Store.Rooms,
	})
}

func AddNewRoom(c *gin.Context) {
	var req CreateRoomRequest
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	Store.Rooms = append(Store.Rooms, Room{
		Name: req.RoomName,
	})
	c.JSON(http.StatusCreated, GeneralResponse{
		Status: "ok",
		Time:   time.Now(),
		Data:   Store.Rooms,
	})
}

type GeneralResponse struct {
	Status  string      `json:"status"`
	Time    time.Time   `json:"time"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type HelloRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type CreateRoomRequest struct {
	RoomName string `json:"room_name"`
}

func (hr HelloRequest) GetFullName() string {
	return hr.FirstName + " " + hr.LastName
}

type Storage struct {
	Rooms []Room
}

type Room struct {
	Name string
}
