package g

import (
	"encoding/json"
	"log"
	"mp"
	"sync"

	"github.com/toolkits/file"
)

var (
	ConfigFile string
	config     *message.GlobalConfig
	configLock = new(sync.RWMutex)
)

func Config() *message.GlobalConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

// ParseConfig 解析配置文件
func ParseConfig(cfg string) {
	ConfigLock := new(sync.RWMutex)
	if cfg == "" {
		log.Fatalln("config file not specified: use -c $filename")
	}
	if !file.IsExist(cfg) {
		log.Fatalln("config file specified not found:", cfg)
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file", cfg, "error:", err.Error())
	}

	var c message.GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("parse config file", cfg, "error:", err.Error())
	}
	// set config
	ConfigLock.Lock()
	defer ConfigLock.Unlock()
	config = &c
	log.Println("g.ParseConfig ok, file", cfg)
}
