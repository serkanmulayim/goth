package server

import (
	"encoding/base64"
	"goth/cryptoutils"
	"goth/objects/admin"
	"goth/services/configservice"
	"goth/services/etcdclientservice"
	"goth/services/loginservice"
	"net/http"
	"strconv"

	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"go.etcd.io/etcd/clientv3"
)

// StartAdminApp starts the admin servers with the apis
func StartAdminApp(appConfig *configservice.AppConfig, etcdClient *clientv3.Client) {

	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("./adminbuild", true)))
	router.Use(cors.Default())
	router.Use(preflightHeadersMiddleWare())

	router.Use(applicationMiddleWare(appConfig, true))
	router.Use(etcdClientMiddleWare(etcdClient))

	public := router.Group("/api")
	api := router.Group("/api")
	api.Use(loginservice.AdminAuthenticationFilterMiddleWare())

	api.GET("/auth", func(c *gin.Context) {
		out, _ := cryptoutils.GenerateRandomBytes(32)
		o := base64.StdEncoding.EncodeToString(out)
		c.JSON(http.StatusOK, gin.H{
			"message": "auth-" + o,
		})
	})

	admins := [3]admin.Object{
		admin.Object{FirstName: "Serkan", LastName: "Mulayim", Email: "ser$%kan@gmail.com", Phone: "5555555", Address: "400 3rd street", UserID: 1},
		admin.Object{FirstName: "Gozde", LastName: "Nulayim", Email: "gozde@gmail.com", Phone: "33333333", Address: "500 4th street", UserID: 2},
		admin.Object{FirstName: "Ali", LastName: "Mulayim", Email: "ali@gmail.com", Phone: "33333333", Address: "500 4th street", UserID: 3},
	}

	api.GET("/admins", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
			"admins":  admins,
		})
	})

	api.GET("/checkauth", loginservice.AdminCheckAuthRequestHandler)
	public.GET("/logout", loginservice.AdminLogoutHandler)
	public.POST("/login", loginservice.AdminLoginRequestHandler)
	router.RunTLS(":"+strconv.Itoa(appConfig.AdminPort), appConfig.AdminTLSCert, appConfig.AdminTLSPrivateKey)

}

//StartClientApp starts the user application
func StartClientApp(appConfig *configservice.AppConfig, etcdClient *clientv3.Client) {

	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("./clientbuild", true)))
	router.Use(cors.Default())
	router.Use(preflightHeadersMiddleWare())

	router.Use(applicationMiddleWare(appConfig, false))
	router.Use(etcdClientMiddleWare(etcdClient))

	public := router.Group("/api")
	api := router.Group("/api")
	api.Use(loginservice.AppAuthenticationFilterMiddleWare())

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
	api.GET("/checkauth", loginservice.AppCheckAuthRequestHandler)
	public.GET("/logout", loginservice.AppLogoutHandler)
	public.POST("/login", loginservice.AppLoginRequestHandler)
	router.RunTLS(":"+strconv.Itoa(appConfig.AppPort), appConfig.AppTLSCert, appConfig.AppTLSPrivateKey)

}

func applicationMiddleWare(appconfig *configservice.AppConfig, isAdminAPI bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(configservice.ConfigContextName, appconfig)
		c.Set(loginservice.IsAdminApiTypeContextName, isAdminAPI)
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
