package main

import (
	"os"
	"io"
	"log"
	"fmt"
	"time"
	"bufio"
	"strconv"
	"net/http"
	"encoding/json"
	"html/template"
	"github.com/SoerenHa/crud"
	"github.com/SoerenHa/utilities"
)

var Err error
var myTemplates = template.Must(template.ParseGlob("./templates/*"))

type headTemplate struct {
	Title string
}

type bodyTemplate struct {
	BaseDevices	[]crud.BaseDevice
	User		crud.User
}

type simTemplate struct {
	Simulator crud.Simulator
}

func main() {
	crud.Connect("localhost:27017")

	//Static file server
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/", handleRequest)
	Err = http.ListenAndServe(":9090", nil) // setting listening port
	if Err != nil {
		log.Fatal("ListenAndServe: ", Err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	requestedPath := r.URL.Path[1:]

	switch requestedPath {
	case "control":
		control(w, r)
		break
	case "simulator":
		simulator(w, r)
		break
	default:
		http.Redirect(w, r, "http://localhost:9090/control", http.StatusNotFound)
		break
	}
}

/* Handler Functions **********************************************************/

func control(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		devices :=crud.GetBaseDevices()
		user := crud.GetUserData()

		myTemplates.ExecuteTemplate(w, "head.html", headTemplate{"Control Hub"})
		myTemplates.ExecuteTemplate(w, "body.html", bodyTemplate{devices, user})
	}

	if r.Method == "POST" {
		var resp struct{
			Msg bool
		}

		r.ParseForm()
		action := r.FormValue("action")

		switch action {
		case "insertRoom":
			room := r.FormValue("room")
			crud.InsertRoom(room)
			http.Redirect(w, r, "http://localhost:9090/control", http.StatusFound)
			break
		case "renameRoom":
			name := r.FormValue("name")
			room := r.FormValue("room")

			if name != "" {
				crud.RenameRoom(room, name)
			}

			http.Redirect(w, r, "http://localhost:9090/control", http.StatusFound)
			break
		case "deleteRoom":
			crud.DeleteRoom(r.FormValue("room"))
			http.Redirect(w, r, "http://localhost:9090/control", http.StatusFound)
			break
		case "insertDevice":
			device := r.FormValue("deviceName")
			deviceType := r.FormValue("deviceType")
			room := r.FormValue("deviceRoom")

			crud.InsertDevice(device, deviceType, room)
			http.Redirect(w, r, "http://localhost:9090/control", http.StatusFound)

			break
		case "renameDevice":
			name := r.FormValue("name")
			device := r.FormValue("device")

			if name != "" {
				crud.RenameDevice(device, name)
			}

			http.Redirect(w, r, "http://localhost:9090/control", http.StatusFound)
			break
		case "deleteDevice":
			crud.DeleteDevice(r.FormValue("room"), r.FormValue("device"))
			http.Redirect(w, r, "http://localhost:9090/control", http.StatusFound)
			break
		case "insertScene":
			name := r.FormValue("sceneName")
			date := r.FormValue("sceneDate")
			offset, _ := strconv.Atoi(r.FormValue("offset"))
			var daily bool
			var dateTime time.Time

			if r.FormValue("daily") == "on" {
				daily = true
			} else {
				daily = false
			}

			if r.FormValue("clockTime") != "" {
				sceneTime := r.FormValue("clockTime")
				dateTime = utilities.ParseTime(date, sceneTime)
			} else {
				dateTime,_ = time.Parse(time.RFC3339, date + "T00:00:00+00:00")
			}

			crud.InsertScene(name, dateTime, offset, daily,true)

			http.Redirect(w, r, "http://localhost:9090/control", http.StatusFound)
			break
		case "updateState":
			device := r.FormValue("device")
			state := r.FormValue("state")

			crud.UpdateDeviceState(device, state)
			w.Header().Set("Content-Type", "application/json")

			resp.Msg = true
			js, _ := json.Marshal(resp)

			w.Write(js)
			break
		case "getDevices":
			sim := crud.GetSimulator()
			var devices []crud.Device

			if sim.State == "running" {
				devices = crud.GetDevices()

				var resp struct{
					Devices []crud.Device
					Msg		bool
				}

				resp.Devices = devices
				resp.Msg = true

				js, _ := json.Marshal(resp)

				w.Header().Set("Content-Type", "application/json")
				w.Write(js)
			} else {
				resp.Msg = false
				js, _ := json.Marshal(resp)
				w.Header().Set("Content-Type", "application/json")
				w.Write(js)
			}
			break
		}
	}
}

func simulator(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		sim := crud.GetSimulator()
		myTemplates.ExecuteTemplate(w, "head.html", headTemplate{"Simulator"})
		myTemplates.ExecuteTemplate(w, "simBody.html", simTemplate{sim})
	}

	if r.Method == "POST" {
		var resp struct{
			Msg string	`json:"msg"`
		}

		r.ParseForm()
		action := r.FormValue("action")

		switch action {
		case "start":
			simDate := r.FormValue("date")
			simClock := r.FormValue("time")
			fff := r.FormValue("fff") //fast forward factor
			dateTime := time.Time{}

			if simDate == "" || simClock == "" {
				dateTime = time.Now()
			} else {
				dateTime = utilities.ParseTime(simDate, simClock)
			}

			if fff == "" {
				fff = "1"
			}

			fastForward, _ := strconv.Atoi(fff)
			crud.StartSimulation(fastForward, dateTime)

			resp.Msg = "running"
			js, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			break
		case "toggle":
			resp.Msg = crud.ToggleSimulator()
			js, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			break
		case "getSim":
			sim := crud.GetSimulator()
			js, _ := json.Marshal(sim)
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			break
		case "getXML":
			xml := crud.CreateXML()

			file, Err := os.Create("tmp/toHuusDB.xml")
			if Err != nil {
				fmt.Print(Err.Error())
			}
			defer file.Close()
			wr := bufio.NewWriter(file)

			fmt.Fprintf(wr, string(xml))
			wr.Flush()

			file, Err = os.OpenFile("tmp/toHuusDB.xml", os.O_RDWR, 0644)

			if Err != nil {
				fmt.Print(Err.Error())
			}

			w.Header().Set("Content-Disposition", "attachment; filename=toHuusDB.xml")
			w.Header().Set("Content-Type", "tmp/text/xml")
			io.Copy(w, file)
			file.Close()
			os.Remove("tmp/toHuusDB.xml")
			break
		}
	}
}