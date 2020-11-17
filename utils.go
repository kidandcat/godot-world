package main

func utilsBoolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func utilsStringToBool(b string) bool {
	if b == "true" {
		return true
	}
	return false
}

func utilsCheck(e error) {
	if e != nil {
		panic(e)
	}
}
