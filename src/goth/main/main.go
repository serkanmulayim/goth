package main

import (
	"encoding/base64"
	"goth/cryptoutils"
	"goth/services/configservice"
	"goth/services/etcdclientservice"
	"goth/services/etcdservice"
	"goth/services/loginservice"
	"log"
	"net/http"
	"strconv"

	"go.etcd.io/etcd/clientv3"

	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

const (
	etcdConfigFilePath        = "./etcd.config.yaml"
	applicationConfigFilePath = "./app.yaml"
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

	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("./clientbuild", true)))
	router.Use(cors.Default())
	router.Use(preflightHeadersMiddleWare())

	router.Use(applicationMiddleWare(appConfig))
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
	api.GET("/logout", loginservice.LogoutHandler)
	public.POST("/login", loginservice.LoginRequestHandler)
	router.RunTLS(":"+strconv.Itoa(appConfig.AppPort), appConfig.TLSCert, appConfig.TLSPrivateKey)
	//router.Run(":" + strconv.Itoa(appConfig.AppPort))

}

func applicationMiddleWare(appconfig *configservice.AppConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(configservice.ConfigContextName, appconfig)
		c.Next()
	}
}

//I think there is a bug in gin
func preflightHeadersMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "https://fe.localhost.goth.com:3000") //for testing
		c.Header("Access-Control-Allow-Headers", "GET, POST, PUT, DELETE, OPTIONS, HEAD")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0, post-check=0, pre-check=0")
		c.Header("Access-Control-Allow-Headers", "X-Requested-With, X-HTTP-Method-Override, Content-Type, Accept")
		c.Next()
	}
}

func etcdClientMiddleWare(client *clientv3.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(etcdclientservice.EtcdClientAPIMiddleWareName, client)
		c.Next()
	}
}
