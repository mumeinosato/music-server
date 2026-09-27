package main

import (
	"errors"

	"music-server/config"
	"music-server/db"
	"music-server/src"
	"music-server/yt"

	"github.com/gin-gonic/gin"
)

func main() {
	config.Get()

	db.Create_DB()

	r := gin.Default()

	r.GET("/sync", func(c *gin.Context) {
		err := src.Sync_List()

		var auth_err *yt.AuthRequiredError
		if errors.As(err, &auth_err) {
			// c.JSON だと URL 内の & が & にエスケープされるため PureJSON を使う
			c.PureJSON(401, gin.H{
				"message":  "auth required",
				"auth_url": auth_err.URL,
			})
			return
		}
		if err != nil {
			c.JSON(500, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "sync completed",
		})
	})

	// Google 認証後のリダイレクト先（SERVER_URL + "/callback"）
	r.GET("/callback", func(c *gin.Context) {
		if err := yt.HandleCallback(c.Query("code"), c.Query("state")); err != nil {
			c.JSON(400, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "auth completed",
		})
	})

	r.Run(":" + config.Get().Port)
}
