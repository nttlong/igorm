package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

type TodoItem struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	IsComplete bool   `json:"isComplete"`
}

var (
	todos  = []TodoItem{}
	nextID = 1
	mu     sync.Mutex
)

func createTodo(c echo.Context) error {
	var t TodoItem
	if err := c.Bind(&t); err != nil {
		fmt.Println("Bind error:", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}

	mu.Lock()
	t.ID = nextID
	nextID++
	todos = append(todos, t)
	mu.Unlock()

	return c.JSON(http.StatusCreated, t)
}

func main() {
	e := echo.New()

	e.POST("/api/todo", createTodo)

	e.Logger.Fatal(e.Start(":8080"))
}
