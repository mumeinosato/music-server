package main

import (
	"errors"

	"music-server/src/config"
	"music-server/src/db"
	"music-server/src"
	"music-server/src/yt"

	"github.com/gin-gonic/gin"
)

func main() {
	config.Get()

	db.Create_DB()

	r := gin.Default()

	r.GET("/sync", func(c *gin.Context) {
		changed, err := src.Sync_List()

		if errors.Is(err, src.ErrSyncInProgress) {
			c.JSON(409, gin.H{
				"message": err.Error(),
			})
			return
		}

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

		if !changed {
			c.JSON(200, gin.H{
				"message": "already up to date",
			})
			return
		}

		// ダウンロード・アップロードはバックグラウンドで続行中（結果はサーバーのログに出る）
		c.JSON(202, gin.H{
			"message": "sync started in background",
		})
	})

	r.GET("/is_syncing", func(c *gin.Context) {
		if src.IsSyncing() {
			c.JSON(200, gin.H{
				"syncing": true,
				"message": "sync is in progress",
			})
		} else {
			c.JSON(200, gin.H{
				"syncing": false,
				"message": "no sync in progress",
			})
		}
	})

	r.GET("/latest", func(c *gin.Context) {
		hash := c.Query("hash")
		l := src.Get_Latest(hash)

		c.JSON(200, gin.H{
			"hash":       l.Hash,
			"up_to_date": l.UpToDate,
			"add":        l.Add,
			"remove":     l.Remove,
			"all_id":     l.AllID,
			"all_name":   l.AllName,
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
