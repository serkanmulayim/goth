package main

import (
	"goth/services/configservice"
	"goth/services/etcdclientservice"
	"goth/services/etcdservice"
	"goth/services/server"
	"log"
)

const (
	etcdConfigFilePath        = "./configs/etcd.config.yaml"
	applicationConfigFilePath = "./configs/app.yaml"
)

func main() {

	appConfig, _ := configservice.GetApplicationConfig(applicationConfigFilePath)

	startChan := make(chan bool)

	go etcdservice.Start(etcdConfigFilePath, startChan)
	<-startChan

	etcdClient, err := etcdclientservice.GetClient(etcdConfigFilePath)
	if err != nil {
		log.Fatalf("Etcd client could not start, exiting")
	}

	go server.StartClientApp(appConfig, etcdClient)
	server.StartAdminApp(appConfig, etcdClient)

}
