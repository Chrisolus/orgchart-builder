package initializers

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func InitViper(path string) {
	// Separate Directory & Filename
	dir, filename := filepath.Dir(path), filepath.Base(path)
	// Extract type/extension of the file
	confType := filepath.Ext(filename)
	// Extract the name of the file
	name := strings.TrimSuffix(filename, confType)
	// Remove . from the confType
	confType = strings.TrimPrefix(confType, ".")

	viper.SetConfigName(name)
	viper.SetConfigType(confType)
	viper.AddConfigPath(dir)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	} else {
		fmt.Println("Initiated", viper.GetString("title"))
	}
}
