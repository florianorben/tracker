package tracker

import (
	"bytes"
	"fmt"
	"github.com/florianorben/tracker/helpers"
	"io/ioutil"
	"os"
	"strings"
	"unicode"

	"github.com/BurntSushi/toml"
	"github.com/spf13/viper"
)

type Config struct {
	Core struct {
		Editor string
	}

	Colors struct {
		Enabled bool
	}

	Log struct {
		DefaultStartDate int
		DefaultEndDate   int
	}

	Backend struct {
		Token          string
		Url            string
		User           string
		AutoAddWorkLog bool
	}
}

var DefaultConfig = []byte(`
[core]
editor = "vim"

[colors]
enabled = true

[log]
defaultStartDate = -14
defaultEndDate = 0

[backend]
token = ""
url = ""
user = ""
autoAddWorkLog = false
`)

func SetConfig(key string, value interface{}) {

	keyArgs := strings.Split(key, ".")
	if len(keyArgs) != 2 {
		fmt.Printf("Error: %s\n", helpers.PrintRed("invalid configuration key"))
		os.Exit(1)
	}

	current := viper.GetStringMap(keyArgs[0])
	subKey := []rune(keyArgs[1])
	subKey[0] = unicode.ToUpper(subKey[0])
	current[string(subKey)] = value
	viper.Set(keyArgs[0], current)
}

func WriteConfigToFile() error {
	var c Config

	err := viper.Unmarshal(&c)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(c); err != nil {
		return err
	}

	return writeToConfigFile(buf.Bytes())
}

func EditConfig() error {
	b, err := ioutil.ReadFile(viper.ConfigFileUsed())
	if err != nil {
		return err
	}

	c, err := helpers.OpenInEditor(b)
	if err != nil {
		return err
	}

	return writeToConfigFile(c)
}

func writeToConfigFile(b []byte) error {
	if file, err := os.OpenFile(viper.ConfigFileUsed(), os.O_WRONLY|os.O_TRUNC, 0666); err == nil {
		defer file.Close()

		_, err := file.Write(b)
		if err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}
