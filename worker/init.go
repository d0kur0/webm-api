package worker

import (
	"time"

	"github.com/spf13/viper"
)

func schedule(what func(), delay time.Duration, firstInstant bool) chan bool {
	stop := make(chan bool)

	if firstInstant {
		what()
	}

	go func() {
		for {
			what()
			select {
			case <-time.After(delay):
			case <-stop:
				return
			}
		}
	}()

	return stop
}

func init() {
	schedule(grabbing, time.Duration(viper.GetInt("updateInterval"))*time.Minute, true)
}
