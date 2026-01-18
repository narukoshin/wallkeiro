package main

import "wallkeiro/core"

// main is the entry point for the Wallkeiro application. It creates the
// profiles folder if it does not exist, and then starts the application by
// calling core.Start(). If there is an error creating the folder or starting
// the application, it panics with the error message.
func main(){
	var err error
	err = core.CreateFolder()
	if err != nil {
		panic(err)
	}
	err = core.Start()
	if err != nil {
		panic(err)
	}
}