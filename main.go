package main

import (
	"log"
	"time"

	"simple-go-server/db"
	"simple-go-server/handler"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r, ": The server will attempt to recover after 5 seconds.")
			time.Sleep(5 * time.Second)
			main()
		}
	}()

	if err := db.Init(); err != nil {
		panic(err)
	}

	d, err := db.Get()
	if err != nil {
		panic(err)
	}
	defer d.Close()

	r := handler.GetRouter()
	r.LoadAll()

	r.Run(":3000")
}
