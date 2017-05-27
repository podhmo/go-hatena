package main

import (
	"flag"
	"log"

	"github.com/podhmo/hatena"
	"github.com/podhmo/hatena/store"
)

var aliasFlag = flag.String("alias", "", "alias name of uploaded gists")

func run(filename string, alias string) error {
	config, err := hatena.LoadConfig()
	if err != nil {
		return err
	}

	app := hatena.App{Client: &hatena.DummyClient{}}
	if alias == "" {
		alias = config.DefaultAlias
	}
	latest, err := store.LoadCommit(config.HistFile, alias)
	if err != nil {
		return err
	}

	var commit store.Commit
	if latest == nil {
		commit, err = app.CreateArticle(filename, alias)
	} else {
		commit, err = app.EditArticle(filename, alias, latest.ID)
	}
	if err != nil {
		return err
	}
	return store.SaveCommit(config.HistFile, commit)
}

func main() {
	flag.Parse()
	filename := flag.Arg(0)
	err := run(filename, *aliasFlag)
	if err != nil {
		log.Fatal(err)
	}
}
