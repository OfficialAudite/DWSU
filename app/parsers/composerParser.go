package parsers

import (
	"encoding/json"
	"os"
)

type Composer struct {
	Require map[string]string `json:"require"`
}

func ParseComposer(path string) (*Composer, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var composer Composer
	if err := json.NewDecoder(file).Decode(&composer); err != nil {
		return nil, err
	}
	return &composer, nil
}
