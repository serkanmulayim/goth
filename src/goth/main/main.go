package main

import (
	"encoding/base64"
	"goth/cryptoutils"
	"goth/services/etcdclientservice"
	"goth/services/etcdservice"
	"goth/services/loginservice"
	"log"
	"net/http"

	"go.etcd.io/etcd/clientv3"

	"github.com/gin-gonic/gin"
)

const (
	etcdConfigFilePath = "./etcd.config.yaml"
)

func main() {

	startChan := make(chan bool)

	go etcdservice.Start(etcdConfigFilePath, startChan)
	<-startChan

	etcdClient, err := etcdclientservice.GetClient(etcdConfigFilePath)
	if err != nil {
		log.Fatalf("Etcd client could not start, exiting")
	}

	router := gin.Default()
	router.Use(etcdClientMiddleWare(etcdClient))

	api := router.Group("/api")
	api.Use(loginservice.AuthenticationFilterMiddleWare())
	public := router.Group("/api")

	api.GET("/token", func(c *gin.Context) {
		out, _ := cryptoutils.GenerateRandomBytes(32)
		o := base64.StdEncoding.EncodeToString(out)

		c.JSON(http.StatusOK, gin.H{
			"message": "token-" + o,
		})
	})
	api.GET("/auth", func(c *gin.Context) {
		out, _ := cryptoutils.GenerateRandomBytes(32)
		o := base64.StdEncoding.EncodeToString(out)
		c.JSON(http.StatusOK, gin.H{
			"message": "auth-" + o,
		})
	})
	api.GET("/checkauth", loginservice.CheckAuthRequestHandler)
	public.POST("/login", loginservice.LoginRequestHandler)
	router.Run(":3000")
}

func etcdClientMiddleWare(client *clientv3.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(etcdclientservice.EtcdClientAPIMiddleWareName, client)
		c.Next()
	}
}
