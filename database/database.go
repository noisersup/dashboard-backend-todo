package database

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"kwiatek.xyz/todo-backend/models"
	u "kwiatek.xyz/todo-backend/utils"
)

type Database struct{
	client	*mongo.Client
	coll 	*mongo.Collection
	last	int
}

func ConnectToDatabase(uri string, dbName string, collName string) (*Database,error){
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err !=nil{ return nil,u.Err("NewClient",err) }

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	err = client.Connect(ctx)
	if err !=nil{ return nil,u.Err("Connect",err) }

	err = client.Ping(ctx, readpref.Primary())
	if err !=nil{ return nil,u.Err("Ping",err) }
 
	coll := client.Database(dbName).Collection(collName)

	return &Database{client,coll,0},nil
}

func (db *Database) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return db.client.Disconnect(ctx)
}

func (db *Database) GetTasks() ([]*models.Task, error){
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var tasks []*models.Task

	cur, err := db.coll.Find(ctx, bson.D{{}},options.Find())
	if err != nil {return nil,err}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for cur.Next(ctx){
		var task models.Task
		err := cur.Decode(&task)
		if err != nil {return nil, err}

		tasks = append(tasks, &task)
	}
	if err := cur.Err(); err != nil{
		return nil,err
	}
	
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cur.Close(ctx)

	return tasks, nil
}

func (db *Database) AddTask(title string, desc string) (*mongo.InsertOneResult,error){
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	task := models.NewTask{
		Title: title,
		Desc: desc,
		Done: false,
		Order: db.last+1,
	}

	result, err := db.coll.InsertOne(ctx,&task)
	if err == nil{
		db.last++
	}
	log.Printf("Last: %d",db.last)
	return result, err
}

func (db *Database) DoneTask(id primitive.ObjectID) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	return db.coll.UpdateOne(
		ctx,
		bson.M{"_id": bson.M{"$eq": id}},
		bson.M{"$set": bson.M{"done": true}},
	)
}

func (db *Database) UndoneTask(id primitive.ObjectID) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	return db.coll.UpdateOne(
		ctx,
		bson.M{"_id": bson.M{"$eq": id}},
		bson.M{"$set": bson.M{"done": false}},
	)
}
func (db *Database) GetTask(id primitive.ObjectID) (*models.Task, error){
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var task models.Task

	res := db.coll.FindOne(ctx, bson.M{"_id": bson.M{"$eq": id}}).Decode(&task)
	
	return &task, errors.New(res.Error())
}

func (db *Database) ReplaceTask(id primitive.ObjectID, secondId primitive.ObjectID) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	task1,err := db.GetTask(id)
	if err != nil {return nil, err}

	task2,err := db.GetTask(secondId)
	if err != nil {return nil, err}

	res,err := db.coll.UpdateOne(
		ctx,
		bson.M{"_id": bson.M{"$eq": id}},
		bson.M{"$set": bson.M{"order": task2.Order}},
	)
	if err != nil { return res,err }
	return db.coll.UpdateOne(
		ctx,
		bson.M{"_id": bson.M{"$eq": secondId}},
		bson.M{"$set": bson.M{"order": task1.Order}},
	)
}

func (db *Database) RemoveTask(id primitive.ObjectID) (*mongo.DeleteResult, error){
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return db.coll.DeleteOne(
		ctx,
		bson.M{"_id": bson.M{"$eq": id}},
	)
}
//BLOAT (+not working)
// func (db *Database) ReorderTasks() error{
// 	pointers,err := db.GetTasks()
// 	if err != nil {return err}
// 	var tasks []models.Task 

// 	for _,pointer := range pointers{
// 		tasks= append(tasks, *pointer)
// 	}

// 	for i:=0;i<len(tasks)-1;i++ {
// 		if tasks[i].Order > tasks[i+1].Order{
// 			buf := tasks[i]
// 			tasks[i]=tasks[i+1]
// 			tasks[i+1]=buf
// 		}
// 	}
// 	offset :=0
// 	for i,task := range tasks{
// 		if i != int(task.Order)+offset{
// 			_,err := db.coll.UpdateMany(
// 				context.TODO(),
// 				bson.M{"order": bson.M{"$gt": i}},
// 				bson.M{"$inc": bson.M{"order":-1}},
// 			)
// 			if err != nil {return err} 
// 			offset++
// 		}
// 	}
// 	err = db.GetLastOrder()
// 	if err != nil {return err}
// 	return nil
// }

func (db *Database) GetLastOrder() error{
	pointers,err := db.GetTasks()
	if err != nil {return err}
	var tasks []models.Task 

	for _,pointer := range pointers{
		tasks= append(tasks, *pointer)
	}

	if len(tasks)==0{
		db.last=0
		log.Printf("Last: %d",db.last)
		return nil
	}
	for i:=0;i<len(tasks)-1;i++ {
		if tasks[i].Order > tasks[i+1].Order{
			buf := tasks[i]
			tasks[i]=tasks[i+1]
			tasks[i+1]=buf
		}
	}
	db.last = tasks[len(tasks)-1].Order
	log.Printf("Last: %d",db.last)
	return nil
}