package models

import (
	"golang.org/x/oauth2"
	"gopkg.in/mgo.v2/bson"
)

type (
	// User represents the structure of our resource
	User struct {
		Id         bson.ObjectId `json:"id" bson:"_id"`
		TeleId     int           `json:"tele_id" bson:"tele_id"`
		FirstName  string        `json:"first_name" bson:"first_name"`
		LastName   string        `json:"last_name" bson:"last_name"`
		Username   string        `json:"username" bson:"username"`
		Sensitcode string        `json:"sensit_code" bson:"sensit_code"`
		Token      oauth2.Token  `json:"token" bson:"token"`
		Scope      string        `json:"scope" bson:"scope"`
		ExpiresIn  int           `json:"expires_in" bson:"expires_in"`
		Email      string        `json:"email" bson:"email"`
	}
)
