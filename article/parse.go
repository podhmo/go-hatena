package article

import (
	"errors"
	"regexp"
	"strings"
)

var topLevelSectionRE = regexp.MustCompile("#[^#]")
var parseRE = regexp.MustCompile("# *((?:\\[[^\\]]+\\])*)? *(.+)")

// ParseArticle :
func ParseArticle(body string) (Article, error) {
	lines := strings.Split(body, "\n")
	for _, line := range lines {
		if topLevelSectionRE.MatchString(line) {
			title := parseTitle(line)
			return Article{Title: title, Body: body}, nil
		}
	}
	return Article{}, errors.New("title is not found")
}

// parseTitle :
func parseTitle(line string) Title {
	separeted := parseRE.FindStringSubmatch(line)
	categoryArea := separeted[1]
	title := separeted[2]
	categories := parseCategories(categoryArea)
	return Title{Title: title, Categories: categories}
}

func parseCategories(s string) []string {
	if !strings.Contains(s, "]") {
		return []string{}
	}
	nodes := strings.Split(s, "]")
	categories := make([]string, 0, len(nodes))
	for _, node := range nodes {
		category := strings.TrimSpace(strings.TrimLeft(strings.TrimSpace(node), "["))
		if category == "" {
			continue
		}
		categories = append(categories, category)
	}
	return categories
}
