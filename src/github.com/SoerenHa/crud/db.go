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
	"github.com/kelvins/sunrisesunset"
	"time"
	"encoding/xml"
	"math/rand"
)

type User struct {
	Username 	string	`xml:"username"`
	Password	string	`xml:"password,omitempty"`
	Room		[]Room	`xml:"rooms"`
	Scene		[]Scene	`xml:"scenes"`
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

type Scene struct {
	Id		bson.ObjectId
	Name	string
	Date	time.Time
	Time	string
	Offset	int
	Daily	bool
	Active	bool
	Action	[]Action
}

type Action struct {
	DeviceId	string
	Action		string
}

type Room struct {
	Id		bson.ObjectId	`bson:"id",json:"id",xml:"id"`
	Name	string			`bson:"name",json:"name",xml:"name"`
	Device  []Device		`bson:"device",xml:"devices"`
}

type RoomContainer struct {
	Rooms []Room
}

type Times struct {
	Sunrise	time.Time
	Sunset	time.Time
}

const DB = "tohuus"
var MongoSession *mgo.Session
var Err error

/******************** PRIVATE FUNCTIONS ********************/

func getRoom (room string) Room {
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

	return user.Room[0]
}

func getDeviceIndex(room, device string) string {
	roomStruct := getRoom(room)

	var deviceIndex string

	for i, v := range roomStruct.Device  {
		if v.Id == bson.ObjectIdHex(device) {
			deviceIndex = strconv.Itoa(i)
		}
	}

	return deviceIndex
}

func getUser () string {
	return "smarty"
}

func getRandomDeviceWithAction () []Action {
	rooms := GetRooms()
	seed := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(seed)

	if len(rooms) > 0 {

		//determine, from how many rooms devices should get picked (1 - rooms)
		roomsToSelect := rng.Intn(len(rooms)) + 1

		var selectedRooms []Room
		for len(selectedRooms) < roomsToSelect {
			tmp := rng.Intn(len(rooms))
			possibleRoom := rooms[tmp]

			if !contains(selectedRooms, possibleRoom) {
				selectedRooms = append(selectedRooms, rooms[tmp])
			}
		}

		var actions []Action

		for _, room := range selectedRooms {
			for _, device := range room.Device {
				var action string
				if device.Type == "Light" || device.Type == "Shutter" {
					action = strconv.Itoa(rng.Intn(100))
				} else {
					coffeeActions := [4]string{"Coffee", "Espresso", "Latte", "Cocoa"}
					action = coffeeActions[rng.Intn(len(coffeeActions))]
				}

				var tmp = Action{
					DeviceId: 	device.Id.Hex(),
					Action:		action,
				}

				actions = append(actions, tmp)
			}
		}

		return actions
	} else {
		return []Action{}
	}
}

func contains(s []Room, e Room) bool {
	for _, a := range s {
		if a.Id.Hex() == e.Id.Hex() {
			return true
		}
	}
	return false
}

/******************** PUBLIC FUNCTIONS ********************/

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

func GetUserData () User {
	c := MongoSession.DB(DB).C("user")

	var user User
	Err = c.Find(bson.M{"username": getUser()}).One(&user)

	if Err != nil {
		fmt.Print(Err.Error())
	}

	return user
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
	Err = c.Find(nil).All(&devices)

	if Err != nil {
		fmt.Print(Err.Error())
	}

	return  devices
}

// Returns Array of Rooms
func GetRooms () []Room {
	c := MongoSession.DB(DB).C("user")

	var result User
	Err = c.Find(bson.M{"username": getUser()}).Select(bson.M{"room": 1}).One(&result)

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

func RenameRoom (room, name string) {
	c := MongoSession.DB(DB).C("user")
	//roomStruct := getRoom(room)
	//roomStruct.Name = name
	query := bson.M{
		"username": getUser(),
		"room.id": bson.ObjectIdHex(room),
	}
	c.Update(query, bson.M{"$set": bson.M{"room.$.name": name}})
}

func DeleteRoom (room string) {
	c := MongoSession.DB(DB).C("user")
	roomToDelete := getRoom(room)
	c.Update(bson.M{"username": getUser()}, bson.M{"$pull": bson.M{"room": roomToDelete}})
}

func DeleteDevice (room, device string) {
	c := MongoSession.DB(DB).C("user")

	roomStruct := getRoom(room)

	var deviceToDelete Device
	for _, v := range roomStruct.Device {
		if v.Id == bson.ObjectIdHex(device) {
			deviceToDelete = v
		}
	}

	query := bson.M{
		"username": getUser(),
		"room.id": bson.ObjectIdHex(room),
	}

	change := bson.M{
		"$pull": bson.M{
			"room.$.device": deviceToDelete,
		},
	}

	Err := c.Update(query, change)

	if Err != nil {
		fmt.Println(Err.Error())
	}

}

func RenameDevice (room, device, name string) {
	c := MongoSession.DB(DB).C("user")

	index := getDeviceIndex(room, device)

	query := bson.M{
		"username": getUser(),
		"room.id": bson.ObjectIdHex(room),
	}

	change := bson.M{
		"$set": bson.M{
			"room.$.device." + index + ".name": name}}

	Err := c.Update(query, change)

	if Err != nil {
		fmt.Println(Err.Error())
	}
}

func UpdateState (room, device, state string) {
	c := MongoSession.DB(DB).C("user")

	index := getDeviceIndex(room, device)

	query := bson.M{
		"username": getUser(),
		"room.id": bson.ObjectIdHex(room),
	}

	change := bson.M{
		"$set": bson.M{
			"room.$.device." + index + ".state": state}}

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

func InsertScene (name, time string, timestamp time.Time, offset int, daily, active bool) {
	c := MongoSession.DB(DB).C("user")

	scene := Scene{
		Id:		bson.NewObjectId(),
		Name:	name,
		Date:	timestamp,
		Time:	time,
		Offset:	offset,
		Daily:	daily,
		Active: active,
		Action:	getRandomDeviceWithAction(),
	}

	query 	:= bson.M{"username": getUser()}
	change 	:= bson.M{"$addToSet": bson.M{"scene": scene}}

	Err := c.Update(query, change)

	if Err != nil {
		fmt.Println(Err.Error())
	}
}

func GetScenes () []Scene {
	c := MongoSession.DB(DB).C("user")

	var result User
	Err = c.Find(bson.M{"username": getUser()}).Select(bson.M{"scene": 1}).One(&result)

	if Err != nil {
		fmt.Print(Err.Error())
	}

	return result.Scene
}

func CreateXML () []byte {
	c := MongoSession.DB(DB).C("user")

	var foo User
	Err = c.Find(bson.M{"username": getUser()}).One(&foo)

	xmlString, Err := xml.MarshalIndent(foo," ", "  ")

	if Err != nil {
		fmt.Print(Err.Error())
	}

	return xmlString
}

func GetTimes () Times {
	timeData := sunrisesunset.Parameters{
		Latitude:  54.774727032665766,
		Longitude: 9.447391927242279,
		UtcOffset: 1.0,
		Date:      time.Now(),
	}

	sunrise, sunset, Err := timeData.GetSunriseSunset()

	if Err != nil {
		fmt.Print(Err.Error())
	}

	return Times{
		Sunrise: sunrise,
		Sunset:	 sunset,
	}
}
