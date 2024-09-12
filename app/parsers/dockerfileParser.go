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

		// Convert line to lowercase for case-insensitive "AS" detection
		lowerLine := strings.ToLower(line)

		// Check for "FROM php:" line and remove any trailing "AS" clause (case-insensitive)
		if strings.Contains(lowerLine, "from php:") {
			// Remove anything after "AS" (case-insensitive)
			phpLine := strings.SplitN(lowerLine, "as", 2)[0]
			phpVersion = strings.TrimSpace(strings.Split(phpLine, ":")[1])
		}

		// Check for "FROM wordpress:" line and remove any trailing "AS" clause (case-insensitive)
		if strings.Contains(lowerLine, "from wordpress:") {
			// Remove anything after "AS" (case-insensitive)
			wpLine := strings.SplitN(lowerLine, "as", 2)[0]
			fullWordpressVersion := strings.TrimSpace(strings.Split(wpLine, ":")[1])

			// Check if the WordPress version includes a PHP version (e.g., 6.6.1-php8.3)
			if strings.Contains(fullWordpressVersion, "-php") {
				parts := strings.Split(fullWordpressVersion, "-php")
				wordpressVersion = parts[0]          // WordPress version
				wpPhpVersion = parts[1]      // Extracted PHP version
			} else {
				wordpressVersion = fullWordpressVersion // No PHP version in WordPress
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", "", "", err
	}

	return phpVersion, wordpressVersion, wpPhpVersion, nil
}
