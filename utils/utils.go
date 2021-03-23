package utils

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Err(msg string, err error) error{
	return errors.New("["+msg+"] "+err.Error())
}

func IidtoObj(id interface{}) primitive.ObjectID{
	return id.(primitive.ObjectID)
}

func SendResponse(w http.ResponseWriter, response interface{}, statusCode int){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("JSON encoding error: %s",err) //TODO: Log file
	}
}