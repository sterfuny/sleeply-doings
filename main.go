package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	cfg "sleeply-alive/internal/config"
)

func main() {
	var path string
	flag.StringVar(&path, "f", "", "select cfg file default:"+cfg.Path)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(2)
	} else {
		if os.Args[1] == "run" {
			os.Args = append([]string{os.Args[0]}, os.Args[2:]...)
			flag.Parse()
			fmt.Println("path:", path)
			cfg.Init(path)
			run()
			return
		}
	}

	flag.Usage()
	os.Exit(2)
}

func run() {
	if cfg.Get().Mode == "server" {
		startServer()
		return
	}

	if cfg.Get().Mode == "client" && cfg.Get().Addrs != nil {
		startClients()
	} else {
		log.Printf("Addrs值需要至少一项")
		return
	}
	select {}
}
