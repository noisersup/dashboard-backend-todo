package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"kwiatek.xyz/todo-backend/database"
	"kwiatek.xyz/todo-backend/handlers"
)

func main(){
	if len(os.Args)<2 { log.Fatalf("You must specify url address to database!")}
	uri := os.Args[1] 
	
	db,err := database.ConnectToDatabase(uri,"tasks","tasks")
	if err != nil { log.Panic(err) }

	defer func(){
		if err = db.Disconnect(); err!=nil{
			log.Fatalf("Problem with disconnecting: %s",err.Error())
		}
	}()

	tasks,err := db.GetTasks()
	if err != nil { log.Panic(err) }

	for _,task := range tasks {
		log.Print(task)
	}

	h := handlers.CreateHandlers(db)

	r := mux.NewRouter()
	
	r.HandleFunc("/tasks", h.GetTasks).Methods("GET")

	http.ListenAndServe(":8000",r)
}