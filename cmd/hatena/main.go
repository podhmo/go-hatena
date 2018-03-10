package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/garyburd/go-oauth/oauth"
	"github.com/podhmo/hatena"
	"github.com/podhmo/hatena/auth"
	"github.com/podhmo/hatena/store"
)

func makeApp(config *hatena.Config, debug bool, dryRun bool) *hatena.App {
	httpclient := http.DefaultClient
	if debug {
		httpclient = &http.Client{Transport: &auth.DebugTransport{Base: http.DefaultTransport, Verbose: false}}
	}

	// oauth dance
	wrap := func(req *http.Request) error {
		client := auth.NewClient(config.ConsumerKey, config.ConsumerSecret)
		if config.ClientID == "" || config.ClientSecret == "" {
			credential, err := client.AuthDance(httpclient)
			if err != nil {
				return err
			}
			config.ClientID = credential.Token
			config.ClientSecret = credential.Secret
			hatena.SaveConfig(config)
		}

		credential := oauth.Credentials{Token: config.ClientID, Secret: config.ClientSecret}
		req.Header.Set("Authorization", client.AuthorizationHeader(&credential, req.Method, req.URL, nil))
		return nil
	}

	return &hatena.App{
		Client: hatena.NewClient(config.HatenaID, config.BlogID, dryRun, httpclient, wrap),
		Config: config,
	}
}

func list(app *hatena.App) error {
	return app.ListRecentlyArticles()
}

func post(app *hatena.App, filename string, alias string) error {
	latest, err := store.LoadCommit(app.Config.HistFile, alias)
	if err != nil {
		return err
	}

	if alias == "" {
		alias = app.Config.DefaultAlias
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
	return store.SaveCommit(app.Config.HistFile, commit)
}

var aliasFlag = flag.String("alias", "", "alias name of uploaded gists")
var debugFlag = flag.Bool("debug", false, "debug")
var dryRunFlag = flag.Bool("dry-run", false, "dry-run")
var listFlag = flag.Bool("list", false, "list latest entries")

func run() error {
	flag.Parse()
	config, err := hatena.LoadConfig()
	if err != nil {
		return err
	}
	app := makeApp(config, *debugFlag, *dryRunFlag)
	if *listFlag {
		return list(app)
	}

	filename := flag.Arg(0)
	return post(app, filename, *aliasFlag)
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
