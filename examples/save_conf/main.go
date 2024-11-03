package main

import (
	_ "embed"
	"log"

	"github.com/flyfy1/confdir"
)

type config struct {
	StayHappy bool   `yaml:"stay_happy"`
	Path      string `yaml:"path"`
}

func main() {
	loader := confdir.NewYamlLoader(".tools", confdir.WithLogFunc(log.Printf))

	exampleCfg := &config{
		StayHappy: true,
		Path:      "/var/etc/sshd",
	}
	loader.RegisterFile("example.yml", exampleCfg)
	panicOnError(loader.LoadAll())

	panicOnError(loader.SaveConfig("example.yml"))
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
