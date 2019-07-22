package loginservice

import (
	"context"
	"encoding/hex"
	"goth/cryptoutils"
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

	cookiePath        = "/"
	userSessionPrefix = "/usersession/"

	messageInternalServerError = "INTERNAL_SERVER_ERROR"
	messageUnauthorized        = "UNAUTHORIZED"
	messageOK                  = "OK"
	returnMessageFieldName     = "message"
)

//LoginRequestHandler Handle Login Request
func LoginRequestHandler(c *gin.Context) {
	sessionName := appSessionName
	isAdmin, _ := c.MustGet(IsAdminApiTypeContextName).(bool)
	if isAdmin {
		sessionName = adminSessionName
	}

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
		etcdSessionPath := etcdGetSessionPath(sessionID)
		err2 := etcdclientservice.EtcdPut(etcdSessionPath, username, etcdClient)

		if err2 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				returnMessageFieldName: messageInternalServerError,
			})
		} else {
			//userJSON, _ := json.Marshal(*u)
			c.SetCookie(sessionName, sessionID, expireSeconds, cookiePath, conf.Fqdn, true, true)
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

func etcdGetSessionPath(sessionID string) string {
	sessionKeySlice := []string{userSessionPrefix, sessionID}
	sessionPath := strings.Join(sessionKeySlice, "")
	return sessionPath
}

//CheckAuthRequestHandler checks auth
func CheckAuthRequestHandler(c *gin.Context) {
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

//LogoutHandler log
func LogoutHandler(c *gin.Context) {
	sessionName := appSessionName
	isAdmin, _ := c.MustGet(IsAdminApiTypeContextName).(bool)
	if isAdmin {
		sessionName = adminSessionName
	}
	etcdClient, ok := c.MustGet(etcdclientservice.EtcdClientAPIMiddleWareName).(*clientv3.Client)
	if !ok {
		log.Fatal("EtcdClient does not exist in LoginService")
	}

	conf, ok := c.MustGet(configservice.ConfigContextName).(*configservice.AppConfig)
	if !ok {
		log.Fatal("Application configuration could not be found")
	}

	sessionCookie, err := c.Request.Cookie(sessionName)
	if err != nil {
		c.JSON(http.StatusNoContent, gin.H{
			returnMessageFieldName: messageOK,
		})
		c.Abort()
		return
	}

	sessionID := sessionCookie.Value

	etcdRemoveSession(sessionID, etcdClient)
	c.SetCookie(sessionName, "x", -1, cookiePath, conf.Fqdn, true, true)
	c.JSON(http.StatusNoContent, gin.H{
		"message": "OK",
	})

}

//AuthenticationFilterMiddleWare authentication filter
func AuthenticationFilterMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionName := appSessionName
		isAdmin, _ := c.MustGet(IsAdminApiTypeContextName).(bool)
		if isAdmin {
			sessionName = adminSessionName
		}

		etcdClient, ok := c.MustGet(etcdclientservice.EtcdClientAPIMiddleWareName).(*clientv3.Client)
		if !ok {
			log.Fatal("EtcdClient does not exist in AuthenticationFilterMiddleWare")
		}

		sessionCookie, err := c.Request.Cookie(sessionName)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				returnMessageFieldName: messageUnauthorized,
			})
			c.Abort()
			return
		}

		sessionID := sessionCookie.Value
		username, err := etcdGetSession(sessionID, etcdClient)
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

func etcdRemoveSession(sessionID string, cli *clientv3.Client) {
	sessionPath := etcdGetSessionPath(sessionID)
	ctx, cancel := context.WithTimeout(context.Background(), etcdclientservice.EtcdClientClientTimeout*time.Millisecond)
	_, err := cli.Delete(ctx, sessionPath)
	cancel()
	if err != nil {
		log.Println("Error in removing sessionId for sessionID\". Deleting cookie anyways", sessionID, "\" :", err)
	}

}

func etcdGetSession(sessionID string, cli *clientv3.Client) (string, error) {
	sessionPath := etcdGetSessionPath(sessionID)
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
