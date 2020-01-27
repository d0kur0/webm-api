package grabberSchemaHelper

import (
	"os"

	"github.com/d0kur0/webm-grabber/sources/fourChannel"
	"github.com/d0kur0/webm-grabber/sources/twoChannel"

	"github.com/d0kur0/webm-grabber/sources/types"
)

type grabberSchema struct {
	defaultGrabberSchema     []types.GrabberSchema
	defaultAllowedExtensions types.AllowedExtensions
	configFilePath           string
}

func (schema *grabberSchema) parseConfig() (grabberSchema []types.GrabberSchema, err error) {
	// TODO: Доделать
	if _, err := os.Stat(schema.configFilePath); os.IsNotExist(err) {

	}

	return
}

func (schema *grabberSchema) Get() []types.GrabberSchema {
	return schema.defaultGrabberSchema
}

func Make() (schema *grabberSchema) {
	schema = &grabberSchema{configFilePath: "config.json"}

	schema.defaultAllowedExtensions = types.AllowedExtensions{".webm", ".mp4"}
	schema.defaultGrabberSchema = []types.GrabberSchema{
		{
			twoChannel.Make(schema.defaultAllowedExtensions),
			[]types.Board{
				{"b", "Бред"},
				{"vg", "Видео Игры"},
				{"a", "Аниме"},
				{"mu", "Музыка"},
				{"e", "Extreme Porn"},
				{"h", "Хентай"},
				{"fur", "Фурри"},
				{"kpop", "K-Pop"},
				{"asmr", "ASMR"},
			},
		},
		{
			fourChannel.Make(schema.defaultAllowedExtensions),
			[]types.Board{
				{"a", "Anime & Manga"},
				{"c", "Anime/Cute"},
				{"cgl", "Cosplay & EGL"},
				{"vg", "Video Game Generals"},
				{"co", "Comics & Cartoons"},
				{"g", "Technology"},
				{"b", "Random"},
				{"mu", "Music"},
				{"s", "Sexy Beautiful Women"},
				{"hc", "Hardcore"},
				{"h", "Hentai"},
				{"e", "Ecchi"},
				{"u", "Yuri"},
				{"d", "Hentai/Alternative"},
			},
		},
	}

	return
}
