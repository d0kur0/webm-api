package filesDaemon

import (
	"log"

	webmGrabber "github.com/d0kur0/webm-grabber"
	"github.com/d0kur0/webm-grabber/sources/fourChannel"
	"github.com/d0kur0/webm-grabber/sources/twoChannel"
	"github.com/d0kur0/webm-grabber/sources/types"
)

var output types.Output
var grabberSchema []types.GrabberSchema

func Start() {
	allowedExtension := types.AllowedExtensions{".webm", ".mp4"}
	grabberSchema = []types.GrabberSchema{
		{
			twoChannel.Make(allowedExtension),
			[]types.Board{"b", "h", "fur"},
		},
		{
			fourChannel.Make(allowedExtension),
			[]types.Board{"b", "e", "h", "u"},
		},
	}

	log.Println("First start, grabbing files...")
	output = webmGrabber.GrabberProcess(grabberSchema)
	log.Println("Grabbing ended")
}

func GetGrabberSchema() []types.GrabberSchema {
	return grabberSchema
}

func GetOutput() types.Output {
	return output
}
