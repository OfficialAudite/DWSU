package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

type WordPressPlugin struct {
	Name          string            `json:"name"`
	Slug          string            `json:"slug"`
	LatestVersion string            `json:"version"`
	Versions      map[string]string `json:"versions"`
  VersionList   []string           `json:"-"`
}

type DockerImageTag struct {
	Name string `json:"name"`
}

type DockerImageTagsResponse struct {
	Results []DockerImageTag `json:"results"`
	Next    string           `json:"next"`
}

type PHPVersion struct {
	VersionID           int    `json:"versionId"`
	Name                string `json:"name"`
	ReleaseDate         string `json:"releaseDate"`
}

type PHPVersionData struct {
	Data map[string]PHPVersion `json:"data"`
}

func GetWordpressPluginVersions(pluginSlug string) (WordPressPlugin, error) {
	WordPressAPIUrl := "https://api.wordpress.org/plugins/info/1.2/?action=plugin_information&request[slug]=" + pluginSlug
	WordPressPlugin := WordPressPlugin{}

	res, err := http.Get(WordPressAPIUrl)
	if err != nil {
		return WordPressPlugin, err
	}
	defer res.Body.Close()

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return WordPressPlugin, readErr
	}

	err = json.Unmarshal(body, &WordPressPlugin)
	if err != nil {
		return WordPressPlugin, err
	}

	// Generate a list of version numbers from the map
	WordPressPlugin.VersionList = make([]string, 0, len(WordPressPlugin.Versions))
	for version := range WordPressPlugin.Versions {
		WordPressPlugin.VersionList = append(WordPressPlugin.VersionList, version)
	}

	return WordPressPlugin, nil
}

func GetDockerImageTags(repository string) ([]DockerImageTag, error) {
	var uniqueTags []DockerImageTag
	seenTags := make(map[string]struct{})
	url := fmt.Sprintf("https://hub.docker.com/v2/repositories/library/%s/tags", repository)

	pattern := regexp.MustCompile(`^\d+\.\d+(\.\d+)?-php\d+\.\d+$`)

	for url != "" {
		res, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		var dockerResponse DockerImageTagsResponse
		err = json.Unmarshal(body, &dockerResponse)
		if err != nil {
			return nil, err
		}

		for _, tag := range dockerResponse.Results {
			if pattern.MatchString(tag.Name) {
				if _, exists := seenTags[tag.Name]; !exists {
					seenTags[tag.Name] = struct{}{}
					uniqueTags = append(uniqueTags, tag)
				}
			}
		}

		url = dockerResponse.Next
	}

	return uniqueTags, nil
}

func GetPHPVersions() (PHPVersionData, error) {
	PHPWatchAPIUrl := "https://php.watch/api/v1/versions"
	phpVersionData := PHPVersionData{}

	res, err := http.Get(PHPWatchAPIUrl)
	if err != nil {
		return phpVersionData, err
	}
	defer res.Body.Close()

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return phpVersionData, readErr
	}

	err = json.Unmarshal(body, &phpVersionData)
	if err != nil {
		return phpVersionData, err
	}

	return phpVersionData, nil
}
