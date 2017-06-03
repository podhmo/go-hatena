package hatena

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/podhmo/hatena/article"
	"github.com/podhmo/hatena/store"
)

// App :
type App struct {
	Client Client
}

// ListRecentlyArticles :
func (app *App) ListRecentlyArticles() error {
	entries, err := app.Client.List()
	if err != nil {
		return errors.Wrap(err, "client")
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(entries)
}

// CreateArticle :
func (app *App) CreateArticle(filename string, alias string) (store.Commit, error) {
	article, err := app.loadArticle(filename)
	if err != nil {
		return store.Commit{}, err
	}
	id, err := app.Client.Create(article)
	if err != nil {
		return store.Commit{}, errors.Wrap(err, "client")
	}
	return store.NewCommit(id, alias, "create"), nil
}

// EditArticle :
func (app *App) EditArticle(filename string, alias string, latestID string) (store.Commit, error) {
	article, err := app.loadArticle(filename)
	if err != nil {
		return store.Commit{}, err
	}
	id, err := app.Client.Edit(article, latestID)
	if err != nil {
		return store.Commit{}, errors.Wrap(err, "client")
	}
	return store.NewCommit(id, alias, "edit"), nil
}

func (app *App) loadArticle(filename string) (article.Article, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return article.Article{}, errors.Wrap(err, "open")
	}
	body := string(b)
	a, err := article.ParseArticle(body)
	if err != nil {
		return article.Article{}, errors.Wrap(err, "parse")
	}
	return a, nil
}
