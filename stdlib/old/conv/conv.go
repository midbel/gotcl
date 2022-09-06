package conv

const (
	zero = "0"
	one  = "1"
)

func Bool(ok bool) string {
	if ok {
		return True()
	}
	return False()
}

func False() string {
	return zero
}

func True() string {
	return one
}
