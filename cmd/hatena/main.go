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

func run(filename string, alias string, debug bool, dryRun bool) error {
	config, err := hatena.LoadConfig()
	if err != nil {
		return err
	}

	httpclient := http.DefaultClient
	if debug {
		httpclient = &http.Client{Transport: &auth.DebugTransport{Base: http.DefaultTransport, Verbose: false}}
	}

	// oauth dance
	wrap := func(req *http.Request) error {
		client := auth.NewClient(config.ConsumerKey, config.ConsumerSecret)
		if config.ClientID == "" || config.ClientSecret == "" {
			credential, err := client.Auth(httpclient)
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

	app := hatena.App{Client: hatena.NewClient(config.HatenaID, config.BlogID, dryRun, httpclient, wrap)}
	latest, err := store.LoadCommit(config.HistFile, alias)
	if err != nil {
		return err
	}

	if alias == "" {
		alias = config.DefaultAlias
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

var aliasFlag = flag.String("alias", "", "alias name of uploaded gists")
var debugFlag = flag.Bool("debug", false, "debug")
var dryRunFlag = flag.Bool("dry-run", false, "dry-run")

func main() {
	flag.Parse()
	filename := flag.Arg(0)
	err := run(filename, *aliasFlag, *debugFlag, *dryRunFlag)
	if err != nil {
		log.Fatal(err)
	}
}
