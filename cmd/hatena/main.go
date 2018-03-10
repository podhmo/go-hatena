package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/garyburd/go-oauth/oauth"
	"github.com/pkg/errors"
	"github.com/podhmo/commithistory"
	"github.com/podhmo/hatena"
	"github.com/podhmo/hatena/auth"
)

func makeApp(c *commithistory.Config, config *hatena.Config, debug bool, dryRun bool) *hatena.App {
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
			return hatena.SaveConfig(c, config)
		}

		credential := oauth.Credentials{Token: config.ClientID, Secret: config.ClientSecret}
		req.Header.Set("Authorization", client.AuthorizationHeader(&credential, req.Method, req.URL, nil))
		return nil
	}

	return &hatena.App{
		C:      c,
		Client: hatena.NewClient(config.HatenaID, config.BlogID, dryRun, httpclient, wrap),
		Config: config,
	}
}

func list(app *hatena.App) error {
	return app.ListRecentlyArticles()
}

func findLatestCommit(c *commithistory.Config, filename, alias string) (*hatena.Commit, error) {
	var commit hatena.Commit
	if err := c.LoadCommit(filename, alias, &commit); err != nil {
		if c.IsNotFound(err) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "load commit")
	}
	return &commit, nil
}

func post(app *hatena.App, filename string, alias string) error {
	latest, err := findLatestCommit(app.C, app.Config.HistFile, alias)
	if err != nil {
		return err
	}
	if latest == nil {
		commit, err := app.CreateArticle(filename, app.Config.ResolveAlias(alias))
		if err != nil {
			return err
		}
		return app.C.SaveCommit(app.Config.HistFile, commit)
	}

	commit, err := app.EditArticle(filename, app.Config.ResolveAlias(alias), latest.ID)
	if err != nil {
		return err
	}
	return app.C.SaveCommit(app.Config.HistFile, commit)
}

var (
	aliasFlag   = flag.String("alias", "", "alias name of uploaded gists")
	debugFlag   = flag.Bool("debug", false, "debug")
	dryRunFlag  = flag.Bool("dry-run", false, "dry-run")
	listFlag    = flag.Bool("list", false, "list latest entries")
	profileFlag = flag.String("profile", "", "using another profile")
)

func run() error {
	flag.Parse()
	c := commithistory.New("hatena", commithistory.WithProfile(*profileFlag))
	config, err := hatena.LoadConfig(c)
	if err != nil {
		return err
	}
	if config.ConsumerKey == "" {
		fmt.Println("please setup consumekey")
		os.Exit(1)
	}

	app := makeApp(c, config, *debugFlag, *dryRunFlag)
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
