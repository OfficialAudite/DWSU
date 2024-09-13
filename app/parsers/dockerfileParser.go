package parsers

import (
	"bufio"
	"os"
	"strings"
)

func ParseDockerfile(path string) (phpVersion string, wordpressVersion string, wpPhpVersion string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return "", "", "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lowerLine := strings.ToLower(line)
		
		if strings.Contains(lowerLine, "from php:") {
			phpLine := strings.SplitN(lowerLine, "as", 2)[0]
			phpVersion = strings.TrimSpace(strings.Split(phpLine, ":")[1])
		}

		if strings.Contains(lowerLine, "from wordpress:") {
			wpLine := strings.SplitN(lowerLine, "as", 2)[0]
			fullWordpressVersion := strings.TrimSpace(strings.Split(wpLine, ":")[1])

			if strings.Contains(fullWordpressVersion, "-php") {
				parts := strings.Split(fullWordpressVersion, "-php")
				wordpressVersion = parts[0]          // WordPress version
				wpPhpVersion = parts[1]      // Extracted PHP version
			} else {
				wordpressVersion = fullWordpressVersion
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", "", "", err
	}

	return phpVersion, wordpressVersion, wpPhpVersion, nil
}
