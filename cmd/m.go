package main

import (
	"fmt"
	"os"

	"github.com/echa/log"
)

func main() {
	buf := log.NewRingBuffer(16)
	log.Log.Attach(buf)
	log.SetLevel(log.LevelTrace)
	log.Tracef("HEllo")
	log.Debugf("HEllo")
	log.Infof("HEllo")
	log.Warnf("HEllo")
	log.Errorf("HEllo")

	fmt.Println("Catured History:")
	for msg := range buf.Last(16) {
		os.Stdout.Write(msg)
	}

	log.Panicf("HEllo")
}
