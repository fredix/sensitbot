package db

import (
	//"log"
	"gopkg.in/mgo.v2"
	"time"
)

func GetSession(uri string, database string, dbowner string, dbpass string) *mgo.Session {
	//func GetSession(uri string) *mgo.Session {

	// Connect to our local mongo
	//session, err := mgo.Dial(uri)

	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{uri},
		Database: database,
		Username: dbowner,
		Password: dbpass,
		Timeout:  60 * time.Second,
	})

	// Check if connection error, is mongo running?
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	return session.Copy()
}
