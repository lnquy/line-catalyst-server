package bot

import (
	"math/rand"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/lnquy/line-catalyst-server/internal/pkg/joke"
)

var (
	jokers = make([]joke.Joker, 1)
)

func (c *Catalyst) initJokers() error {
	pwd, _ := os.Getwd()
	lj, err := joke.NewLocalJSON(filepath.Join(pwd, c.conf.Joke.Folder, "jokes.json"))
	if err != nil {
		return errors.Wrapf(err, "failed to init local JSON joker")
	}
	jokers[0] = lj
	return nil
}

func (c *Catalyst) joke(cmdArgs []string, replyTo string) error {
	idx := rand.Intn(len(jokers))

	jokeStr, err := jokers[idx].Get()
	if err != nil {
		return errors.Wrapf(err, "failed to get joke")
	}
	c.replyTo(replyTo, jokeStr)
	return nil
}
