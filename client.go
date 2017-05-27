package hatena

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/pkg/errors"
	"github.com/podhmo/hatena/article"
)

// Client :
type Client interface {
	Create(article article.Article) (string, error)
	Edit(article article.Article, ID string) (string, error)
}

// ClientConfig :
type ClientConfig struct {
	HatenaID string
	BlogID   string
}

// NewClient :
func NewClient(hatenaID, blogID string, dryRun bool, client *http.Client, auth func(*http.Request) error) Client {
	config := ClientConfig{HatenaID: hatenaID, BlogID: blogID}
	if dryRun {
		return &dummyClient{Config: config}
	}
	return &actualClient{Config: config, Client: client, Auth: auth}
}

// dummyClient :
type dummyClient struct {
	Config ClientConfig
}

// Create :
func (c *dummyClient) Create(article article.Article) (string, error) {
	fmt.Println("Create: ")
	err := tmpl.Execute(os.Stdout, article)
	return "xxx", err
}

// Edit :
func (c *dummyClient) Edit(article article.Article, ID string) (string, error) {
	fmt.Println("Edit: ")
	err := tmpl.Execute(os.Stdout, article)
	return ID, err
}

// actualClient :
type actualClient struct {
	Config ClientConfig
	Client *http.Client
	Auth   func(*http.Request) error
}

// Link :
type Link struct {
	Rel  string `xml:"rel,attr,omitempty"`
	Href string `xml:"href,attr"`
}

// Entry :
type Entry struct {
	Title   string    `xml:"title"`
	Id      string    `xml:"id"`
	Link    []Link    `xml:"link"`
	Updated time.Time `xml:"updated"`
}

// Create :
func (c *actualClient) Create(article article.Article) (string, error) {
	url := fmt.Sprintf("https://blog.hatena.ne.jp/%s/%s/atom/entry", c.Config.HatenaID, c.Config.BlogID)
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, article)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/xml; charset=utf-8")
	if err := c.Auth(req); err != nil {
		return "", err
	}
	resp, err := c.Client.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		io.Copy(os.Stdout, resp.Body)
		return "", errors.Errorf("something wrong: %s", resp.Status)
	}
	fmt.Println(resp.Status)
	decoder := xml.NewDecoder(resp.Body)
	var entry Entry
	err = decoder.Decode(&entry)
	if err != nil {
		io.Copy(os.Stdout, resp.Body)
		return "", errors.Wrap(err, "parse error")
	}
	for _, link := range entry.Link {
		if link.Rel == "edit" {
			return link.Href, err
		}
	}
	return "", errors.Errorf("hmm.")
}

// Edit :
func (c *actualClient) Edit(article article.Article, url string) (string, error) {
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, article)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("PUT", url, &buf)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/xml; charset=utf-8")
	if err := c.Auth(req); err != nil {
		return "", err
	}
	resp, err := c.Client.Do(req)
	defer resp.Body.Close()
	fmt.Println(resp.Status)
	if resp.StatusCode != http.StatusOK {
		io.Copy(os.Stdout, resp.Body)
		return "", errors.Errorf("something wrong: %s", resp.Status)
	}
	return url, err
}

var tmpl *template.Template

func init() {
	body := `<?xml version="1.0" encoding="utf-8"?>
<entry xmlns="http://www.w3.org/2005/Atom"
       xmlns:app="http://www.w3.org/2007/app">
  <title>{{.Title.Title}}</title>
  <author><name>name</name></author>
  <content type="text/plain">
{{.Body}}
  </content>
{{range .Title.Categories}}
  <category term="{{.}}" />
{{end}}
  <app:control>
    <app:draft>no</app:draft>
  </app:control>
</entry>`
	tmpl = template.Must(template.New("req.tmpl").Parse(body))
}
