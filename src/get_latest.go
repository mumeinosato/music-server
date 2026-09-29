package src

import "music-server/src/db"

// Latest は /latest のレスポンス内容
type Latest struct {
	Hash string
	// UpToDate が true のとき、クライアントは最新。Add/Remove/AllID/AllName は空
	UpToDate bool
	Add      []string
	Remove   []string
	AllID    []string
	AllName  []string
}

// Get_Latest はクライアントが持つ hash を最新の hash と比べ、一致すれば何も走査せず返す。
// 一致しないときだけ差分と全ID・曲名を作る
func Get_Latest(hash string) Latest {
	all_id, latest_hash := db.Get_List(db.Mode.Latest)

	if hash != "" && hash == latest_hash {
		return Latest{Hash: latest_hash, UpToDate: true, Add: []string{}, Remove: []string{}, AllID: []string{}, AllName: []string{}}
	}

	if all_id == nil {
		all_id = []string{}
	}
	all_name := db.Get_Music_Name(all_id)

	if hash == "" {
		return Latest{Hash: latest_hash, Add: all_id, Remove: []string{}, AllID: all_id, AllName: all_name}
	}

	// 知らない hash なら old_ids は nil になり、全件が add になる
	old_ids, _ := db.Get_List(db.Mode.From_Hash, hash)
	add, remove := Diff_List(all_id, old_ids)
	if add == nil {
		add = []string{}
	}
	if remove == nil {
		remove = []string{}
	}
	return Latest{Hash: latest_hash, Add: add, Remove: remove, AllID: all_id, AllName: all_name}
}
