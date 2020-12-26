package main

import "log"

func Check(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
