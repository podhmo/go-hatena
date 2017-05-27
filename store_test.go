package hatena

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
	"time"
)

func TestSaveCommit(t *testing.T) {
	r := bytes.NewBufferString(strings.Trim(`
https://blog.hatena.ne.jp/{はてなID}/{ブログID}/atom/entry/{entry_id}@head@Thu Aug 18 12:33:45 +0000 2016@create
`, "\n"))
	now, err := time.Parse(time.RubyDate, "Fri Aug 19 04:50:21 +0000 2016")
	if err != nil {
		t.Errorf("time parse error: %s", err)
	}
	commit := Commit{
		ID:        "https://blog.hatena.ne.jp/{はてなID}/{ブログID}/atom/entry/{entry_id}",
		CreatedAt: now,
		Alias:     "newItem",
		Action:    "create",
	}
	w := &bytes.Buffer{}

	err = saveCommit(w, r, commit)
	if err != nil {
		t.Errorf("saveCommit failured: %s", err)
	}

	text, err := ioutil.ReadAll(w)
	if err != nil {
		t.Errorf("invalid contents: %s", err)
	}
	result := strings.Split(string(text), "\n")[0]
	expected := "https://blog.hatena.ne.jp/{はてなID}/{ブログID}/atom/entry/{entry_id}@newItem@Fri Aug 19 04:50:21 +0000 2016@create"
	if result != expected {
		t.Errorf("commit line is must be %q but %q", expected, result)
	}
}

func TestLoadCommit(t *testing.T) {
	r := bytes.NewBufferString(strings.Trim(`
https://blog.hatena.ne.jp/{はてなID}/{ブログID}/atom/entry/{entry_id2}@head@Fri Aug 19 04:50:21 +0000 2016@update
https://blog.hatena.ne.jp/{はてなID}/{ブログID}/atom/entry/{entry_id}@head@Thu Aug 18 12:33:45 +0000 2016@create
`, "\n"))

	cases := []struct {
		alias      string
		expectedID string
		found      bool
	}{
		{alias: "head", expectedID: "https://blog.hatena.ne.jp/{はてなID}/{ブログID}/atom/entry/{entry_id2}", found: true},
		{alias: "hmm", found: false},
	}

	for _, c := range cases {
		commit, err := loadCommit(r, c.alias)
		if err != nil {
			t.Error(err)
		}
		if c.found {
			if commit == nil {
				t.Errorf("%q should be found. but not found", c.alias)
			}
			if commit.ID != c.expectedID {
				t.Errorf("expected id is %q but found id is %q", c.expectedID, commit.ID)
			}
		}
	}
}
