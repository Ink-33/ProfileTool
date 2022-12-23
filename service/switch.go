package service

import (
	"net/http"

	"github.com/Ink-33/ProfileTool/internal/image"
	log "github.com/Ink-33/logger"
)

func Switch(storage *image.Images) func(w http.ResponseWriter, r *http.Request) {
	f := func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusOK
		if r.Method != http.MethodGet {
			code = http.StatusBadRequest
			return
		}

		w.Header().Add("cache-control", "no-cache,max-age=0")
		redirect := r.URL.Query().Get("redirect")
		if redirect != "" {
			code = http.StatusFound
			http.Redirect(w, r, redirect, http.StatusFound)
		}

		err := storage.Switch()
		if err != nil {
			log.Error(err.Error())
		}

		defer func() {
			log.Info("%v %v %v %v", r.RemoteAddr, r.Method, r.URL.Path, code)
			if code != http.StatusOK && code != http.StatusFound {
				w.WriteHeader(code)
			}
			r.Body.Close()
		}()
	}
	return f
}
