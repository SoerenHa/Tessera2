/*
	Package db contains utility functions for working with the database
*/

package crud

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"bytes"
)

type User struct {
	Username 	string
	Password	string
	ID			bson.ObjectId
}

type Entity struct {
	Type 	string
	Room 	string
	Actions	string
	Name	string
	User 	string
}

type BaseEntity struct {
	Name	string
}

type Room struct {
	Type	string
}

const DB = "tohuus"
var MongoSession *mgo.Session
var Err error

func Connect(url string) {
	MongoSession, Err = mgo.Dial(url)

	if Err != nil {
		panic(Err)
	} else {
		fmt.Println("Connected!")
	}

	MongoSession.SetMode(mgo.Monotonic, true)

}

func InsertUser (username, password string) {
	c := MongoSession.DB(DB).C("users")
	fmt.Println(Err)

	var result User
	Err = c.Find(bson.M{"name": username}).One(&result)

	fmt.Println(Err)

	if Err != nil {
		Err = c.Insert(&User{Username: username, Password: password, ID: bson.NewObjectId()})
		fmt.Println(Err)

	}
	if Err != nil {
		panic(Err)
	}
}

func GetBaseEntities () []BaseEntity  {
	c := MongoSession.DB(DB).C("baseEntities")

	var names []BaseEntity
	Err := c.Find(nil).Select(bson.M{"name": 1}).All(&names)

	if Err != nil {

	}

	return  names
}

func GetRooms () []Room {
	c := MongoSession.DB(DB).C("room")

	var rooms []Room
	Err := c.Find(nil).Select(bson.M{"type": 1}).All(&rooms)

	if Err != nil {

	}

	return  rooms
}

func getUser () string {
	return "Test"
}

func InsertRoom () bool {
	c := MongoSession.DB("test").C("coll")


	//if (Room{}) != r {
		query 	:= bson.M{"user": "test"}
		change 	:= bson.M{"$addToSet": bson.M{"rooms": "Fickraum1"}}

		Err := c.Update(query, change)

		if Err != nil {
			fmt.Println(Err.Error())
		}
	//}

	return true
}

func GetAllUsernames () []User {
	c := MongoSession.DB("tessera").C("users")

	var result []User
	Err = c.Find(bson.M{}).Sort("Username").All(&result)

	return result
}

func GetRandomName() string {

	var source = "1234567890abcdefghijklmnopqrstuvwxyz"
	var buffer bytes.Buffer


	for i := 0; i < 10; i++ {
		buffer.WriteString(string(source[rand.Intn(len(source) - 5) + 5]))
	}

	return buffer.String()
}

func insertEntity (entity Entity) {
	c := MongoSession.DB(DB).C("Entity")

	Err = c.Insert(entity)
}