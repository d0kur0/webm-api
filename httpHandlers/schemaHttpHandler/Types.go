package schemaHttpHandler

type responseBoard struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type responseSchema map[string][]responseBoard
