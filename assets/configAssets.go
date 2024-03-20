package assets

import (
	_ "embed"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

//go:embed configs/names.yaml
var namesYaml []byte

var (
	NameConfig nameConfig
)

func initializeConfigs(useEmbededAssets bool) {
	var namesFile []byte = namesYaml
	var err error

	if !useEmbededAssets {
		namesFile, err = os.ReadFile("./assets/configs/names.yaml")
		if err != nil {
			fmt.Println("Couldn't read the names file:", err)

			return
		}
	}

	err = yaml.Unmarshal(namesFile, &NameConfig)
	if err != nil {
		fmt.Println("Couldn't (...) yaml file:", err)

		return
	}
}

type nameConfig struct {
	Firstnames []string `yaml:"firstnames"`
	Lastnames  []string `yaml:"lastnames"`
}
