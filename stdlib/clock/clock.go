package clock

import (
	"fmt"
	"strings"
	"time"
)

func Format(unix int64, pat string) (string, error) {
	when := time.Unix(unix, 0)
	pat, err := parseFormat(pat)
	if err != nil {
		return "", err
	}
	return when.Format(pat), nil
}

func Scan(str, pat string) (int64, error) {
	pat, err := parseFormat(pat)
	if err != nil {
		return 0, err
	}
	when, err := time.Parse(pat, str)
	if err != nil {
		return 0, err
	}
	return when.Unix(), nil
}

var rules = map[rune]string{
	'd': "02",   // day of month two digits
	'j': "002",  // day of year
	'm': "01",   // months 2 digits
	'M': "04",   // minutes 2 digits
	'Y': "2006", // year 4 digits
	'H': "15",   // hour of the day 2 digits
}

func parseFormat(str string) (string, error) {
	var (
		rs  = strings.NewReader(str)
		buf strings.Builder
	)
	for rs.Len() > 0 {
		r, _, _ := rs.ReadRune()
		if r == '%' {
			r, _, _ = rs.ReadRune()
			part, ok := rules[r]
			if !ok {
				return "", fmt.Errorf("%c: unknown specifier", r)
			}
			buf.WriteString(part)
		} else {
			buf.WriteRune(r)
		}
	}
	return buf.String(), nil
}
