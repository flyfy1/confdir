# ConfDir

ConfDir provides an easy interface, to store conf into your home dir.

## How to Use

Can see this example file: (the file also at `examples/load_conf/main.go`)

```go
package main

import (
	"log"
	
	_ "embed"

	"github.com/flyfy1/confdir"
)

// secret is your config object
type secret struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// config is another kind of your config object 
type config struct {
	StayHappy bool   `yaml:"stay_happy"`
	Path      string `yaml:"path"`
}

// you can prepare an example `secret_config.yml`, and use `go:embed` to assign its content to a variable
//go:embed secret_config.yml
var exampleSecretContent string

func main() {
	// use `~/.tools` as the config folder; files registered would be put under this config folder
	loader := confdir.NewYamlLoader(".tools", confdir.WithLogFunc(log.Printf))

	// register config object with its corresponding file; use `config.RegWithExampleContent` to define the content of 
	// the example config when the file doesn't exist
	secretCfg := &secret{}
	loader.RegisterFile(".secret.yml", secretCfg, confdir.RegWithExampleContent(exampleSecretContent))

	// you can assign default configs. If `confdir.RegWithExampleContent` isn't passed into RegisterFile, it simply 
	// ignores it (unless, `confdir.RegErrorOnNoFile` is set, then `loader.LoadAll()` would return error instead)
	normalCfg := &config{
		StayHappy: true,
		Path:      "/var/etc/sshd",
	}
	loader.RegisterFile("config.yml", normalCfg)

	// loader.LoadAll() would: 1. check if folder exist in home directory; 2. load config files defined
	err := loader.LoadAll()

	// if `confdir.RegErrorOnNoFile` is set, and no ExampleContent provided, err would be `os.ErrNoFile` when 
	// ConfigFile not exist
	log.Println("load err: ", err)      // load err: nil
	
	log.Println("secret config loaded: ", secretCfg)    // secret config loaded:  &{example }
	log.Println("normal config not touched: ", normalCfg)   // normal config not touched:  &{true /var/etc/sshd}
}

```