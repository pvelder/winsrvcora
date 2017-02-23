package oracle

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type tomlConfig struct {
	DB database   `toml:"database"`
	FS filesystem `toml:"filesystem"`
}

type database struct {
	Type          string
	Connectstring string
}

type filesystem struct {
	watch string
}

var config tomlConfig

func WatchOracleArchiveLogs() {

	if _, err := toml.DecodeFile("D:\\config.toml", &config); err != nil {
		fmt.Println(err)
		return
	}

}
