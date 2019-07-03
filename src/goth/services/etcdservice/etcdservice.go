package etcdservice

import (
	"log"
	"time"

	"go.etcd.io/etcd/embed"
)

const (
	memberName  = "simple"
	clusterName = "simple-cluster"
	tempPrefix  = "simple-etcd-"

	// No peer URL exists but etcd doesn't allow the value to be empty.
	peerURL    = "http://localhost:0"
	clusterCfg = memberName + "=" + peerURL
)

// SimpleEtcd provides a single node etcd server.
// type SimpleEtcd struct {
// 	Port     int
// 	listener net.Listener
// 	server   *etcdserver.EtcdServer
// 	dataDir  string
// }

//Start start etcd service
func Start(configFilePath string, isStarted chan bool) {
	cfg, err := embed.ConfigFromFile(configFilePath)
	if err != nil {
		log.Fatal(err)
	}
	cfg.Dir = "default.etcd"
	e, err := embed.StartEtcd(cfg)

	if err != nil {
		log.Fatal(err)
	}
	defer e.Close()
	select {
	case <-e.Server.ReadyNotify():
		log.Println("Etcd Server is Ready")
		isStarted <- true

	case <-time.After(60 * time.Second):
		e.Server.Stop()
		log.Fatal("Server could not start due to timeout")
	}
	log.Fatal(<-e.Err())

}
