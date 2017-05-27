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
	article := Article{}
	found := false
	buf := []string{}

	for _, line := range lines {
		if found {
			buf = append(buf, line)
			continue
		}
		if topLevelSectionRE.MatchString(line) {
			article.Title = parseTitle(line)
			found = true
		}
	}
	if !found {
		return article, errors.New("title is not found")
	}
	article.Body = strings.Trim(strings.Join(buf, "\n"), "\n")
	return article, nil
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
