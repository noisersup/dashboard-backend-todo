package handlers

import (
	"encoding/json"
	"net/http"

	"kwiatek.xyz/todo-backend/database"
	"kwiatek.xyz/todo-backend/models"
)
type TodoServer struct{
	db 	*database.Database
}
func CreateHandlers(db *database.Database) TodoServer{
	return TodoServer{db}
}

func (todo *TodoServer)GetTasks(w http.ResponseWriter, r *http.Request){
	tasks := models.Tasks{}
	t ,err := todo.db.GetTasks()

	if err !=nil {}

	for _,pointer := range t {
		tasks.Tasks= append(tasks.Tasks, *pointer)
	}

	// w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}