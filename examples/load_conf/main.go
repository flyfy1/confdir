package main

import (
	_ "embed"
	"github.com/flyfy1/confdir"
	"log"
)

type secret struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type config struct {
	StayHappy bool   `yaml:"stay_happy"`
	Path      string `yaml:"path"`
}

//go:embed secret_config.yml
var exampleSecretContent string

func main() {
	loader := confdir.NewYamlLoader(".tools", confdir.WithLogFunc(log.Printf))

	secretCfg := &secret{}
	loader.RegisterFile(".secret.yml", secretCfg, confdir.RegWithExampleContent(exampleSecretContent))

	normalCfg := &config{
		StayHappy: true,
		Path:      "/var/etc/sshd",
	}
	loader.RegisterFile("config.yml", normalCfg)

	panicOnError(loader.LoadAll())

	log.Println("secret config loaded: ", secretCfg)
	log.Println("normal config not touched: ", normalCfg)
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
