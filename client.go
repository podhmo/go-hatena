package hatena

import (
	"fmt"

	"github.com/k0kubun/pp"
	"github.com/podhmo/hatena/article"
)

// Client :
type Client interface {
	Create(article article.Article) (string, error)
	Edit(article article.Article, ID string) (string, error)
}

// DummyClient :
type DummyClient struct {
}

// Create :
func (c *DummyClient) Create(article article.Article) (string, error) {
	fmt.Println("Create: ")
	pp.Println(article)
	return "xxx", nil
}

// Edit :
func (c *DummyClient) Edit(article article.Article, ID string) (string, error) {
	fmt.Println("Edit: ")
	pp.Println(article)
	return ID, nil
}
