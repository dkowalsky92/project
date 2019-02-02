package main

import (
	"fmt"
	"net/http"

	"github.com/dkowalsky/todo/db"
	"github.com/dkowalsky/todo/task"
	"github.com/go-chi/chi"
)

// Router - hub for networking
type Router struct {
	mux      *chi.Mux
	database *db.DB
}

// NewRouter - creates a new router
func NewRouter(db *db.DB) *Router {
	mux := chi.NewRouter()

	taskRouter := task.NewRouter(db)
	mux.Mount("/task", taskRouter.Mux)

	return &Router{database: db, mux: mux}
}

// Run - starts the server
func (r *Router) Run() {
	path := "localhost:5555"
	err := http.ListenAndServe(path, r.mux)
	panic(err)
}

func main() {
	fmt.Println("Connecting to database...")
	db, err := db.Connect()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected.")

	router := NewRouter(db)

	fmt.Println("Server is running.")
	router.Run()
}
