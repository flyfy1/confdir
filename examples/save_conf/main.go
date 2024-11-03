package main

import (
	"log"

	"github.com/your-project/confdir"
)

type secret struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

//go:embed secret_config.yml
var exampleSecretContent string

func main() {
	loader := confdir.NewYamlLoader(".tools", confdir.WithLogFunc(log.Printf))

	secretCfg := &secret{}
	loader.RegisterFile(".secret.yml", secretCfg, confdir.RegWithExampleContent(exampleSecretContent))

	err := loader.LoadAll()
	if err != nil {
		log.Fatal(err)
	}

	// 修改配置
	secretCfg.Username = "newuser"
	secretCfg.Password = "newpassword"

	// 保存修改后的配置
	err = loader.SaveConfig(".secret.yml")
	if err != nil {
		log.Fatal(err)
	}
}
