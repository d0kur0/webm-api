package grabberSchemaHelper

import "github.com/d0kur0/webm-grabber/sources/types"

type configStruct struct {
	AllowedExtensions []types.AllowedExtensions
	GrabberSchema     map[string][]string
}
