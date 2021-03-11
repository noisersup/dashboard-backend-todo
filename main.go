package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"kwiatek.xyz/todo-backend/database"
	han "kwiatek.xyz/todo-backend/handlers"
)

func main(){	
	corsPtr := flag.Bool("cors",false,"Enable CORS mode for locally debugging purposes.")
	uriPtr := flag.String("uri","","Specify uri do mongodb database.")
	flag.Parse()
	
	if *uriPtr=="" { log.Fatalf("You must specify url address to database!")}

	log.Printf("Connecting to database on %s",*uriPtr)
	db,err := database.ConnectToDatabase(*uriPtr,"tasks","tasks")
	if err != nil { log.Panic(err) }

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
	r.HandleFunc("/tasks", h.AddTask).Methods("POST")
	r.HandleFunc("/tasks/{id}", h.RemoveTask).Methods("DELETE")
	
	var httpHandler http.Handler 

	if *corsPtr {
		log.Printf("Warning!!! You are running CORS enabled mode. Do not use it on production")
		headersOk := handlers.AllowedHeaders([]string{"*","Content-Type"})
		originsOk := handlers.AllowedOrigins([]string{"*"})
		methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})
		httpHandler = handlers.CORS(originsOk, headersOk, methodsOk)(r)
	}else{
		httpHandler = r
	}

	log.Printf("Starting http server on port :8000...")
	log.Fatal(http.ListenAndServe(":8000",httpHandler))
}