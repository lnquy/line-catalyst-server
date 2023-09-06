package utils

import (
	"log"
	"os"
	"time"
	_ "time/tzdata"
)

var GlobalLocation *time.Location

func init() {
	locStr := os.Getenv("LOCATION")
	if locStr == "" {
		locStr = "Asia/Bangkok"
	}
	var err error
	GlobalLocation, err = time.LoadLocation(locStr)
	if err != nil {
		log.Panicf("failed to load location: %s", locStr)
	}
}
