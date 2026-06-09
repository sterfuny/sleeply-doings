package main

import (
	"log"

	cfg "sleeply-alive/internal/config"
)

func main() {
	if cfg.Get().Mode == "server" {
		startServer()
	}

	if cfg.Get().Mode == "client" && cfg.Get().Addrs != nil {
		startClients()
	} else {
		log.Printf("Addrs值需要至少一项")
		return
	}
	select {}
}
