/*
	Package db contains utility functions for working with the database
*/

package crud

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"encoding/json"
	"strconv"
)

type User struct {
	Username 	string
	Password	string
	Room		[]Room
}

type BaseDevice struct {
	Type	string `bson:"type"`
}

type Device struct {
	Id		bson.ObjectId	`bson:"id",json:"id"`
	Name	string			`bson:"name",json:"name"`
	Type	string
	State	string			`bson:"state",json:"state"`
}

type Room struct {
	Id		bson.ObjectId	`bson:"id",json:"id"`
	Name	string			`bson:"name",json:"name"`
	Device  []Device		`bson:"device"`
}

type RoomContainer struct {
	Rooms []Room
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

		light := BaseDevice{
			Type:	"Light",
		}

		shutter := BaseDevice{
			Type:	"Shutter",
		}

		coffeeMachine := BaseDevice{
			Type:	"Coffee machine",
		}

		Err = c.Insert(light, shutter, coffeeMachine)
	}

	c = MongoSession.DB(DB).C("user")

	smarty, _ := c.Find(bson.M{"username": "smarty"}).Count()

	if smarty == 0 {
		Err = c.Insert(User{
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
		Err = c.Insert(&User{Username: username, Password: password})
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
		fmt.Print(Err.Error())
	}

	return  devices
}

// Returns Array of Rooms
func GetRooms () []Room {
	c := MongoSession.DB(DB).C("user")

	user := getUser()

	var result User
	Err = c.Find(bson.M{"username": user}).Select(bson.M{"room": 1}).One(&result)

	if Err != nil {
		fmt.Print(Err.Error())
	}

	return result.Room
}

func GetRoomJson () []byte {
	rooms := GetRooms()
	json, Err := json.Marshal(rooms)

	if Err != nil {
		fmt.Print(Err.Error())
	}
	return json
}

func getUser () string {
	return "smarty"
}

func InsertRoom (room string) {
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

}

func DeleteRoom (room string) {
	c := MongoSession.DB(DB).C("user")

	var user User
	c.Find(bson.M{
		"username": getUser(),
	}).Select(bson.M{
		"room": bson.M{
			"$elemMatch": bson.M{
				"id": bson.ObjectIdHex(room),
			},
		},
	}).One(&user)

	var roomToDelete Room
	if len(user.Room) > 0 {
		roomToDelete = user.Room[0]
		c.Update(bson.M{"username": getUser()}, bson.M{"$pull": bson.M{"room": roomToDelete}})

	}
}

func UpdateState (room, device, state string) {
	c := MongoSession.DB(DB).C("user")

	var user User
	c.Find(bson.M{
		"username": getUser(),
	}).Select(bson.M{
		"room": bson.M{
			"$elemMatch": bson.M{
				"id": bson.ObjectIdHex(room),
			},
		},
	}).One(&user)

	var deviceIndex string

	for i, v := range user.Room[0].Device  {
		if v.Id == bson.ObjectIdHex(device) {
			deviceIndex = strconv.Itoa(i)
		}
	}

	query := bson.M{
		"username": getUser(),
		"room.id": bson.ObjectIdHex(room),
	}

	change := bson.M{
		"$set": bson.M{
			"room.$.device." + deviceIndex + ".state": state}}

	Err := c.Update(query, change)

	if Err != nil {
		fmt.Println(Err.Error())
	}
}

func InsertDevice (name, deviceType, room string) {
	c := MongoSession.DB(DB).C("user")

	deviceStruct := Device {
		Id:		bson.NewObjectId(),
		Name:	name,
		Type:	deviceType,
	}

	query := bson.M{
		"username": getUser(),
		"room.id": bson.ObjectIdHex(room),
	}

	change := bson.M{"$addToSet": bson.M{"room.$.device": deviceStruct}}

	Err := c.Update(query, change)

	if Err != nil {
		fmt.Println(Err.Error())
	}
}