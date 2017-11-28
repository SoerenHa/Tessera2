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
	Id			bson.ObjectId
	Username 	string
	Password	string
}

type Entity struct {
	Type 	string
	Room 	string
	Actions	string
	Name	string
	User 	string
}

type BaseDevice struct {
	Id		bson.ObjectId `json:"id"`
	Name	string
}

type Device struct {
	Id		bson.ObjectId
	Name	string
}

type Room struct {
	Id		bson.ObjectId
	Name	string
}

const DB = "tohuus"
var MongoSession *mgo.Session
var Err error

func Connect(url string) {
	MongoSession, Err = mgo.Dial(url)

	MongoSession.SetMode(mgo.Monotonic, true)

	// Make sure that The baseDevices and user smarty are available
	c := MongoSession.DB(DB).C("baseDevice")
	count, _ := c.Find(nil).Count()

	if count < 1 {

		light := Device{
			Id:		bson.NewObjectId(),
			Name:	"Light",
		}

		shutter := Device{
			Id:		bson.NewObjectId(),
			Name:	"Shutter",
		}

		coffeeMachine := Device{
			Id:		bson.NewObjectId(),
			Name:	"Coffee machine",
		}

		Err = c.Insert(light, shutter, coffeeMachine)
	}

	c = MongoSession.DB(DB).C("user")

	smarty, _ := c.Find(bson.M{"username": "smarty"}).Count()

	if smarty == 0 {
		Err = c.Insert(User{
			Id:			bson.NewObjectId(),
			Username:	"smarty",
			Password:	"123",
		})
	}

	if Err != nil {
		panic(Err)
	} else {
		fmt.Println("Connected!")
	}

}

func InsertUser (username, password string) {
	c := MongoSession.DB(DB).C("users")
	fmt.Println(Err)

	var result User
	Err = c.Find(bson.M{"name": username}).One(&result)

	fmt.Println(Err)

	if Err != nil {
		Err = c.Insert(&User{Username: username, Password: password, Id: bson.NewObjectId()})
		fmt.Println(Err)

	}
	if Err != nil {
		panic(Err)
	}
}

func GetBaseDevices () []BaseDevice {
	c := MongoSession.DB(DB).C("baseDevice")

	var devices []BaseDevice
	Err := c.Find(nil).All(&devices)

	if Err != nil {

	}

	return  devices
}

func GetRooms () []Room {
	c := MongoSession.DB(DB).C("user")

	user := getUser()
	var rooms []Room
	var haha = c.Find(bson.M{"username": user}).Select(bson.M{"room":1})

	if Err != nil {
	}

	fmt.Print(haha)

	return  rooms
}

func getUser () string {
	return "smarty"
}

func InsertRoom (room string) bool {
	c := MongoSession.DB(DB).C("user")

	user := getUser()
	roomStruct := Room {
		Id:		bson.NewObjectId(),
		Name:	room,
	}

	query 	:= bson.M{"username": user}
	change 	:= bson.M{"$addToSet": bson.M{"room": roomStruct}}

	Err := c.Update(query, change)

	if Err != nil {
		fmt.Println(Err.Error())
	}

	return true
}

func InsertDevice (device string) bool {
	c := MongoSession.DB("test").C("coll")

	user := getUser()
	deviceStruct := Room {
		Id:		bson.NewObjectId(),
		Name:	device,
	}

	query 	:= bson.M{"user": user}
	change 	:= bson.M{"$addToSet": bson.M{"device": deviceStruct}}

	Err := c.Update(query, change)

	if Err != nil {
		fmt.Println(Err.Error())
	}

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

func InsertEntity () {
	c := MongoSession.DB(DB).C("baseDevice")

	Err = c.Insert(BaseDevice{
		Id: 	bson.NewObjectId(),
		Name:	"Light",
	},
		BaseDevice{
			Id: 	bson.NewObjectId(),
			Name:	"Shutter",
		},
		BaseDevice{
			Id: 	bson.NewObjectId(),
			Name:	"Coffee machine",
		})

	if Err != nil {
		fmt.Printf(Err.Error())
	}
}