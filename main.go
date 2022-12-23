package main

import (
	"flag"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Ink-33/ProfileTool/config"
	"github.com/Ink-33/ProfileTool/internal/image"
	"github.com/Ink-33/ProfileTool/service"
	"github.com/Ink-33/ProfileTool/utils"
	log "github.com/Ink-33/logger"
)

func init() {
	log.SetProductName("ProfileTool")
}

var (
	flagConfig string
)

func main() {
	wd, _ := os.Getwd()
	flag.StringVar(&flagConfig, "c", filepath.Join(wd, "config.json"), "Configuration filename")
	flag.Parse()

	log.Info("Reading config: %v", flagConfig)
	conf, err := config.Parse(flagConfig)
	if err != nil {
		log.Error("Invalid config: %v", err.Error())
	}

	imgs, err := image.Init(conf)
	if err != nil {
		log.Fatal("Initing images failed: %v", err.Error())
	}

	log.Info("Loaded %v images.", imgs.Length())

	mux := &http.ServeMux{}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Info("%v %v %v %v", utils.ReadUserIP(r), r.Method, r.URL.Path, http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
	})

	mux.HandleFunc(conf.Endpoint.Image, service.Image(imgs))
	mux.HandleFunc(conf.Endpoint.Switch, service.Switch(imgs))

	listener, err := net.Listen("tcp", conf.Listen)
	if err != nil {
		log.Fatal("Cannot start server: %v", err.Error())
	}
	log.Info("Listening on: %v", listener.Addr().String())

	log.Fatal("%v", http.Serve(listener, mux))
}
