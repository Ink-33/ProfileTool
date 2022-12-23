package service

import (
	"net/http"

	"github.com/Ink-33/ProfileTool/internal/image"
	"github.com/Ink-33/ProfileTool/utils"
	log "github.com/Ink-33/logger"
)

func Image(storage *image.Images) func(w http.ResponseWriter, r *http.Request) {
	f := func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusOK
		if r.Method != http.MethodGet {
			code = http.StatusBadRequest
			return
		}
		w.Header().Add("cache-control", "no-cache,max-age=0")
		file, name := storage.Get()
		_, err := w.Write(file)
		if err != nil {
			log.Error(err.Error())
		}

		defer func() {
			log.Info("%v %v %v %v %v", utils.ReadUserIP(r), r.Method, r.URL.Path, code, name)
			if code != http.StatusOK {
				w.WriteHeader(code)
			}
			r.Body.Close()
		}()
	}
	return f
}
