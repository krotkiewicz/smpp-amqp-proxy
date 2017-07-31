package main

import (
	"esme/proxy"
	"esme/config"
	log "github.com/sirupsen/logrus"
)

func main() {
	config.Init()
	if config.Config.LoggingFormat == "JSON" {
		log.SetFormatter(&log.JSONFormatter{})
	}
	p := proxy.NewProxy()
	go p.Serve()
	forever := make(chan bool)
	<-forever
}
