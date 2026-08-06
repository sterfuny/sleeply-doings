package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sleeply-alive/internal/ctrl"

	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Mode  string   `yaml:"mode"`
	ID    string   `yaml:"id,omitempty"`
	Port  int      `yaml:"port,omitempty"`
	Addrs []string `yaml:"addrs"`
}

var config *Config

func Get() Config {
	return *config
}

var home, _ = os.UserHomeDir()
var Path = filepath.Join(home, ".config", "alive")

func Init(path string) {
	ctrl.OpenDB(Path)
	ctrl.ReadDB()
	if path == "" {
		path = filepath.Join(Path, "config.yaml")
	}

	if err := Load(path); err != nil {
		fmt.Println("Failed to load configuration")
		panic(err)
	}

	if Get().ID == "" {
		config.ID = uuid.New().String()
	}
	if Get().Mode == "" {
		config.Mode = "server"
	}
	if Get().Port == 0 {
		if Get().Mode == "server" {
			config.Port = 9180
		} else if Get().Mode == "client" {
			config.Port = 9090
		}
	}

	if err := Rewrite(path); err != nil {
		fmt.Println("Failed to load configuration")
		panic(err)
	}
}

func Load(path string) error {
	result, err := os.ReadFile(path) // Go 1.16+
	if os.IsNotExist(err) {
		// 文件不存在，创建默认配置
		config = &Config{Mode: "server"}
		os.MkdirAll(filepath.Dir(path), 0755)
		return Rewrite(path)
	}

	return yaml.Unmarshal(result, &config)
}

func Rewrite(path string) error {
	data, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return err
	}
	return nil
}
