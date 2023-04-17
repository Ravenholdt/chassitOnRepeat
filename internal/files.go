package internal

import (
	"chassit-on-repeat/internal/utils"
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog/log"
	"net/url"
	"path/filepath"
	"regexp"
	"sync"
)

type VideoFile struct {
	Name string
	Id   string
	Url  string
}

type FileHandler struct {
	watcher     *fsnotify.Watcher
	stopChannel chan bool
	mu          sync.RWMutex
	videos      map[string]VideoFile
}

var videoNameRegex = regexp.MustCompile("^/+files/(.*)-([A-Za-z0-9_-]{11}).mp4$")

func NewFileHandler() (*FileHandler, func()) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal().Str("tag", "files").Err(err).Send()
	}

	channel := make(chan bool)
	return &FileHandler{
		watcher:     watcher,
		stopChannel: channel,
		videos:      map[string]VideoFile{},
	}, func() { channel <- true }
}

func (f *FileHandler) GetVideos() []VideoFile {
	f.mu.RLock()
	defer f.mu.RUnlock()

	var videos []VideoFile
	for _, v := range f.videos {
		videos = append(videos, v)
	}
	return videos
}

func (f *FileHandler) GetVideoIds() []string {
	videos := f.GetVideos()
	if len(videos) <= 0 {
		return []string{}
	}
	var ids []string
	for _, file := range videos {
		ids = append(ids, file.Id)
	}
	return ids
}

// GetVideoFile Tries to get a video by id, error if no video was found
func (f *FileHandler) GetVideoFile(id string) (*VideoFile, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if val, ok := f.videos[id]; ok {
		return &val, nil
	}
	return nil, errors.New("no video found with id")
}

// Updates the internal list of videos to lower
// the amount of file requests
func (f *FileHandler) updateFiles() {
	glob, err := filepath.Glob(fmt.Sprintf("%s/*.mp4", utils.GetStringEnv("FILES_PATH", "/files")))
	if err != nil {
		return
	}

	log.Debug().Str("tag", "files").Msg("Updating video files")
	f.mu.Lock()
	defer f.mu.Unlock()
	f.videos = map[string]VideoFile{}
	for _, s := range glob {
		parts := videoNameRegex.FindAllStringSubmatch(s, -1)
		if len(parts[0]) == 3 {
			name := parts[0][1]
			id := parts[0][2]
			log.Debug().Str("tag", "files").Str("file", s).Str("name", name).Str("id", id).Msg("Found valid file")
			f.videos[id] = VideoFile{
				Name: name,
				Id:   id,
				Url:  fmt.Sprintf("/files/%s-%s.mp4", url.PathEscape(name), id),
			}
		}
	}
}

func (f *FileHandler) Watch() {
	defer f.watcher.Close()

	go func() {
		for {
			select {
			case _, ok := <-f.watcher.Events:
				if !ok {
					return
				}
				f.updateFiles()
			case err, ok := <-f.watcher.Errors:
				if !ok {
					return
				}
				log.Error().Str("tag", "files").Err(err).Msg("Error watching files")
			}
		}
	}()

	err := f.watcher.Add(utils.GetStringEnv("FILES_PATH", "/files"))
	if err != nil {
		log.Fatal().Str("tag", "files").Err(err).Msg("Error adding file path")
	}

	f.updateFiles()

	// Block until program exists
	<-f.stopChannel
}
