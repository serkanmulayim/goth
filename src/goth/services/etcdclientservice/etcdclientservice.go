package etcdclientservice

import (
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"time"

	"go.etcd.io/etcd/clientv3"
	"gopkg.in/yaml.v2"
)

const (
	defaultDialTimeout          = 1000
	EtcdClientAPIMiddleWareName = "etcdClient"
	EtcdClientClientTimeout     = 1000
)

type clientConfig struct {
	ClientDialTimeout int    `yaml:"client-dial-timeout,omitempty"`
	Endpoints         string `yaml:"listen-client-urls"`
}

//GetClient get etcdClient
func GetClient(filename string) (*clientv3.Client, error) {

	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Etcd ClientConfig file %v could not be found.", filename)
		return nil, err
	}
	var config clientConfig
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("ClientConfig could not be unmarshalled from file %v %v", filename, err)
		return nil, err
	}

	space := regexp.MustCompile(`\s+`)

	nospaceEndpoints := space.ReplaceAllString(config.Endpoints, "")
	endpoints := strings.Split(nospaceEndpoints, ",")

	var dialTimeout int
	if config.ClientDialTimeout != 0 {
		dialTimeout = config.ClientDialTimeout
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Duration(dialTimeout) * time.Millisecond,
	})

	if err != nil {
		log.Fatalf("Etcd Client could not be crated: %v", err)
		return nil, err
	}

	return cli, nil
}
