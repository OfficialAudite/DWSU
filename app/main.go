package main

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/OfficialAudite/DWSU/app/controllers"
	"github.com/OfficialAudite/DWSU/app/parsers"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	wordpressPlugins  []parsers.ComposerPlugin
	wordpressVersion  string
	PhpVersion        string
	width             int
	height            int
	cursor            int
	selectingVersion  bool
	selectedPlugin    int
	pluginVersions    map[string]string
	selectedVersion   int
	scrollOffset      int
	versionScroll     int
	visibleItems      int
}

func New(wordpressPlugins []parsers.ComposerPlugin) *model {
	return &model{wordpressPlugins: wordpressPlugins}
}

func (m model) Init() tea.Cmd {
	return nil
}

func getVersionKeys(versions map[string]string) []string {
	keys := make([]string, 0, len(versions))
	for key := range versions {
		keys = append(keys, key)
	}
	return keys
}

func compareVersions(v1, v2 string) bool {
	v1Parts := strings.Split(v1, ".")
	v2Parts := strings.Split(v2, ".")

	for i := 0; i < len(v1Parts) && i < len(v2Parts); i++ {
		if v1Parts[i] != v2Parts[i] {
			return v1Parts[i] > v2Parts[i]
		}
	}

	return len(v1Parts) > len(v2Parts)
}

func sortVersions(versions []string) {
	sort.Slice(versions, func(i, j int) bool {
		return compareVersions(versions[i], versions[j])
	})
}


func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.visibleItems = 10

	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if !m.selectingVersion {
				if m.cursor > 0 {
					m.cursor--
					if m.cursor < m.scrollOffset {
						m.scrollOffset--
					}
				}
			} else if m.selectedVersion > 0 {
				m.selectedVersion--
				if m.selectedVersion < m.versionScroll {
					m.versionScroll--
				}
			}

		case "down":
			if !m.selectingVersion {
				if m.cursor < len(m.wordpressPlugins)-1 {
					m.cursor++
					if m.cursor >= m.scrollOffset+m.visibleItems {
						m.scrollOffset++
					}
				}
			} else if m.selectedVersion < len(m.pluginVersions)-1 {
				m.selectedVersion++
				if m.selectedVersion >= m.versionScroll+m.visibleItems {
					m.versionScroll++
				}
			}

		case "enter":
			if !m.selectingVersion {
				m.selectedPlugin = m.cursor
				selectedPlugin := m.wordpressPlugins[m.cursor].Name
				plugin, err := controllers.GetWordpressPluginVersions(selectedPlugin)
				if err != nil {
					fmt.Println("Error fetching plugin versions:", err)
					return m, nil
				}

				m.pluginVersions = plugin.Versions
				m.selectingVersion = true
				m.selectedVersion = 0
				m.versionScroll = 0
			} else {
				m.selectingVersion = false
			}

		case "q":
			if m.selectingVersion {
				m.selectingVersion = false
			} else {
				return m, tea.Quit
			}

		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.width == 0 {
		return "loading..."
	}

	var sb strings.Builder

	if m.selectingVersion {
		sb.WriteString(fmt.Sprintf("Available versions for %s (current version: %s):\n\n", m.wordpressPlugins[m.selectedPlugin].Name, m.wordpressPlugins[m.selectedPlugin].Version))

		versionKeys := getVersionKeys(m.pluginVersions)
		sortVersions(versionKeys)

		latestIndex := -1
		for i, version := range versionKeys {
			if version != "trunk" {
				latestIndex = i
				break
			}
		}

		for i := m.versionScroll; i < m.versionScroll+m.visibleItems && i < len(versionKeys); i++ {
			cursor := " "
			if m.selectedVersion == i {
				cursor = ">"
			}

			currentMarker := ""
			if versionKeys[i] == m.wordpressPlugins[m.selectedPlugin].Version {
				currentMarker = " (current)"
			}

			latestMarker := ""
			if i == latestIndex {
				latestMarker = " (latest)"
			}

     	sb.WriteString(fmt.Sprintf("%s %s%s%s\n", cursor, versionKeys[i], currentMarker, latestMarker))
		}

		sb.WriteString("\nPress Enter to select version or q to go back.")
		return sb.String()
	}

	sb.WriteString("WordPress Plugins (Navigate with up/down, press Enter to select):\n\n")
	for i := m.scrollOffset; i < m.scrollOffset+m.visibleItems && i < len(m.wordpressPlugins); i++ {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		sb.WriteString(fmt.Sprintf("%s %s (current version: %s)\n", cursor, m.wordpressPlugins[i].Name, m.wordpressPlugins[i].Version))
	}

	return sb.String()
}

func main() {
	mainlocation := "./test"

	f, err := tea.LogToFile("dwsu.log", "/tmp")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	plugins, err := parsers.ParseComposer(mainlocation + "/composer.json")

	if err != nil {
		panic(err)
	}

	p := tea.NewProgram(model{wordpressPlugins: plugins}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
