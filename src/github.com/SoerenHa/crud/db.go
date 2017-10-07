/*
	Package db contains utility functions for working with the database
*/

package crud

import (
	"fmt"
	"time"
)

func Connect() {

	var startT = time.Now()
	var end = false

	fmt.Println("Connecting to the database...")
	fmt.Printf("start ime is %v\n", startT)

	for !end {

		if time.Now().After(startT.Add(time.Duration(time.Second * 5))) {
			fmt.Println("Connected!")
			end = true
		}
	}
}