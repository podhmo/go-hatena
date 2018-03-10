package hatena

import (
	"github.com/pkg/errors"
	"github.com/podhmo/commithistory"
)

// Config is mapping object for application config
type Config struct {
	DefaultAlias string `json:"default_alias"`
	HistFile     string `json:"hist_file"`

	ConsumerKey    string `json:"consumer_key"`
	ConsumerSecret string `json:"consumer_secret"`

	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`

	HatenaID string `json:"hatena_id"`
	BlogID   string `json:"blog_id"`
}

const (
	defaultAlias    = "head"
	defaultHistFile = "hatena.history"
)

// ResolveAlias :
func (c *Config) ResolveAlias(alias string) string {
	if alias == "" {
		return c.DefaultAlias
	}
	return alias
}

// LoadConfig :
func LoadConfig(c *commithistory.Config) (*Config, error) {
	var conf Config
	if err := c.Load("config.json", &conf); err != nil {
		return nil, errors.Wrap(err, "load config")
	}
	if conf.DefaultAlias == "" {
		conf.DefaultAlias = defaultAlias
	}
	if conf.HistFile == "" {
		conf.HistFile = defaultHistFile
	}
	return &conf, nil
}

// SaveConfig :
func SaveConfig(c *commithistory.Config, config *Config) error {
	return c.Save("config.json", config)
}
