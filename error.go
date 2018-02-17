package main

var err error

func PanicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
