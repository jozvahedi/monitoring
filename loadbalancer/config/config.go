package config

import (
	"encoding/json"
	"os"

	"gopkg.in/yaml.v2"
)

type ConfigYmlFile struct {
	HTTPServer struct {
		HTTPServerServer string `yaml:"httpServerServer"`
		HTTPServerPort   string `yaml:"httpServerPort"`
	} `yaml:"httpServer"`
	BasicAutentication struct {
		BasicAutenticationUsername string `yaml:"basicAutenticationUsername"`
		BasicAutenticationPassword string `yaml:"basicAutenticationPassword"`
	} `yaml:"basicAutentication"`
}

type ConfigJsonFile struct {
	Whitelistip []struct {
		ID          string `json:"id"`
		IP          string `json:"ip"`
		Description string `json:"description"`
	} `json:"whitelistip"`
	Blacklistip []struct {
		ID          string `json:"id"`
		IP          string `json:"ip"`
		Description string `json:"description"`
	} `json:"blacklistip"`
	Middelwarepath []struct {
		Path       string `json:"path"`
		Middelware []struct {
			Name string `json:"name"`
		} `json:"middelware"`
	} `json:"middelwarepath"`
}

var JsonConfigFile ConfigJsonFile

func ReadJsonFileOrPanic() error {
	JsonFile, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(JsonFile, &JsonConfigFile)
	if err != nil {
		panic(err)
	}
	return nil
}

var ConfFile ConfigYmlFile

func ReadYamlFileOrPanic() error {
	//fmt.Println("Load Config ")
	YamlFile, err := os.Open("config.yml")
	if err != nil {
		panic(err)
	}
	defer YamlFile.Close()
	decoder := yaml.NewDecoder(YamlFile)
	err = decoder.Decode(&ConfFile)
	if err != nil {
		panic(err)
	}

	return nil

}
