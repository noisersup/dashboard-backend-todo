package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"kwiatek.xyz/todo-backend/database"
	han "kwiatek.xyz/todo-backend/handlers"
)

type DbConfig struct{
	Address		string 
	Port		int    
}

func main(){	
	corsPtr := flag.Bool("cors",false,"Enable CORS mode for locally debugging purposes.")
	flag.Parse()
	
	config := getVars() 

	log.Printf("Connecting to database on %s:%d",config.Address,config.Port)
	db,err := database.ConnectToDatabase(config.Address+":"+fmt.Sprint(config.Port),"tasks","tasks")
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

	h := han.CreateHandlers(db)

	r := mux.NewRouter()

	r.HandleFunc("/tasks", h.GetTasks).Methods("GET")
	r.HandleFunc("/tasks/due", h.GetDues).Methods("GET")
	
	r.HandleFunc("/tasks", h.AddTask).Methods("POST")
	r.HandleFunc("/tasks/{id}", h.RemoveTask).Methods("DELETE")
	r.HandleFunc("/tasks/{id}", h.DoneTask).Methods("PATCH")
	
	var httpHandler http.Handler 

	if *corsPtr {
		log.Printf("Warning!!! You are running CORS enabled mode. Do not use it on production")
		headersOk := handlers.AllowedHeaders([]string{"*","Content-Type"})
		originsOk := handlers.AllowedOrigins([]string{"*"})
		methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"})
		httpHandler = handlers.CORS(originsOk, headersOk, methodsOk)(r)
	}else{
		httpHandler = r
	}

	log.Printf("Starting http server on port :8000...")
	log.Fatal(http.ListenAndServe(":8000",httpHandler))
}

func getVars() *DbConfig {	
	var config DbConfig

	config.Address = os.Getenv("DB_ADDRESS")
	config.Port,_ = strconv.Atoi(os.Getenv("DB_PORT")) //default: 27017

	if config.Address=="" { 
		log.Fatal("ENV variables did not set")
	}

	config.Address = "mongodb://"+config.Address

	if config.Port == 0 {
		log.Print("Port invalid or not provided. Setting default (27017)")
		config.Port=27017
	}
	return &config
}