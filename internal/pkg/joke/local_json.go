package joke

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type (
	localJSON struct {
		jokes   map[int]*localJoke
		mux     *sync.RWMutex
		randSrc *rand.Rand
	}

	localJoke struct {
		ID        int    `json:"id"`
		Type      string `json:"type"`
		Setup     string `json:"setup"`
		Punchline string `json:"punchline"`
	}
)

func NewLocalJSON(path string) (Joker, error) {
	// pwd, _ := os.Getwd()
	// files, err := ioutil.ReadDir(filepath.Join(pwd, path))
	// if err != nil {
	// 	return nil, errors.Wrapf(err, "failed to read JSON jokes from local directory")
	// }
	//
	// for _, f := range files {
	// 	f.IsDir()
	// 	fmt.Println(f.Name())
	// }

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read JSON jokes from local directory")
	}
	var jokes []*localJoke
	if err := json.Unmarshal(b, &jokes); err != nil {
		return nil, errors.Wrapf(err, "failed to decode JSON jokes")
	}

	var lj = &localJSON{
		mux:     &sync.RWMutex{},
		jokes:   make(map[int]*localJoke, len(jokes)),
		randSrc: rand.New(rand.NewSource(time.Now().Unix())),
	}
	lj.mux.Lock()
	for i, joke := range jokes {
		lj.jokes[i] = joke
	}
	lj.mux.Unlock()
	return lj, nil
}

func (l *localJSON) Get() (string, error) {
	l.mux.RLock()
	idx := l.randSrc.Intn(len(l.jokes))
	joke := l.jokes[idx]
	l.mux.RUnlock()

	return strings.TrimSpace(joke.Setup + "\n" + joke.Punchline), nil
}
