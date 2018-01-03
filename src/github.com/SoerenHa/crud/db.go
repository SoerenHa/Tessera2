/*
	Package crud contains function for creating, reading, updating and deleting elements in the mongoDB
*/

package crud

import (
	"fmt"
	"time"
	"strconv"
	"math/rand"
	"encoding/xml"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Username 	string	`xml:"username"`
	Password	string	`xml:"password,omitempty"`
	Room		[]Room	`xml:"rooms"`
	Scene		[]Scene	`xml:"scenes"`
	Simulator	Simulator
}

type BaseDevice struct {
	Type	string `bson:"type"`
}

type Device struct {
	Id		bson.ObjectId	`bson:"id",json:"id",xml:"id"`
	Name	string			`bson:"name",json:"name",xml:"name"`
	Type	string			`xml:"type"`
	State	string			`bson:"state",json:"state",xml:"state"`
}

type Scene struct {
	Id		bson.ObjectId
	Name	string
	Time	time.Time
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

type Simulator struct {
	State		string			`bson:"state",json:"state"`
	SimTime 	time.Time		`json:"simTime"`
	TimeStep	time.Duration	`json:"timeStep"`

}

const DB = "HA17DB_Soeren_Hansen_550548"
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

func getDeviceIndex(device string) string {
	rooms := GetRooms()

	var deviceIndex string

	for _, room := range rooms {
		for i, v := range room.Device {
			if v.Id == bson.ObjectIdHex(device) {
				deviceIndex = strconv.Itoa(i)
			}
		}
	}

	return deviceIndex
}

func getUser() string {
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
				if device.Type == "Light" {
					if rng.Intn(2) == 1 {
						action = "on"
					} else {
						action = "off"
					}
				} else if device.Type == "Shutter" {
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

func deleteAction(deviceId string) {
	c := MongoSession.DB(DB).C("user")

	scenes := GetScenes()

	for _, scene := range scenes {
		for _, action := range scene.Action {
			if deviceId == action.DeviceId {
				query := bson.M{
					"username": getUser(),
					"scene.action.deviceid": deviceId,
				}

				change := bson.M{
					"$pull": bson.M{
						"scene.$.action": action,
					},
				}
				Err := c.Update(query, change)

				if Err != nil {
					fmt.Println(Err.Error())
				}
			}
		}
	}
}

/******************** PUBLIC FUNCTIONS ********************/

func Connect (url string) {
	MongoSession, Err = mgo.Dial(url)

	if Err != nil {
		panic(Err)
	}

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
			Simulator:	Simulator{
				State:	"paused",
				TimeStep: time.Second,
			},
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

func GetBaseDevices () []BaseDevice {
	c := MongoSession.DB(DB).C("baseDevice")

	var devices []BaseDevice
	Err = c.Find(nil).All(&devices)

	if Err != nil {
		fmt.Print(Err.Error())
	}

	return  devices
}

func GetRooms () []Room {
	c := MongoSession.DB(DB).C("user")

	var result User
	Err = c.Find(bson.M{"username": getUser()}).Select(bson.M{"room": 1}).One(&result)

	if Err != nil {
		fmt.Print(Err.Error())
	}

	return result.Room
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
		"room.id":  bson.ObjectIdHex(room),
	}
	c.Update(query, bson.M{"$set": bson.M{"room.$.name": name}})
}

func DeleteRoom (room string) {
	c := MongoSession.DB(DB).C("user")
	roomToDelete := getRoom(room)

	for _, device := range roomToDelete.Device {
		deleteAction(device.Id.Hex())
	}

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

	deleteAction(deviceToDelete.Id.Hex())

	query := bson.M{
		"username": getUser(),
		"room.id":  bson.ObjectIdHex(room),
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

func RenameDevice (device, name string) {
	c := MongoSession.DB(DB).C("user")

	index := getDeviceIndex(device)

	query := bson.M{
		"username": getUser(),
		"room.device.id": bson.ObjectIdHex(device),
	}

	change := bson.M{
		"$set": bson.M{
			"room.$.device." + index + ".name": name}}

	Err := c.Update(query, change)

	if Err != nil {
		fmt.Println(Err.Error())
	}
}

func UpdateDeviceState(device, state string) {
	c := MongoSession.DB(DB).C("user")

	index := getDeviceIndex(device)

	query := bson.M{
		"username":       getUser(),
		"room.device.id": bson.ObjectIdHex(device),
	}

	change := bson.M{
		"$set": bson.M{
			"room.$.device." + index + ".state": state,
		},
	}

	Err := c.Update(query, change)

	if Err != nil {
		fmt.Println(Err.Error())
	}
}

func InsertDevice (name, deviceType, room string) {
	c := MongoSession.DB(DB).C("user")

	var state string

	if deviceType == "Light" {
		state = "off"
	} else if deviceType == "Shutter" {
		state = "0"
	} else {
		state = "default"
	}

	deviceStruct := Device {
		Id:		bson.NewObjectId(),
		Name:	name,
		Type:	deviceType,
		State:	state,
	}

	query := bson.M{
		"username": getUser(),
		"room.id":  bson.ObjectIdHex(room),
	}

	change := bson.M{"$addToSet": bson.M{"room.$.device": deviceStruct}}

	Err := c.Update(query, change)

	if Err != nil {
		fmt.Println(Err.Error())
	}
}

func GetDevices() []Device {
	rooms := GetRooms()
	var devices []Device

	for _, room := range rooms {
		for _, device := range room.Device {
			devices = append(devices, device)
		}
	}

	return devices
}

func InsertScene (name string, time time.Time, offset int, daily, active bool) {
	c := MongoSession.DB(DB).C("user")

	scene := Scene{
		Id:		bson.NewObjectId(),
		Name:	name,
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
	xmlData, Err := xml.MarshalIndent(GetRooms(),"", "  ")
	xmlScenes, Err := xml.MarshalIndent(GetScenes(),"", "  ")

	if Err != nil {
		fmt.Print(Err.Error())
	}

	for _,x := range xmlScenes {
		xmlData = append(xmlData, x)
	}

	return xmlData
}

func StartSimulation (fff int, simTime time.Time) {
	c := MongoSession.DB(DB).C("user")

	timestep := time.Duration(fff) * time.Second

	query := bson.M{"username": getUser()}
	change := bson.M{
		"$set": bson.M{
			"simulator.state": "running",
			"simulator.simtime": simTime,
			"simulator.timestep": timestep,
		},
	}

	Err = c.Update(query, change)

	if Err != nil {
		fmt.Println(Err.Error())
	}
}

func GetSimulator () Simulator {
	c := MongoSession.DB(DB).C("user")

	var result User
	Err = c.Find(bson.M{"username": getUser()}).Select(bson.M{"simulator": 1}).One(&result)

	if Err != nil {
		fmt.Print(Err.Error())
	}

	return result.Simulator
}

func UpdateSimTime (simTime time.Time) {
	c := MongoSession.DB(DB).C("user")

	query := bson.M{"username": getUser()}
	change := bson.M{
		"$set": bson.M{
			"simulator.simtime": simTime,
		},
	}

	Err = c.Update(query, change)

	if Err != nil {
		fmt.Println(Err.Error())
	}
}

func ToggleSimulator () string {
	c := MongoSession.DB(DB).C("user")
	sim := GetSimulator()
	query := bson.M{"username": getUser()}
	var change bson.M

	if sim.State == "running" {
		change = bson.M{
			"$set": bson.M{
				"simulator.state": "paused",
			},
		}
	} else {
		change = bson.M{
			"$set": bson.M{
				"simulator.state": "running",
			},
		}
	}

	Err = c.Update(query, change)

	if Err != nil {
		fmt.Println(Err.Error())
	}

	if sim.State == "running" {
		return "paused"
	} else {
		return "running"
	}
}
