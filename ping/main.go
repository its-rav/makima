// print ping each 10 seconds

package main

import (
	"time"

	logger "github.com/its-rav/makima/pkg/logger"
)

func main() {
	logger.InitLogrusLogger()
	var log = logger.Log
	for {
		log.Infof("[Ping] %s", time.Now().Format(time.RFC3339))
		time.Sleep(10 * time.Second)
	}
}
