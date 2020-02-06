package grabberTask

import (
	"log"

	"github.com/d0kur0/webm-api/helpers/grabberSchemaHelper"

	"github.com/jasonlvhit/gocron"

	webmGrabber "github.com/d0kur0/webm-grabber"
	"github.com/d0kur0/webm-grabber/sources/types"
)

var output types.Output
var grabberSchema = grabberSchemaHelper.Make()

func Start(updateInterval uint64) {
	log.Println("First start, grabbing files...")
	output = webmGrabber.GrabberProcess(grabberSchema.Get())
	log.Println("Grabbing ended")

	log.Println("Start grabberTask")
	gocron.Every(updateInterval).Minutes().DoSafely(func() {
		output = webmGrabber.GrabberProcess(grabberSchema.Get())
		log.Println("GrabberTask: update files")
	})

	<-gocron.Start()
}

func GetGrabberSchema() []types.GrabberSchema {
	return grabberSchema.Get()
}

func GetOutput() types.Output {
	return output
}
