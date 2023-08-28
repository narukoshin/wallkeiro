package main

import "wallkeiro/core"

func main(){
	var err error
	err = core.Start()
	if err != nil {
		panic(err)
	}
}