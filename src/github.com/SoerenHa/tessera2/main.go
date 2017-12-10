package main

import (
	"encoding/json"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"log"
	"net/http"
	"time"
	"github.com/SoerenHa/crud"
	"fmt"
)

func main() {
	//Verbindung zur Datenbank herstellen
	crud.Connect("localhost:27017")
	//Static file server
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))


	http.HandleFunc("/", handleRequest)
	err := http.ListenAndServe(":9090", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}


}


/* Globale Variablen **********************************************************/

var MongoSession *mgo.Session
var Err error
var myTemplates = template.Must(template.ParseGlob("./templates/*"))

/* Structs ********************************************************************/

type Users struct {
	ID       bson.ObjectId //`bson:"_id,omitempty"`
	Username string        //`json:"username"`
	Password string        //`json:"password"`
	Session  string
}

type HeadTemplate struct {
	Title string
}

type BodyTemplate struct {
	BaseEntities	[]crud.BaseDevice
	Rooms			[]crud.Room
}

type RoomTemplate struct {
	RoomId		string
	RoomName	string
}

type Test struct {
	Name	string
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	requestedPath := r.URL.Path[1:]

	switch requestedPath {
		case "test":
			test(w, r)
			break
	}
}

/* Handler Functions **********************************************************/

//func mainMenu(w http.ResponseWriter, r *http.Request) {
//	if r.Method == "GET" {
//		if checkLogInStatus(r) {
//			myTemplates.ExecuteTemplate(w, "head.html", HeadTemplate{Title: "tessera Menu"})
//			myTemplates.ExecuteTemplate(w, "mainMenu.html", BodyTemplate{Username: getUsername(r)})
//		} else {
//			myTemplates.ExecuteTemplate(w, "head.html", HeadTemplate{Title: "tessera Menu"})
//			myTemplates.ExecuteTemplate(w, "home.html", BodyTemplate{})
//		}
//
//	}
//
//}

func test(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		devices :=crud.GetBaseDevices()
		rooms := crud.GetRooms()


		myTemplates.ExecuteTemplate(w, "head.html", HeadTemplate{})
		myTemplates.ExecuteTemplate(w, "body.html", BodyTemplate{devices, rooms})
		myTemplates.ExecuteTemplate(w, "foot.html", BodyTemplate{})
	}

	if r.Method == "POST" {
		var resp struct{
			success string
		}
		resp.success = "ok"

		//var resp []byte

		r.ParseForm()

		action := r.FormValue("action")

		switch action {
		case "insertRoom":
			room := r.FormValue("room")
			crud.InsertRoom(room)
			http.Redirect(w, r, "http://localhost:9090/test", http.StatusFound)
			break
		case "deleteRoom":
			crud.DeleteRoom(r.FormValue("room"))
			http.Redirect(w, r, "http://localhost:9090/test", http.StatusFound)

			break
		case "insertDevice":
			device := r.FormValue("deviceName")
			deviceType := r.FormValue("deviceType")
			room := r.FormValue("deviceRoom")

			crud.InsertDevice(device, deviceType, room)
			http.Redirect(w, r, "http://localhost:9090/test", http.StatusFound)

			break
		case "getRooms":
			fmt.Print("Request")
			rooms := crud.GetRoomJson()
			fmt.Print(rooms)
			w.Header().Set("Content-Type", "application/json")
			w.Write(rooms)


		}




		//fmt.Fprint(w, resp.success)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		myTemplates.ExecuteTemplate(w, "head.html", HeadTemplate{Title: "tessera Login"})
		myTemplates.ExecuteTemplate(w, "login.html", BodyTemplate{})
		//POST
	} else {
		decoder := json.NewDecoder(r.Body)

		var user Users

		err := decoder.Decode(&user)
		if err != nil {
			panic(err)
		}
		c := MongoSession.DB("local").C("Users")

		//Actual User
		var actUser Users

		_ = c.Find(bson.M{"username": user.Username}).One(&actUser)

		if actUser.Username != "" && actUser.Password == user.Password {
			//fmt.Println("login")
			//startSession(w, r, actUser.Username)

		} else {
			http.Error(w, "<span style=\"color: red;\">Login fehlgeschlagen</span>", 401)
		}
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	//Cookie verfällt sofort
	expiration := time.Now()
	cookie := http.Cookie{Name: "session", Value: "", Expires: expiration}
	http.SetCookie(w, &cookie)

	c := MongoSession.DB("local").C("Sessions")
	c.RemoveAll(bson.M{"username": getUsername(r)})

	myTemplates.ExecuteTemplate(w, "head.html", HeadTemplate{Title: "tessera Home"})
	myTemplates.ExecuteTemplate(w, "home.html", BodyTemplate{})

}

//func register(w http.ResponseWriter, r *http.Request) {
//	if r.Method == "GET" {
//		myTemplates.ExecuteTemplate(w, "head.html", HeadTemplate{Title: "tessera Registrierung"})
//		myTemplates.ExecuteTemplate(w, "register.html", BodyTemplate{})
//
//	} else {
//		decoder := json.NewDecoder(r.Body)
//
//		var user Users
//
//		err := decoder.Decode(&user)
//		if err != nil {
//			panic(err)
//		}
//
//		c := MongoSession.DB("local").C("Users")
//
//		n, _ := c.Find(bson.M{"username": user.Username}).Count()
//
//		if n != 0 {
//			fmt.Fprint(w, "Username bereits vergeben")
//		} else {
//			fmt.Fprint(w, "Benutzername frei")
//		}
//	}
//}

func insertNewUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)

		var user Users

		err := decoder.Decode(&user)
		if err != nil {
			panic(err)
		}

		c := MongoSession.DB("local").C("Users")

		var id = bson.NewObjectId()
		_ = c.Insert(&Users{ID: id, Username: user.Username, Password: user.Password})
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		myTemplates.ExecuteTemplate(w, "head.html", HeadTemplate{Title: "tessera Home"})
		myTemplates.ExecuteTemplate(w, "home.html", BodyTemplate{})
	}
}

func account(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		myTemplates.ExecuteTemplate(w, "head.html", HeadTemplate{Title: "tessera Konto"})
		myTemplates.ExecuteTemplate(w, "account.html", BodyTemplate{})
	}
}

/* Andere Functions ***********************************************************/

func getUsername(r *http.Request) string {
	cookie, Err := r.Cookie("session")
	if Err != nil {
		panic(Err)
	}

	return crypt(cookie.Value)
}

func crypt(username string) string {
	var out string
	for i := 0; i < len(username); i++ {
		out += string([]rune(username)[i] ^ 10)
	}
	return out
}