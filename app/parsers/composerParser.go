package parsers

import (
	"encoding/json"
	"os"
	"strings"
)

type Composer struct {
	Require map[string]string `json:"require"`
}

type ComposerPlugin struct {
	Name    string
	Version string
}

func ParseComposer(path string) ([]ComposerPlugin, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var composer Composer
	if err := json.NewDecoder(file).Decode(&composer); err != nil {
		return nil, err
	}

	var plugins []ComposerPlugin
	for pluginName, version := range composer.Require {
		// Only keep plugins from wpackagist-plugin
		if strings.HasPrefix(pluginName, "wpackagist-plugin/") {
			pluginName = strings.TrimPrefix(pluginName, "wpackagist-plugin/")
			plugins = append(plugins, ComposerPlugin{
				Name:    pluginName,
				Version: version,
			})
		}
	}

	return plugins, nil
}

