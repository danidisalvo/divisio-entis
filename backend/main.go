package main

import (
	"backend/internal/graph"
	"backend/internal/rest"
	_ "embed"
	log "github.com/sirupsen/logrus"
)

const filename = "volume/graph.json"

// main starts the backend http server
func main() {
	g, err := graph.NewGraph("ens", filename)
	if err != nil {
		log.Fatalf("Failed to create the root node [%v]", err)
	}
	g.Load()

	server := rest.NewHttpServer(g)
	err = server.StartHttpServer()
	if err != nil {
		log.Fatalf(err.Error())
	}
}
