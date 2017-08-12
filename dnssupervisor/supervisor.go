package main

import (
	"flag"
	"log"
)

var (
	port = flag.Int("port", 9999, "The server port")
)

func main() {
	flag.Parse()

	supervisor := new(DnsSupervisor)
	supervisor.Run(*port)

	log.Println("done")
}
