package main

import (
	"log"
	"flag"
	"fmt"
	"os"

	cfg "sleeply-alive/internal/config"
)

func main() {
	// start := flag.Bool("boot", false, "e")
	var path string
	flag.StringVar(&path, "f", "","run with select file. default:"+cfg.Path)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  A demo program that requires at least one argument.\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NFlag() == 0 && flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "Error: no arguments provided")
		flag.Usage()
		os.Exit(2)
	}

	args := flag.Args()
	for i, arg := range args {
		if i==0 && arg=="run" {
			run()
		}
	}

	flag.Usage()
	os.Exit(2)
}

func run() {
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
