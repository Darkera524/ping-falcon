package main

import (
	"github.com/toolkits/file"
	"encoding/json"
	"sync"
	"time"
)

var (
	config *Config
	lock = new(sync.RWMutex)
)

type Config struct {
	Interval int `json:"interval"`
	Portal_path	string `json:"portal_path"`
	MaxConn int `json:"maxConn"`
	MaxIdle int `json:"maxIdle"`
	Agent_path string `json:"agent_path"`
}

func ParseConfig(cfg string){
	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		Logger().Println(err.Error())
	}

	var c Config
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		Logger().Println(err.Error())
	}

	lock.Lock()
	defer lock.Unlock()

	config = &c
}

func GetConfig() *Config {
	lock.RLock()
	defer lock.RUnlock()
	return config
}

func CronConfig(interval int, cfg string){
	for {
		time.Sleep(time.Duration(interval) * time.Second)
		ParseConfig(cfg)
	}
}


