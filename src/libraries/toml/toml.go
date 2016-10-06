package toml

import (
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

type TomlConfig struct {
	Etcd          DBConfig
	Pillx         PillConfig
	//Elasticsearch DBConfig
	//DB       database `toml:"database"`
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBname   string
}

type PillConfig struct {
	GatewayOuterHost string
	GatewayOuterPort int
	GatewayInnerHost string
	GatewayInnerPort int
	WorkerInnerHost  string
	WorkerInnerPort  int
	GatewayName      string
	WorkerName       string
}

func LoadTomlConfig(filename string) (tc TomlConfig, err error) {
	tomlData, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Read failed", err)
		return
	}

	if _, err1 := toml.Decode(string(tomlData), &tc); err1 != nil {
		err = err1
		fmt.Println("ReadToml failed", err)
		return
	}
	return
}
