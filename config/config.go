package config

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"

	log "github.com/Ink-33/logger"
)

// Parse config
func Parse(path string) (conf *Config, err error) {
	file, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			b, err := json.Marshal(defaultConfig)
			if err != nil {
				return nil, err
			}
			err = os.WriteFile(path, b, 0o644)
			if err != nil {
				return nil, err
			}
			log.Warn("Config file doesn't exist. Running with default config.")
			configcopy := defaultConfig
			return &configcopy, nil
		}
		return nil, err
	}
	conf = &defaultConfig
	err = json.Unmarshal(file, conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

// Config is the structure of config file.
type Config struct {
	Endpoint struct {
		Image  string `json:"image"`
		Switch string `json:"switch"`
	} `json:"endpoint"`
	Listen       string `json:"listen"`
	Assets       string `json:"assets"`
	SwitchCunter int    `json:"switch_counter"`
}

var defaultConfig = Config{
	Endpoint: struct {
		Image  string "json:\"image\""
		Switch string "json:\"switch\""
	}{
		Image:  "/v1/profile/image",
		Switch: "/v1/profile/switch",
	},
	Listen:       "127.0.0.1:9080",
	Assets:       "assets",
	SwitchCunter: 5,
}
