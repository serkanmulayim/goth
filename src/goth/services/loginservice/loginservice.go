package loginservice

import (
	"context"
	"encoding/hex"
	"goth/cryptoutils"
	"goth/services/etcdclientservice"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/etcdserver/api/v3rpc/rpctypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	expireSeconds          = (14 * 24 * 60 * 60)
	sessionLifeTimeSeconds = (14 * 24 * 60 * 60)
	sessionName            = "gothsess"
	cookieDomain           = ".goth.com"
	cookiePath             = "/"
	userSessionPrefix      = "/usersession/"

	messageInternalServerError = "INTERNAL_SERVER_ERROR"
	messageUnauthorized        = "UNAUTHORIZED"
	messageOK                  = "OK"
	returnMessageFieldName     = "message"
)

//LoginRequestHandler Handle Login Request
func LoginRequestHandler(c *gin.Context) {

	etcdClient, ok := c.MustGet(etcdclientservice.EtcdClientAPIMiddleWareName).(*clientv3.Client)
	if !ok {
		log.Fatal("EtcdClient does not exist in LoginService")
	}

	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")
	auth, err := authenticate(username, password)

	if auth && err == nil {

		sessionID := generateSessionID()
		etcdSessionPath := etcdGetSessionPath(sessionID)
		err2 := etcdPut(etcdSessionPath, username, etcdClient)

		if err2 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				returnMessageFieldName: messageInternalServerError,
			})
		} else {
			c.SetCookie(sessionName, sessionID, expireSeconds, cookiePath, cookieDomain, true, true)
			c.JSON(http.StatusOK, gin.H{
				returnMessageFieldName: messageOK,
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

func CheckAuthRequestHandler(c *gin.Context) {
	username, ok := c.MustGet("username").(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "UNAUTHORIZED",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message":  "OK",
			"username": username,
		})
	}
}

//AuthenticationFilterMiddleWare authentication filter
func AuthenticationFilterMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
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

func etcdPut(key string, val string, cli *clientv3.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), etcdclientservice.EtcdClientClientTimeout*time.Millisecond)

	lease, _ := cli.Grant(ctx, sessionLifeTimeSeconds)

	_, err := cli.Put(ctx, key, val, clientv3.WithLease(lease.ID))
	cancel()
	if err != nil {
		if err == context.Canceled {
			log.Println(err)
		} else if err == context.DeadlineExceeded {
			log.Println(err)
		} else if err == rpctypes.ErrEmptyKey {
			log.Println("Key is not provided, empty key")
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

func authenticate(u string, p string) (bool, error) {
	if u == "serkan" && p == "password" {
		return true, nil
	} else {
		return false, nil
	}
}

func generateSessionID() string {
	randomBytes, _ := cryptoutils.GenerateRandomBytes(64)
	return hex.EncodeToString(randomBytes)
}
