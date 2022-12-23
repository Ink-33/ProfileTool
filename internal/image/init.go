package image

import (
	_ "embed"
	"errors"
	"io/fs"
	"math"
	"mime"
	"os"
	"path/filepath"
	"sync"

	"github.com/Ink-33/ProfileTool/config"
	log "github.com/Ink-33/logger"
)

//go:embed blank.png
var blank []byte

// Images ...
type Images struct {
	lock    sync.RWMutex
	imgs    []string
	current struct {
		file []byte
		name string
	}
	counter struct {
		max uint32
		now uint32
	}
}

// Init images
func Init(conf *config.Config) (images *Images, err error) {
	assets, err := os.Open(conf.Assets)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}
		wd, _ := os.Getwd()
		err := os.Mkdir(filepath.Join(wd, "assets"), 0o644)
		if err != nil {
			return nil, err
		}
		log.Warn("Assets directory doesn't exist. Running with blank png.")
		return &Images{
			current: struct {
				file []byte
				name string
			}{
				file: blank,
				name: "blank.png",
			},
		}, nil
	}
	files, err := assets.Readdir(0)
	if err != nil {
		return nil, err
	}
	imgs := make([]string, 0)
	for k := range files {
		if files[k].IsDir() {
			continue
		}
		ext := filepath.Ext(files[k].Name())
		fmime := mime.TypeByExtension(ext)
		ok := false
		switch fmime {
		case "image/bmp", "image/gif", "image/jpeg", "image/png", "image/webp":
			ok = true
		}
		if !ok {
			continue
		}
		imgs = append(imgs, filepath.Join(conf.Assets, files[k].Name()))
	}
	images = &Images{
		imgs: imgs,
		counter: struct {
			max uint32
			now uint32
		}{
			max: func() uint32 {
				if conf.SwitchCunter < 0 {
					return 0
				}
				if conf.SwitchCunter > math.MaxUint32 {
					return math.MaxUint32
				}
				return uint32(conf.SwitchCunter)
			}(),
			now: 0,
		},
	}
	if images.Length() == 0 {
		images.current = struct {
			file []byte
			name string
		}{file: blank, name: "blank"}
	}

	return images, nil
}

// Length returns the length of valid images.
func (i *Images) Length() int {
	return len(i.imgs)
}
