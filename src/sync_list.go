package src

import (
	"errors"
	"log"
	"music-server/src/config"
	"music-server/src/db"
	"music-server/src/r2"
	"music-server/src/yt"
	"slices"
	"sync"
	"sync/atomic"
)

var ErrSyncInProgress = errors.New("sync is already in progress")

var sync_mu sync.Mutex
var syncing atomic.Bool

func IsSyncing() bool {
	return syncing.Load()
}

func unlock_sync() {
	syncing.Store(false)
	sync_mu.Unlock()
}

func Sync_List() (changed bool, err error) {
	if !sync_mu.TryLock() {
		return false, ErrSyncInProgress
	}
	syncing.Store(true)

	handed_off := false
	defer func() {
		if !handed_off {
			unlock_sync()
		}
	}()

	list_id := config.Get().YouTubePlaylistID
	video_ids, video_names, err := yt.GetPlaylist(list_id)

	db.Add_Music_Name(video_ids, video_names)

	if err != nil {
		return false, err
	}

	// DB には前回までに R2 に実際にあるIDが保存されている
	old, _ := db.Get_List(db.Mode.Latest)

	if slices.Equal(video_ids, old) {
		return false, nil
	}

	added, removed := Diff_List(video_ids, old)

	handed_off = true
	go func() {
		defer unlock_sync()
		finish_sync(video_ids, old, added, removed)
	}()

	return true, nil
}


func finish_sync(video_ids []string, old []string, added []string, removed []string) {
	defer yt.Clean_Temp()

	var errs []error

	failed_add := make(map[string]struct{})

	if len(added) > 0 {
		files, err := yt.Download(added)
		if err != nil {
			errs = append(errs, err)
		}
		for _, id := range added {
			if _, ok := files[id]; !ok {
				failed_add[id] = struct{}{}
			}
		}

		failed, err := r2.Upload(files)
		if err != nil {
			errs = append(errs, err)
		}
		for _, id := range failed {
			failed_add[id] = struct{}{}
		}
	}

	var failed_remove []string
	if len(removed) > 0 {
		var err error
		failed_remove, err = r2.Delete(removed)
		if err != nil {
			errs = append(errs, err)
		}
	}

	saved := make([]string, 0, len(video_ids)+len(failed_remove))
	for _, id := range video_ids {
		if _, ok := failed_add[id]; !ok {
			saved = append(saved, id)
		}
	}
	// 削除に失敗したものは R2 に残っているので、リストに残して次回また削除させる
	saved = append(saved, failed_remove...)

	if err := errors.Join(errs...); err != nil {
		log.Printf("Sync finished with errors: %v", err)
	}

	if !slices.Equal(saved, old) {
		db.Update_List(saved)
	}
	log.Printf("Sync completed (added %d items / removed %d items / failed %d items)",
		len(added)-len(failed_add), len(removed)-len(failed_remove), len(failed_add)+len(failed_remove))
}
