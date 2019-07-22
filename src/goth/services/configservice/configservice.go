package configservice

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

const (
	//ConfigContextName for the applicationMiddleWare
	ConfigContextName = "ApplicationConfig"
)

//AppConfig application config
type AppConfig struct {
	AppPort            int    `yaml:"application-port"`
	AdminPort          int    `yaml:"admin-port"`
	Fqdn               string `yaml:"application-fqdn"`
	AppTLSCert         string `yaml:"app-tls-cert"`
	AppTLSPrivateKey   string `yaml:"app-tls-privatekey"`
	AdminTLSCert       string `yaml:"admin-tls-cert"`
	AdminTLSPrivateKey string `yaml:"admin-tls-privatekey"`
}

//GetApplicationConfig get the pointer for application config
func GetApplicationConfig(filename string) (*AppConfig, error) {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Application configuration file %v could not be found.", filename)
		return nil, err // no need to return, exits 1
	}
	var config AppConfig
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Application config could not be unmarshalled, from file %v %v", filename, err)
		return nil, err
	}

	log.Printf("Application config: %+v\n", config)
	return &config, nil
}
