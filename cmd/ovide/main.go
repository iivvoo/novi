package main

import (
	"log"
	"time"

	"github.com/iivvoo/ovim/logger"
	"github.com/iivvoo/ovim/ovide"
)

func main() {
	logger.OpenLog("ovide.log")
	log.Printf("Starting at %s\n", time.Now())
	defer logger.CloseLog()
	ovide.Run()
}
