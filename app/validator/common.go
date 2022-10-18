package validator

func IsAlphanumericMax30(val string) bool {
	if len(val) == 0 || len(val) > 30 {
		return false
	}

	for _, r := range val {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') {
			return false
		}
	}

	return true
}
