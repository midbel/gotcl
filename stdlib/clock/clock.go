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
	'a': "Mon",
	'A': "Monday",
	'b': "Jan",
	'B': "January",
	'd': "02", // day of month two digits
	'D': "01/02/2006",
	'e': "_2",
	'h': "Jan",
	'H': "15", // hour of the day 2 digits
	'I': "3",
	'j': "002", // day of year
	'k': "15",
	'K': "03",
	'm': "01", // months 2 digits
	'M': "04", // minutes 2 digits
	'N': "_1",
	'p': "PM",
	'R': "15:04",
	's': "",
	'S': "05",
	'T': "15:04:05",
	'y': "06",
	'Y': "2006", // year 4 digits
	'z': "-07:00",
	'Z': "Z",
	'+': "Mon Jan 02 15:04:05 Z 2006",
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
				if r == 's' {

				} else if r == 't' {
					part = "\t"
				} else if r == '%' {
					part = "%"
				} else {
					return "", fmt.Errorf("%c: unknown specifier", r)
				}
			}
			buf.WriteString(part)
		} else {
			buf.WriteRune(r)
		}
	}
	return buf.String(), nil
}
