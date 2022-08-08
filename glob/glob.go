package glob

func Match(str, pattern string) bool {
	return match(str, pattern)
}

func Filter(list []string, pattern string) []string {
	var res []string
	for i := range list {
		if ok := match(list[i], pattern); !ok {
			continue
		}
		res = append(res, list[i])
	}
	return res
}

const (
	star     = '*'
	question = '?'
	lsquare  = '['
	rsquare  = ']'
)

func match(str, pattern string) bool {
	return false
}
