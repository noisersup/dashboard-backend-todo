package main

import (
	"log"
	"os"

	"kwiatek.xyz/todo-backend/database"
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
}