package internal

import (
	"chassit-on-repeat/internal/utils"
	"encoding/csv"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog/log"
	"os"
	"sync"
)

type Overrides struct {
	watcher     *fsnotify.Watcher
	stopChannel chan bool

	overrides map[string]string
	mu        sync.RWMutex
}

func NewOverridesHandler() (*Overrides, func()) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal().Str("tag", "overrides").Err(err).Send()
	}

	channel := make(chan bool)
	return &Overrides{
		watcher:     watcher,
		stopChannel: channel,
		overrides:   map[string]string{},
	}, func() { channel <- true }
}

func (o *Overrides) GetOverride(id string) string {
	o.mu.RLock()
	defer o.mu.RUnlock()

	if newId, ok := o.overrides[id]; ok {
		return newId
	}
	return id
}

func (o *Overrides) updateFiles() {
	fileName := fmt.Sprintf("%s/overrides.csv", utils.GetStringEnv("CONFIG_PATH", "/config"))
	f, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer f.Close()

	reader := csv.NewReader(f)
	all, err := reader.ReadAll()
	if err != nil {
		log.Error().Str("tag", "overrides").Err(err).Msg("Error parsing CSV file")
		return
	}
	o.mu.Lock()
	defer o.mu.Unlock()
	o.overrides = map[string]string{}

	for _, row := range all {
		if len(row) == 2 {
			o.overrides[row[0]] = row[1]
		}
	}
	log.Debug().Str("tag", "overrides").Msgf("Updated overrides with %d overrides", len(o.overrides))
}

func (o *Overrides) Watch() {
	defer o.watcher.Close()

	go func() {
		for {
			select {
			case _, ok := <-o.watcher.Events:
				if !ok {
					return
				}
				o.updateFiles()
			case err, ok := <-o.watcher.Errors:
				if !ok {
					return
				}
				log.Error().Str("tag", "overrides").Err(err).Msg("Error watching files")
			}
		}
	}()

	err := o.watcher.Add(utils.GetStringEnv("CONFIG_PATH", "/config"))
	if err != nil {
		log.Fatal().Str("tag", "overrides").Err(err).Msg("Error adding file path")
	}

	o.updateFiles()

	// Block until program exists
	<-o.stopChannel
}
