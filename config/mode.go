package config

import "strings"

type YamlConfig struct {
	Bsc           NodeConfig
	PancakeRouter string
	FromAddress   string
	PrivateKey    string
	Mnemonic      string
	PriceToken    []string
	ShowSyrupPool []string
	BscToken      map[string]string
	SyrupPool     map[string]string
}

type NodeConfig struct {
	RpcUrl string
	WsUrl  string
}

var CF YamlConfig

func BscToken(symbol string) string {
	return CF.BscToken[strings.ToLower(symbol)]
}
