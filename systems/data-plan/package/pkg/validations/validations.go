package validationErrors

func IsEmpty(ss ...string) bool {
	for _, s := range ss {
		if s == "" {
			return true
		}
	}
	return false
}

func IsInvalidId(id uint64) bool {
	return id == 0
}
