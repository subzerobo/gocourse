package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var Store *Storage

var RoomCreatedChannel chan string

var httpServer http.Server

func main() {

	// Create Context
	parentCtx := context.Background()
	ctx, cancel := context.WithCancel(parentCtx)

	// Listen for exit signals
	done := make(chan bool, 1)
	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)

	Store = new(Storage)
	RoomCreatedChannel = make(chan string)

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

	httpServer = http.Server{
		Addr:    ":9000",
		Handler: r,
	}

	// Run Graceful shutdown
	go GracefulShutdown(exitSignal, done)

	// Run Background Worker
	go BGW(ctx)

	log.Println("Running web server in blocking mode")
	err := httpServer.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

	log.Println("We are blocked in here waiting for done channel being close or getting some value")
	// Wait for HTTP Server to be killed gracefully !
	<-done

	// Calling cancel function to close the context
	cancel()

	time.Sleep(10 * time.Second)

	fmt.Println("Final Output before exiting the code")
}

// BG WORKER

func BGW(ctx context.Context) {
	for {
		select {
		case room := <-RoomCreatedChannel:
			log.Printf("New Room With Name: %s Created! \n", room)
		case <-ctx.Done():
			log.Println("We got the Context Done!")
			return
		}
	}
}

// METHODS

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

	RoomCreatedChannel <- req.RoomName

	Store.Rooms = append(Store.Rooms, Room{
		Name: req.RoomName,
	})
	c.JSON(http.StatusCreated, GeneralResponse{
		Status: "ok",
		Time:   time.Now(),
		Data:   Store.Rooms,
	})
}

func GracefulShutdown(exitSignal <-chan os.Signal, done chan<- bool) {
	log.Println("Graceful shutdown method started!")
	val := <-exitSignal
	log.Printf("We got signal: %s \n", val.String())

	// Create a 5s timeout context or waiting for app to shut down after 5 seconds
	ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelTimeout()

	httpServer.SetKeepAlivesEnabled(false)
	if err := httpServer.Shutdown(ctxTimeout); err != nil {
		log.Println(err)
	}
	log.Println("[OK] HTTPServer graceful shutdown completed")

	close(done)
}

// STRUCTS

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
