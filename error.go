package main

import "log"

var err error

func PanicOnErr(err error) {
	if err != nil {
		log.Println(err)
	}
}
