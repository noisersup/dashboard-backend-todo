package utils

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Err(msg string, err error) error{
	return errors.New("["+msg+"] "+err.Error())
}

func IidtoObj(id interface{}) primitive.ObjectID{
	return id.(primitive.ObjectID)
}
