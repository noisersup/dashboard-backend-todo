package models

import "go.mongodb.org/mongo-driver/bson/primitive"


type Task struct{
	ID			primitive.ObjectID `bson:"_id" json:"id"`
	Title		string `json:"title"`
	Desc		string `json:"desc"`
	Done		bool `json:"done"`
	Order		int `json:"order"` //TODO: Move this to collection struct when you'll create it
}

type NewTask struct{
	Title		string `json:"title"`
	Desc		string `json:"desc"`
	Done		bool `json:"done"`
	Order		int `json:"order"` 
}

type DoneTask struct {
	Done		bool `json:"done"`
}

type GetResponse struct{
	Tasks	[]Task `json:"tasks"`
	Error	string `json:"error"`
}

type ErrorResponse struct{
	Error	string `json:"error"`
}