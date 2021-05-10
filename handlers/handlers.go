package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"kwiatek.xyz/todo-backend/database"
	"kwiatek.xyz/todo-backend/models"
	"kwiatek.xyz/todo-backend/utils"
)
type TodoServer struct{
	db 	*database.Database
}
func CreateHandlers(db *database.Database) TodoServer{
	return TodoServer{db}
}

func (todo *TodoServer) GetTasks(w http.ResponseWriter, r *http.Request){
	log.Print("GET!") //TODO: remove

	response := models.GetResponse{}
	t ,err := todo.db.GetTasks()

	if err != nil { // Database problems [500 code]
		log.Printf("Database error: %s",err) //TODO: Log file

		response.Error = "Database internal error"
		utils.SendResponse(w,response,http.StatusInternalServerError)
		return
	}

	for _,pointer := range t {
		response.Tasks = append(response.Tasks, *pointer)
	}

	utils.SendResponse(w,response,http.StatusOK)
}

func (todo *TodoServer) GetDues(w http.ResponseWriter, r *http.Request){
	log.Print("GET Dues!") //TODO: remove

	response := models.GetResponse{}
	
	t,err := todo.db.GetDueTasks()
	if err != nil { // Database problems [500 code]
		log.Printf("Database error: %s",err) //TODO: Log file

		response.Error = "Database internal error"
		utils.SendResponse(w,response,http.StatusInternalServerError)
		return
	}
	
	response.Tasks = t
	utils.SendResponse(w,response,http.StatusOK)
}

func (todo *TodoServer) RemoveTask(w http.ResponseWriter, r *http.Request) {
	log.Print("DELETE!") //TODO: remove
	params := mux.Vars(r)
	
	response := models.ErrorResponse{}

	id,err := primitive.ObjectIDFromHex(params["id"])
	if err != nil { // ID parsing problems [400 code]
		log.Printf("ID parse error: %s",err) //TODO: Log file

		response.Error = "Cannot parse provided id"
		utils.SendResponse(w,response,http.StatusBadRequest)
		return
	}
	
	resp, err := todo.db.RemoveTask(id)
	if err != nil { // Database problems [500 code]
		log.Printf("Database error: %s",err) //TODO: Log file

		response.Error = "Database internal error"
		utils.SendResponse(w,response,http.StatusInternalServerError)
		return
	}
	if(resp.DeletedCount<1){ // Object not found [404 code]
		response.Error = "Object with provided ID not found"
		utils.SendResponse(w,response,http.StatusNotFound)
		return
	}

	utils.SendResponse(w,response,http.StatusOK)
}

func (todo *TodoServer) AddTask(w http.ResponseWriter, r *http.Request){
	log.Print("POST!") //TODO: remove
	var task models.NewTask
	response := models.ErrorResponse{}

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil { // JSON decoding problems [400 code]
		log.Printf("JSON decoding error: %s",err)

		response.Error = "Cannot parse json to task object."
		utils.SendResponse(w,response,http.StatusBadRequest)
		return
	}

	if task.Title == ""{
		response.Error = "No title provided"
		utils.SendResponse(w,response,http.StatusBadRequest)
		return
	}
	
	_, err = todo.db.AddTask(&task) 

	if err != nil { // Database problems [500 code]
		log.Printf("Database error: %s",err) //TODO: Log file

		response.Error = "Database internal error"
		utils.SendResponse(w,response,http.StatusInternalServerError)
		return
	}
	utils.SendResponse(w,response,http.StatusOK)
}

func (todo *TodoServer) DoneTask(w http.ResponseWriter, r *http.Request) {
	log.Print("PATCH!") //TODO: remove
	params := mux.Vars(r)

	response := models.ErrorResponse{}

	id,err := primitive.ObjectIDFromHex(params["id"])
	if err != nil { // ID parsing problems [400 code]
		log.Printf("ID parse error: %s",err) //TODO: Log file

		response.Error = "Cannot parse provided id"
		utils.SendResponse(w,response,http.StatusBadRequest)
		return
	}

	var doneModel models.DoneTask

	err = json.NewDecoder(r.Body).Decode(&doneModel)
	if err != nil { // JSON decoding problems [400 code]
		log.Printf("JSON decoding error: %s",err)//TODO: Log file

		response.Error = "Cannot parse json to task object."
		utils.SendResponse(w,response,http.StatusBadRequest)
		return
	}
	
	if &doneModel.Done == nil { // Done variable is not present [400 code]
		log.Printf("Done variable is empty.")//TODO: Log file 

		response.Error = "Done variable is empty"
		utils.SendResponse(w,response,http.StatusBadRequest)
		return
	}

	var res *mongo.UpdateResult

	if doneModel.Done {
		res,err = todo.db.DoneTask(id)
	}else{
		res,err = todo.db.UndoneTask(id)
	}

	if err != nil { // Database problems [500 code]
		log.Printf("Database error: %s",err) //TODO: Log file

		response.Error = "Database internal error"
		utils.SendResponse(w,response,http.StatusInternalServerError)
		return
	}

	if res.MatchedCount < 1 { // Object not found [404 code]
		response.Error = "Object with provided ID not found"
		utils.SendResponse(w,response,http.StatusNotFound)
		return
	}

	utils.SendResponse(w,response,http.StatusOK)
}