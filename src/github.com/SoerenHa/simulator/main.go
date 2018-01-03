package main

import (
	"fmt"
	"time"
	"github.com/SoerenHa/crud"
	"github.com/SoerenHa/utilities"
)

func main() {
	crud.Connect("localhost:27017")

	for true {
		time.Sleep(time.Second)

		sim := crud.GetSimulator()

		if sim.State == "running" {
			scenes := crud.GetScenes()
			sim.SimTime = sim.SimTime.Add(sim.TimeStep)
			crud.UpdateSimTime(sim.SimTime)

			for _, scene := range scenes {
				if scene.Active {
					if 	utilities.ExecutionIsNow(sim.SimTime, scene.Time, sim.TimeStep, scene.Daily) {
						for _, action := range scene.Action {
							crud.UpdateDeviceState(action.DeviceId, action.Action)
						}
					}
				}
			}
		}
	}
}
