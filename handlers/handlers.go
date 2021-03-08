package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"kwiatek.xyz/todo-backend/database"
	"kwiatek.xyz/todo-backend/models"
)
type TodoServer struct{
	db 	*database.Database
}
func CreateHandlers(db *database.Database) TodoServer{
	return TodoServer{db}
}

func (todo *TodoServer) GetTasks(w http.ResponseWriter, r *http.Request){
	tasks := models.Tasks{}
	t ,err := todo.db.GetTasks()

	if err !=nil {}

	for _,pointer := range t {
		tasks.Tasks= append(tasks.Tasks, *pointer)
	}
	log.Print("GET!")
	// w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (todo *TodoServer) RemoveTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	log.Print("DELETE!")
	
	id,err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		//TODO
	}
	
	resp, err := todo.db.RemoveTask(id)
	if err != nil {
		log.Printf("DELETE DB ERROR: %s",err.Error())	
	}
	if(resp.DeletedCount<1){
		w.WriteHeader(http.StatusNoContent)
		return
	}

	json.NewEncoder(w).Encode(`{'ID':`+params["id"]+` }`)
}

func (todo *TodoServer) AddTask(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")
	var task models.NewTask

	log.Print("POST!")
	

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		log.Println(err)
		//TODO
	}

	res, err := todo.db.AddTask(task.Title,task.Desc) 
	log.Println(res.InsertedID)
	if err != nil {
		log.Println(err)
		//TODO
	}
	json.NewEncoder(w).Encode(&task)
}