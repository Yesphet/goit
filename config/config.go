package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

var Global *Config

type Config struct {
	Commit Commit
}

type Commit struct {
	Scopes              []string
	Types               []string
	DisableDefaultTypes bool `yaml:"disable_default_types"`
}

func init() {
	Global = &Config{
		Commit: Commit{
			Scopes:              make([]string, 0),
			Types:               make([]string, 0),
			DisableDefaultTypes: false,
		},
	}
	path := "goit.yaml"

	file, err := os.Open(path)
	if err == os.ErrNotExist {
		file, err = os.Open("goit.yml")
	}
	if err != nil {
		return
	}
	defer file.Close()

	raw, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(raw, Global)
	if err != nil {
		fmt.Printf("Read config %s failed, %s.\n", path, err.Error())
	}
}
