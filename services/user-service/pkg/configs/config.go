package configs

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var Config *viper.Viper

func InitViper() {
	config := viper.New()
	wd, err := os.Getwd()
	if err != nil {
		log.Println("error in os Get wd")
	}
	rootDir := findRoot(wd)
	config.AddConfigPath(rootDir)
	config.SetConfigName("config")
	config.SetConfigType("yaml")
	config.AutomaticEnv()
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	err = config.ReadInConfig()
	Config = config
}

func findRoot(path string) string {
	dir := path
	for {
		if _, err := os.Stat(filepath.Join(dir, "config.yaml")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}
