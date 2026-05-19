package main

import (
	// "flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Mode   string   `yaml:"mode"`
	ID     string	`yaml:"id,omitempty"`
	Port   int      `yaml:"port,omitempty"`
	Addrs  []string `yaml:"addrs"`
}

var config *Config

func configInit() error {
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, ".config", "alive", "config.yaml")

	if err := configLoad(path); err != nil {
		return err
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
		} else
		if Get().Mode == "client" {
			config.Port = 9090
		}
	}

	if err := configRewrite(path); err != nil {
		return err
	}
	return nil
}

func configLoad(path string) error {
	result, err := os.ReadFile(path)// Go 1.16+
	if os.IsNotExist(err) {
		// 文件不存在，创建默认配置
		config = &Config{Mode: "server"}
		os.MkdirAll(filepath.Dir(path), 0755)
		return configRewrite(path)
	}

	return yaml.Unmarshal(result, &config)
}

func configRewrite(path string) error {
	data, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return err
	}
	return nil
}

func Get() *Config {
	return config
}
func main() {
	if err := configInit(); err != nil {
		fmt.Println("Failed to load configuration")
		panic(err)
	}

	if Get().Mode == "server" {
		err := startServer(":"+strconv.Itoa(Get().Port))
		if err != nil {
			log.Println(err)
		}
	}

	if Get().Mode == "client" && Get().Addrs != nil {
		startClients(Get().Addrs)
		log.Printf("客户端已启动")
	} else {
		log.Printf("未提供Addrs")
		log.Printf("客户端自动退出")
		return
	}
	select {}
}
