package store

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type SegmentConfig struct {
	MaxSize int `yaml:"maxSize"`
}

type Config struct {
	Segments SegmentConfig `yaml:"segments"`
}

func NewConfig(yamlConfig []byte) Config {
	c := Config{}
	err := yaml.Unmarshal(yamlConfig, &c)
	if err != nil {
		fmt.Println(err)
	}

	return c
}
