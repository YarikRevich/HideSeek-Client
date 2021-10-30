package resource_manager

import (
	"embed"
	_ "image/png"
	"io/fs"
	"strings"
	"sync"

	loaderpostprocessor "github.com/YarikRevich/HideSeek-Client/internal/resource_manager/loader_post_processor"
	metadatapostprocessor "github.com/YarikRevich/HideSeek-Client/internal/resource_manager/loader_post_processor/metadata"
	"github.com/sirupsen/logrus"
)

func processResource(e embed.FS, sourcePath string, files []fs.DirEntry, motherWg *sync.WaitGroup, loaders ...Loader) {
	for _, v := range files {
		path := sourcePath + "/" + v.Name()
		if v.IsDir() {
			processResourceDir(e, path, motherWg, loaders...)
		} else {
			nameSplit := strings.Split(v.Name(), ".")
			extension := nameSplit[len(nameSplit)-1]

			for _, l := range loaders {
				l(e, extension, path, motherWg)
			}
		}
	}
}

func processResourceDir(e embed.FS, path string, wg *sync.WaitGroup, loaders ...Loader) {
	d, err := e.ReadDir(path)
	if err != nil {
		logrus.Fatal("error happened reading dir from embedded fs", err)
	}
	processResource(e, path, d, wg, loaders...)
}

func LoadResources(loaders map[Component][]Loader) {
	var wg sync.WaitGroup
	for c, l := range loaders {
		for _, p := range c.SeparatePath() {
			wg.Add(1)
			go func(c Component, p string, l []Loader) {
				defer wg.Done()
				processResourceDir(c.Embed, p, &wg, l...)

			}(c, p, l)
		}
	}
	wg.Wait()

	loaderpostprocessor.ApplyPostProcessors(
		metadatapostprocessor.ConnectAdditionalStatementsToMetadata)
}
