package article

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseTitle(t *testing.T) {
	t.Run("no category", func(t *testing.T) {
		line := "# yamlの何か"
		title := parseTitle(line)
		assert.Exactly(t, "yamlの何か", title.Title)
		assert.EqualValues(t, []string{}, title.Categories)
	})
	t.Run("categories", func(t *testing.T) {
		line := "# [python][yaml]yamlの何か"
		title := parseTitle(line)
		assert.Exactly(t, "yamlの何か", title.Title)
		assert.EqualValues(t, []string{"python", "yaml"}, title.Categories)
	})
}

func TestParseArticle(t *testing.T) {
	content := `#[tag1][tag2]タイトル

なんか文章

## トップレベルのセクションはひとつだけ

はい

## 見出しにしたい場合はこのレベルにする

はい`
	article, err := ParseArticle(content)
	require.NoError(t, err)
	assert.Exactly(t, "タイトル", article.Title.Title)
	assert.Exactly(t, []string{"tag1", "tag2"}, article.Title.Categories)
	assert.Exactly(t, content, article.Body)
}
