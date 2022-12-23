package image

import (
	"math/rand"
	"os"
	"sync/atomic"
	"time"

	log "github.com/Ink-33/logger"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (i *Images) Get() (b []byte, name string) {
	if len(i.current.file) == 0 {
		i.Switch()
	}
	if i.counter.max == 0 {
		defer i.Switch()
	} else {
		atomic.AddUint32(&i.counter.now, 1)
		if atomic.CompareAndSwapUint32(&i.counter.now, i.counter.max, 0) {
			defer i.Switch()
		}
	}
	i.lock.RLock()
	f, n := i.current.file, i.current.name
	i.lock.RUnlock()
	return f, n
}

func (i *Images) Switch() error {
	if i.Length() == 0 {
		log.Warn("No image!")
		return nil
	}
	filename := ""
	for {
		filename = i.imgs[rand.Intn(i.Length())]
		if filename != i.current.name {
			break
		}
	}
	file, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	defer func(old, new string) {
		log.Info("Switch image cache: %v -> %v", old, new)
	}(i.current.name, filename)
	i.lock.Lock()
	i.current = struct {
		file []byte
		name string
	}{
		file: file,
		name: filename,
	}
	i.lock.Unlock()
	atomic.StoreUint32(&i.counter.now, 0)
	return nil
}
