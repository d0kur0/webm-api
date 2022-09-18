package worker

import (
	"github.com/d0kur0/webm-api/util"
	webmGrabber "github.com/d0kur0/webm-grabber"
	"github.com/d0kur0/webm-grabber/types"
)

var GrabbingOutPut types.Output

func grabbing() {
	GrabbingOutPut = webmGrabber.GrabberProcess(util.ParseSchema())
}
