package main

import (
	"log"
	"os"

	"github.com/nullrish/goauth/run"
)

func main() {
	if err := run.InitializeApp(); err != nil {
		log.Fatalln("Failed to start service:", err)
		os.Exit(1)
	}
}
