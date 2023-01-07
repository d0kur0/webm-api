package worker

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"

	"github.com/spf13/viper"
)

func Init() {
	grabbing()

	s := gocron.NewScheduler(time.UTC)
	_, err := s.Every(viper.GetInt("updateInterval")).Minutes().Do(grabbing)
	if err != nil {
		log.Fatalln(err)
	}
	s.StartAsync()
}
