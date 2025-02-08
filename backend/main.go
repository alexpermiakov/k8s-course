package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type Todo struct {
	ID   int    `json:"id"`
	Task string `json:"task"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var todos []Todo

func main() {
	todos = []Todo{
		{ID: 1, Task: "Buy groceries"},
		{ID: 2, Task: "Walk the dog"},
	}

	r := mux.NewRouter()
	apiRouter := r.PathPrefix("/api/v1").Subrouter()

	apiRouter.HandleFunc("/todos", getTodos).Methods("GET")
	apiRouter.HandleFunc("/todos", createTodo).Methods("POST")
	apiRouter.HandleFunc("/todos/{id}", deleteTodo).Methods("DELETE")
	apiRouter.HandleFunc("/todos/{id}", updateTodo).Methods("PUT")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, todos)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if todo.Task == "" || len(todo.Task) > 200 {
		respondError(w, http.StatusBadRequest, "Task is required and must be less than 200 characters")
		return
	}

	newID := 1
	if len(todos) > 0 {
		newID = todos[len(todos)-1].ID + 1
	}
	todo.ID = newID
	todos = append(todos, todo)

	respondJSON(w, http.StatusCreated, todo)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	for i, todo := range todos {
		if todo.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	respondError(w, http.StatusNotFound, "Todo not found")
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if todo.Task == "" || len(todo.Task) > 200 {
		respondError(w, http.StatusBadRequest, "Task is required and must be less than 200 characters")
		return
	}

	for i, t := range todos {
		if t.ID == id {
			todos[i].Task = todo.Task
			respondJSON(w, http.StatusOK, nil)
			return
		}
	}

	respondError(w, http.StatusNotFound, "Todo not found")
}

func respondError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, ErrorResponse{Code: code, Message: message})
}

func respondJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}
