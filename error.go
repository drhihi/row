package main

import "log"

func PanicOnErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func printLog(value interface{}) {
	log.Println(value)
}
