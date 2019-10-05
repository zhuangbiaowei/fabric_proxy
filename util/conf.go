package util

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var GlobalConfig = new(Conf)

type Conf struct {
	Endpoint         string   `yaml:"endpoint"`
	Path             string   `yaml:"path"`
	CryptoMethod     string   `yaml:"cryptoMethod"`
	SignCert         string   `yaml:"signCert"`
	PrvKey           string   `yaml:"prvKey"`
	OrgId            string   `yaml:"orgId"`
	ChannelId        string   `yaml:"channelId"`
	ChaincodeId      string   `yaml:"chaincodeId"`
	ChaincodeVersion string   `yaml:"chaincodeVersion"`
	UserId           string   `yaml:"userId"`
	Orderers         []string `yaml:"orderers"`
	Peers            []string `yaml:"peers"`
}

func InitConfig(path string) {
	data, _ := ioutil.ReadFile(path)
	if err := yaml.Unmarshal(data, GlobalConfig); err != nil {
		fmt.Println("Pares config file err", err.Error())
	}
}