package main

import "log"

//Check - panics if error is found
func Check(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
