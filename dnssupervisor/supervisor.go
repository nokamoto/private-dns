package main

import (
	"flag"
	"log"
)

var (
	port = flag.Int("port", 9999, "The server port")
	hostsFile = flag.String("hostsfile", "/dev/null", "The dns hosts file")
	apiKey = flag.String("apikey", "default", "The request 'x-api-key'")
)

func main() {
	flag.Parse()

	supervisor := new(DnsSupervisor)
	supervisor.hostsFile = *hostsFile
	supervisor.apiKey = *apiKey
	supervisor.Run(*port)

	log.Println("done")
}
