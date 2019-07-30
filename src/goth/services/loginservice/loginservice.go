package loginservice

import (
	"context"
	"encoding/hex"
	"goth/cryptoutils"
	"goth/objects/admin"
	"goth/objects/gothuser"
	"goth/services/configservice"
	"goth/services/etcdclientservice"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.etcd.io/etcd/clientv3"
)

const (
	expireSeconds = (14 * 24 * 60 * 60)

	IsAdminApiTypeContextName = "apiType"

	appSessionName   = "gothsess"
	adminSessionName = "agothsess"

	cookiePath         = "/"
	userSessionPrefix  = "/usersession/"
	adminSessionPrefix = "/adminsession/"

	messageInternalServerError = "INTERNAL_SERVER_ERROR"
	messageUnauthorized        = "UNAUTHORIZED"
	messageOK                  = "OK"
	returnMessageFieldName     = "message"
)

//AppLoginRequestHandler Handle Login Request
func AppLoginRequestHandler(c *gin.Context) {

	etcdClient, ok := c.MustGet(etcdclientservice.EtcdClientAPIMiddleWareName).(*clientv3.Client)
	if !ok {
		log.Fatal("EtcdClient does not exist in LoginService")
	}

	conf, ok := c.MustGet(configservice.ConfigContextName).(*configservice.AppConfig)
	if !ok {
		log.Fatal("Application configuration could not be found")
	}

	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")
	u, err := authenticate(username, password)

	if u != nil && err == nil {

		sessionID := generateSessionID()
		etcdSessionPath := etcdGetSessionPath(sessionID, false)
		err2 := etcdclientservice.EtcdPut(etcdSessionPath, username, etcdClient)

		if err2 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				returnMessageFieldName: messageInternalServerError,
			})
		} else {
			//userJSON, _ := json.Marshal(*u)
			c.SetCookie(appSessionName, sessionID, expireSeconds, cookiePath, conf.Fqdn, true, true)
			c.JSON(http.StatusOK, gin.H{
				returnMessageFieldName: messageOK,
				"user":                 *u, //string(userJSON),
			})
		}

	} else if err == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			returnMessageFieldName: messageUnauthorized,
		})

	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			returnMessageFieldName: messageInternalServerError,
		})
	}

}

//AdminLoginRequestHandler Handle Login Request
func AdminLoginRequestHandler(c *gin.Context) {

	etcdClient, ok := c.MustGet(etcdclientservice.EtcdClientAPIMiddleWareName).(*clientv3.Client)
	if !ok {
		log.Fatal("EtcdClient does not exist in LoginService")
	}

	conf, ok := c.MustGet(configservice.ConfigContextName).(*configservice.AppConfig)
	if !ok {
		log.Fatal("Application configuration could not be found")
	}

	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")
	u, err := authenticateAdmin(username, password)

	if u != nil && err == nil {

		sessionID := generateSessionID()
		etcdSessionPath := etcdGetSessionPath(sessionID, true)
		err2 := etcdclientservice.EtcdPut(etcdSessionPath, username, etcdClient)

		if err2 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				returnMessageFieldName: messageInternalServerError,
			})
		} else {
			//userJSON, _ := json.Marshal(*u)
			c.SetCookie(adminSessionName, sessionID, expireSeconds, cookiePath, conf.Fqdn, true, true)
			c.JSON(http.StatusOK, gin.H{
				returnMessageFieldName: messageOK,
				"user":                 *u, //string(userJSON),
			})
		}

	} else if err == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			returnMessageFieldName: messageUnauthorized,
		})

	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			returnMessageFieldName: messageInternalServerError,
		})
	}

}

func etcdGetSessionPath(sessionID string, isAdmin bool) string {
	var sessionKeySlice []string
	if isAdmin {
		sessionKeySlice = []string{userSessionPrefix, sessionID}
	} else {
		sessionKeySlice = []string{adminSessionPrefix, sessionID}
	}

	sessionPath := strings.Join(sessionKeySlice, "")
	return sessionPath
}

//AppCheckAuthRequestHandler checks auth
func AppCheckAuthRequestHandler(c *gin.Context) {
	username, ok := c.MustGet("username").(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "UNAUTHORIZED",
		})
	} else {

		u := gothuser.Object{Username: username}
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
			"user":    u,
		})
	}
}

//AdminCheckAuthRequestHandler checks auth
func AdminCheckAuthRequestHandler(c *gin.Context) {
	_, ok := c.MustGet("admin").(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "UNAUTHORIZED",
		})
	} else {

		admin := admin.Object{FirstName: "Serkan", LastName: "Mulayim", Email: "ser$%kan@gmail.com", Phone: "5555555", Address: "400 3rd street", UserID: 1}
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
			"admin":   admin,
		})
	}
}

//AppLogoutHandler log
func AppLogoutHandler(c *gin.Context) {

	etcdClient, ok := c.MustGet(etcdclientservice.EtcdClientAPIMiddleWareName).(*clientv3.Client)
	if !ok {
		log.Fatal("EtcdClient does not exist in LoginService")
	}

	conf, ok := c.MustGet(configservice.ConfigContextName).(*configservice.AppConfig)
	if !ok {
		log.Fatal("Application configuration could not be found")
	}

	sessionCookie, err := c.Request.Cookie(appSessionName)
	if err != nil {
		c.JSON(http.StatusNoContent, gin.H{
			returnMessageFieldName: messageOK,
		})
		c.Abort()
		return
	}

	sessionID := sessionCookie.Value

	etcdRemoveSession(sessionID, false, etcdClient)
	c.SetCookie(appSessionName, "x", -1, cookiePath, conf.Fqdn, true, true)
	c.JSON(http.StatusNoContent, gin.H{
		"message": "OK",
	})

}

//AdminLogoutHandler log
func AdminLogoutHandler(c *gin.Context) {

	etcdClient, ok := c.MustGet(etcdclientservice.EtcdClientAPIMiddleWareName).(*clientv3.Client)
	if !ok {
		log.Fatal("EtcdClient does not exist in LoginService")
	}

	conf, ok := c.MustGet(configservice.ConfigContextName).(*configservice.AppConfig)
	if !ok {
		log.Fatal("Application configuration could not be found")
	}

	sessionCookie, err := c.Request.Cookie(adminSessionName)
	if err != nil {
		c.JSON(http.StatusNoContent, gin.H{
			returnMessageFieldName: messageOK,
		})
		c.Abort()
		return
	}

	sessionID := sessionCookie.Value

	etcdRemoveSession(sessionID, true, etcdClient)
	c.SetCookie(adminSessionName, "x", -1, cookiePath, conf.Fqdn, true, true)
	c.JSON(http.StatusNoContent, gin.H{
		"message": "OK",
	})

}

//AppAuthenticationFilterMiddleWare authentication filter
func AppAuthenticationFilterMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {

		etcdClient, ok := c.MustGet(etcdclientservice.EtcdClientAPIMiddleWareName).(*clientv3.Client)
		if !ok {
			log.Fatal("EtcdClient does not exist in AuthenticationFilterMiddleWare")
		}

		sessionCookie, err := c.Request.Cookie(appSessionName)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				returnMessageFieldName: messageUnauthorized,
			})
			c.Abort()
			return
		}

		sessionID := sessionCookie.Value
		username, err := etcdGetSession(sessionID, false, etcdClient)
		if err != nil || username == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				returnMessageFieldName: messageUnauthorized,
			})
			c.Abort()
			return
		}
		c.Set("username", username)
		c.Next()
	}
}

//AdminAuthenticationFilterMiddleWare authentication filter
func AdminAuthenticationFilterMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {

		etcdClient, ok := c.MustGet(etcdclientservice.EtcdClientAPIMiddleWareName).(*clientv3.Client)
		if !ok {
			log.Fatal("EtcdClient does not exist in AuthenticationFilterMiddleWare")
		}

		sessionCookie, err := c.Request.Cookie(adminSessionName)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				returnMessageFieldName: messageUnauthorized,
			})
			c.Abort()
			return
		}

		sessionID := sessionCookie.Value
		adminname, err := etcdGetSession(sessionID, true, etcdClient)
		if err != nil || adminname == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				returnMessageFieldName: messageUnauthorized,
			})
			c.Abort()
			return
		}
		c.Set("admin", adminname)
		c.Next()
	}
}

func etcdRemoveSession(sessionID string, isAdmin bool, cli *clientv3.Client) {
	sessionPath := etcdGetSessionPath(sessionID, isAdmin)
	ctx, cancel := context.WithTimeout(context.Background(), etcdclientservice.EtcdClientClientTimeout*time.Millisecond)
	_, err := cli.Delete(ctx, sessionPath)
	cancel()
	if err != nil {
		log.Println("Error in removing sessionId for sessionID\". Deleting cookie anyways", sessionID, "\" :", err)
	}

}

func etcdGetSession(sessionID string, isAdmin bool, cli *clientv3.Client) (string, error) {
	sessionPath := etcdGetSessionPath(sessionID, isAdmin)
	ctx, cancel := context.WithTimeout(context.Background(), etcdclientservice.EtcdClientClientTimeout*time.Millisecond)
	gr, err := cli.Get(ctx, sessionPath)
	cancel()

	if err != nil {
		log.Println("Error in getting sessionId for sessionID\"", sessionID, "\" :", err)
		return "", err
	}

	if len(gr.Kvs) == 0 {
		return "", nil
	}

	return string(gr.Kvs[0].Value), nil
}

func authenticateAdmin(u string, p string) (*admin.Object, error) {
	if u == "serkan" && p == "password" {
		admin := admin.Object{FirstName: "Serkan", LastName: "Mulayim", Email: "ser$%kan@gmail.com", Phone: "5555555", Address: "400 3rd street", UserID: 1}
		return &admin, nil
	}
	return nil, nil
}

func authenticate(u string, p string) (*gothuser.Object, error) {
	if (u == "serkan" || u == "serkan2") && p == "password" {
		obj := gothuser.Object{Username: u}
		return &obj, nil
	}
	return nil, nil
}

func generateSessionID() string {
	randomBytes, _ := cryptoutils.GenerateRandomBytes(64)
	return hex.EncodeToString(randomBytes)
}
