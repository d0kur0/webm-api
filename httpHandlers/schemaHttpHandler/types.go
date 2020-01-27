package schemaHttpHandler

type responseBoard struct {
	Name        string
	Description string
}

type responseSchema map[string][]responseBoard
