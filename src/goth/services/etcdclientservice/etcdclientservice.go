package etcdclientservice

import (
	"context"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"time"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/etcdserver/api/v3rpc/rpctypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/yaml.v2"
)

const (
	defaultDialTimeout = 1000
	//EtcdClientAPIMiddleWareName key to retrieve from gin context
	EtcdClientAPIMiddleWareName = "etcdClient"

	//EtcdClientClientTimeout in milliseconds
	EtcdClientClientTimeout  = 1000
	sessionLifeTimeSeconds   = (14 * 24 * 60 * 60)
	etcdClientPutErrorFormat = "Etcd Put error for key:\"%v\" error:%v"
)

type clientConfig struct {
	ClientDialTimeout int    `yaml:"client-dial-timeout,omitempty"`
	Endpoints         string `yaml:"listen-client-urls"`
}

//ETCDClient implements the interface to do all operations
type ETCDClient struct {
	Client *clientv3.Client
}

//ETCDIface mockable interface for accessing ETCD
type ETCDIface interface {
	Put(key string, value string) error
	Get(key string) (*clientv3.GetResponse, error)
	// GetAll(key string)
	Delete(key string) (*clientv3.DeleteResponse, error)
	// DeleteAll(key string)
}

//GetClient get etcdClient
func GetClient(filename string) (*ETCDClient, error) {

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

	out := ETCDClient{Client: cli}

	return &out, nil
}

//Get implementation for get one
func (cli *ETCDClient) Get(key string) (*clientv3.GetResponse, error) {

	ctx, cancel := context.WithTimeout(context.Background(), EtcdClientClientTimeout*time.Millisecond)
	gr, err := cli.Client.Get(ctx, key)
	cancel()
	return gr, err
}

//Delete implementation for delete 1
func (cli *ETCDClient) Delete(key string) (*clientv3.DeleteResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), EtcdClientClientTimeout*time.Millisecond)
	dr, err := cli.Client.Delete(ctx, key)
	cancel()
	return dr, err
}

// Put etcd put implementation
func (cli *ETCDClient) Put(key string, val string) error {
	ctx, cancel := context.WithTimeout(context.Background(), EtcdClientClientTimeout*time.Millisecond)

	lease, _ := cli.Client.Grant(ctx, sessionLifeTimeSeconds)

	_, err := cli.Client.Put(ctx, key, val, clientv3.WithLease(lease.ID))
	cancel()
	if err != nil {
		if err == context.Canceled {
			log.Printf(etcdClientPutErrorFormat, key, err)
		} else if err == context.DeadlineExceeded {
			log.Printf(etcdClientPutErrorFormat, key, err)
		} else if err == rpctypes.ErrEmptyKey {
			log.Println(etcdClientPutErrorFormat, key, err)
			// process (verr.Errors)
		} else if ev, ok := status.FromError(err); ok {
			code := ev.Code()
			if code == codes.DeadlineExceeded {
				// server-side context might have timed-out first (due to clock skew)
				// while original client-side context is not timed-out yet
				//will not happen since this is embedded etcd
				log.Println("server-side context might have timed-out first (due to clock skew)")

			}
		} else {
			//
			log.Println("bad cluster endpoints, which are not etcd servers:", err)
		}
		return err
	}
	return nil
}

//EtcdPut puts with etcd client
func EtcdPut(key string, val string, cli *clientv3.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), EtcdClientClientTimeout*time.Millisecond)

	lease, _ := cli.Grant(ctx, sessionLifeTimeSeconds)

	_, err := cli.Put(ctx, key, val, clientv3.WithLease(lease.ID))
	cancel()
	if err != nil {
		if err == context.Canceled {
			log.Printf(etcdClientPutErrorFormat, key, err)
		} else if err == context.DeadlineExceeded {
			log.Printf(etcdClientPutErrorFormat, key, err)
		} else if err == rpctypes.ErrEmptyKey {
			log.Println(etcdClientPutErrorFormat, key, err)
			// process (verr.Errors)
		} else if ev, ok := status.FromError(err); ok {
			code := ev.Code()
			if code == codes.DeadlineExceeded {
				// server-side context might have timed-out first (due to clock skew)
				// while original client-side context is not timed-out yet
				//will not happen since this is embedded etcd
				log.Println("server-side context might have timed-out first (due to clock skew)")

			}
		} else {
			//
			log.Println("bad cluster endpoints, which are not etcd servers:", err)
		}
		return err
	}
	return nil
}
