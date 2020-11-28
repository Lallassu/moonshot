package main

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	KeyLoadFuel   string            `json:"keyLoadFuel"`
	KeyDumpFuel   string            `json:"keyDumpFuel"`
	KeyRelease    string            `json:"keyRelease"`
	KeyRespawn    string            `json:"keyRespawn"`
	KeyDebugInfo  string            `json:"keyDebugInfo"`
	KeyMenu       string            `json:"keyMenu"`
	KeyWireframe  string            `json:"keyWirefram"`
	KeyMenuUp     string            `json:"keyMenuUp"`
	KeyMenuDown   string            `json:"keyMenuDown"`
	KeyMenuSelect string            `json:"keyMenuSelect"`
	Assets        map[string]string `json:"assets"`
	Sounds        map[string]string `json:"sounds"`
	Shaders       map[string]string `json:"shaders"`
	Colors        map[string]Color  `json:"Colors"`
}

type Color struct {
	R float32 `json:"r"`
	G float32 `json:"g"`
	B float32 `json:"b"`
	A float32 `json:"a"`
}

func LoadConfiguration(file string) Config {
	var conf Config

	cFile, err := os.Open(file)
	if err != nil {
		log.Fatal("Failed to open configuration file", file, err)
		panic(err)
	}

	jsonParser := json.NewDecoder(cFile)
	err = jsonParser.Decode(&conf)
	if err != nil {
		log.Fatal("Failed to parse config file:", err)
	}
	return conf
}
