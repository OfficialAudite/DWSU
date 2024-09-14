package parsers

import (
	"encoding/json"
	"os"
	"strings"
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

	for pluginName, version := range composer.Require {
		if !strings.HasPrefix(pluginName, "wpackagist-plugin/") {
			delete(composer.Require, pluginName)
			continue
		}
		
		composer.Require[strings.TrimPrefix(pluginName, "wpackagist-plugin/")] = version
		delete(composer.Require, pluginName)
	}

	return &composer, nil
}
