package main

import (
	"fmt"
	"log"
	"testing"

	"kwiatek.xyz/todo-backend/database"
	"kwiatek.xyz/todo-backend/models"
)

type TestVars struct{
	Db	*database.Database
}

var test TestVars

func TestMain(m *testing.M) {
	config := getVars() 

	var err error
	log.Printf("Connecting to database...")
	test.Db,err = database.ConnectToDatabase(config.Address+":"+fmt.Sprint(config.Port),"tasks_test","tasks")
	if err != nil {log.Fatal(err)}
	
	log.Printf("Getting tasks...")
	t,err := test.Db.GetTasks()
	if err != nil {log.Fatal(err)}

	for _,pointer := range t {
		_,err = test.Db.RemoveTask(pointer.ID)
		if err != nil {log.Fatal(err)}
	}
	
	m.Run()
}

func TestGetDueTasks(t *testing.T){
	expectedAmount := 3
	for i := 0; i < 5; i++ {
		_,err := test.Db.AddTask(&models.NewTask{Title: fmt.Sprintf("%d",i)})
		if err != nil { t.Error(err) }
	}
	for i := 0; i < expectedAmount; i++ {
		_,err := test.Db.AddTask(&models.NewTask{Title: fmt.Sprintf("%d",i), Due: 12312312})
		if err != nil { t.Error(err) }
	}
	
	for i := 0; i < 2; i++ {
		_,err := test.Db.AddTask(&models.NewTask{Title: fmt.Sprintf("%d",i)})
		if err != nil { t.Error(err) }
	}
	tasks, err := test.Db.GetDueTasks()
	if err != nil { t.Error(err) }

	if len(tasks) != expectedAmount {
		t.Errorf("expected Amount of due tasks (%d) is different than %d",expectedAmount,len(tasks))
	}
}