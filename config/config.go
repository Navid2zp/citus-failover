package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// ServiceConfig represents the config data for the application.
type ServiceConfig struct {
	Monitor struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
	}
	Coordinators []struct {
		Formation string
		Username  string
		Password  string
		DBName    string
	}
	Settings struct {
		CheckInterval int
		Debug         bool
	}
	API struct {
		Enabled bool
		Port    string
		Secret  string
	}
}

var Config ServiceConfig

// getConfigNameAndPath returns the path and name of the config file.
func getConfigNameAndPath(configPath string) (string, string) {
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	name := "config"
	if err != nil {
		panic(err)
	}
	sp := strings.Split(configPath, "/")
	name = strings.Split(sp[len(sp)-1], ".")[0]
	if strings.HasPrefix(configPath, "/") {
		path = strings.Join(sp[0:len(sp)-1], "/")
	}
	return name, path
}

// InitConfig initializes the application config.
func InitConfig(configPath string) (*ServiceConfig, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return &Config, fmt.Errorf("Fatal error getting directory: %s \n", err)
	}
	configName, configPath := getConfigNameAndPath(configPath)
	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)
	viper.AddConfigPath(".")
	viper.AddConfigPath(dir)
	err = viper.ReadInConfig()
	if err != nil {
		return &Config, fmt.Errorf("Fatal error config file: %s \n", err)
	}

	err = viper.Unmarshal(&Config)
	if err != nil {
		return &Config, err
	}
	log.Println("Settings initialized ...")
	return &Config, err
}
