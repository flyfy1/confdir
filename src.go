package confdir

import (
	"errors"
	"os"
	"path"

	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type logFunc func(format string, arguments ...any)

type YamlLoader struct {
	path        string
	logFunc     logFunc
	configFiles map[string]*cfgLoader
}

type OptionFunc func(dir *YamlLoader)

func WithLogFunc(logger logFunc) OptionFunc {
	return func(c *YamlLoader) {
		c.logFunc = logger
	}
}

// NewYamlLoader creates YamlLoader with configFolder relative path to $Home.
// If configFolder is empty, the config file would be put under $HOME directly
func NewYamlLoader(configFolder string, options ...OptionFunc) *YamlLoader {
	c := &YamlLoader{
		path:        configFolder,
		configFiles: map[string]*cfgLoader{},
	}

	for _, option := range options {
		option(c)
	}

	return c
}

// LoadAll loads all configs registered
func (loader *YamlLoader) LoadAll() error {
	configFolder := pathFromHome(loader.path)

	if loader.path != "" {
		err := loader.prepareFolder(configFolder)
		if err != nil {
			loader.logIfNeeded("prepareFolder error: ", err)
			return err
		}
	}

	for fileName, l := range loader.configFiles {
		configFile := path.Join(configFolder, fileName)
		loader.logIfNeeded("loading file: %s", configFile)

		_, err := os.Stat(configFile)
		if errors.Is(err, os.ErrNotExist) {
			if l.exampleContent != "" {
				loader.logIfNeeded("creating example config: %s", configFile)
				err = saveExample(configFile, l.exampleContent)
				if err != nil {
					return err
				}
			} else if l.errorOnNoFile {
				loader.logIfNeeded("returnning error because expected")
				return os.ErrNotExist
			} else {
				// we just skip loading
				continue
			}
		}

		err = loadYaml(configFile, l.config)
		if err != nil {
			loader.logIfNeeded("unexpected error when loading %s: %s", fileName, err)
			return err
		}
	}

	return nil
}

// PrepareFolder prepares folderPath, to make sure folder exist
// it set permission to 700 (i.e., only owner have full access to folder)
func (loader *YamlLoader) prepareFolder(folder string) error {
	return os.MkdirAll(folder, 0700)
}

func (loader *YamlLoader) logIfNeeded(format string, args ...interface{}) {
	if loader.logFunc == nil {
		return
	}

	loader.logFunc(format, args...)
}

type cfgLoader struct {
	config         interface{}
	exampleContent string
	errorOnNoFile  bool
}

type RegisterOptionFunc func(loader *cfgLoader)

func RegWithExampleContent(content string) RegisterOptionFunc {
	return func(loader *cfgLoader) {
		loader.exampleContent = content
	}
}

// RegErrorOnNoFile if registered file not exist, and no exampleContent provided, return error (i.e., loading failed)
func RegErrorOnNoFile() RegisterOptionFunc {
	return func(loader *cfgLoader) {
		loader.errorOnNoFile = true
	}
}

func (loader *YamlLoader) RegisterFile(filename string, cfg interface{}, options ...RegisterOptionFunc) {
	c := &cfgLoader{config: cfg}
	for _, o := range options {
		o(c)
	}
	loader.configFiles[filename] = c
}

// loadYaml loads `config` from given filepath, in yaml format
// expect `filepath` either absolute path, or relative to the directory executing
// to be relative to HomeDir, wrap with `pathFromHome(filepath)`
func loadYaml(filepath string, config interface{}) error {
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(content, config)
}

// pathFromHome appends the home dir, in front of the given relativePath
func pathFromHome(relativePath string) string {
	home := os.Getenv("HOME")
	return path.Join(home, relativePath)
}

// saveExample persist example data
func saveExample(path string, data string) error {
	return os.WriteFile(path, []byte(data), 0700)
}
