package database

import (
	"context"
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

	return &Database{client,coll},nil
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

	task := models.Task{
		Title: title,
		Desc: desc,
		Done: false,
	}

	result, err := db.coll.InsertOne(ctx,&task)
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

func (db *Database) RemoveTask(id primitive.ObjectID) (*mongo.DeleteResult, error){
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return db.coll.DeleteOne(
		ctx,
		bson.M{"_id": bson.M{"$eq": id}},
	)
}