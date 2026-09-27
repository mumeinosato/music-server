package src

import (
	"crypto/sha256"
	"encoding/hex"
	"music-server/src/config"
	"music-server/src/db"
	"music-server/src/yt"
)

func Sync_List() error {
	list_id := config.Get().YouTubePlaylistID
	video_ids, err := yt.GetPlaylist(list_id)

	if err != nil {
		return err
	}

	hasher := sha256.New()

	for _, video_id := range video_ids {
		hasher.Write([]byte(video_id))
	}

	hash := hasher.Sum(nil)
	db.Update_List(video_ids, hex.EncodeToString(hash))
	return nil
}
