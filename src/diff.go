package src

func Diff_List(new []string, old []string) (added []string, removed []string) {
	old_set := make(map[string]struct{}, len(old))
	for _, id := range old {
		old_set[id] = struct{}{}
	}

	new_set := make(map[string]struct{}, len(new))
	for _, id := range new {
		new_set[id] = struct{}{}
	}

	for _, id := range new {
		if _, exists := old_set[id]; !exists {
			added = append(added, id)
		}
	}

	for _, id := range old {
		if _, exists := new_set[id]; !exists {
			removed = append(removed, id)
		}
	}

	return added, removed
}