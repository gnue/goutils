package parser

import (
	"strconv"
	"unicode/utf8"
)

func unquoteString(s string) (string, error) {
	quote, w := utf8.DecodeRuneInString(s)
	if w == 0 {
		return "", errInvalidQuote
	}

	s = s[w:]
	buf := make([]rune, 0, len(s))

	for 0 < len(s) {
		r, w, err := unquoteRune(s, quote)
		if err == errEOS {
			break
		}

		if r != utf8.RuneError {
			buf = append(buf, r)
		}

		s = s[w:]
	}

	return string(buf), nil
}

func unquoteRune(s string, quote rune) (rune, int, error) {
	r, w := utf8.DecodeRuneInString(s)

	switch r {
	case quote:
		return r, w, errEOS
	case '\\':
		rest := s[w:]
		r, size, err := espcaeRune(func() (rune, int) {
			r, w := utf8.DecodeRuneInString(rest)
			rest = rest[w:]
			return r, w
		})
		w += size
		return r, w, err
	}

	return r, w, nil
}

type nextRuneFn func() (rune, int)

func espcaeRune(next nextRuneFn) (rune, int, error) {
	r, w := next()
	if w == 0 {
		return utf8.RuneError, w, errInvalidEscape
	}

	switch r {
	case 'a':
		r = '\a'
	case 'b':
		r = '\b'
	case 'f':
		r = '\f'
	case 'n':
		r = '\n'
	case 'r':
		r = '\r'
	case 't':
		r = '\t'
	case 'v':
		r = '\v'
	case 'x', 'u', 'U':
		n := 0
		switch r {
		case 'x':
			n = 2
		case 'u':
			n = 4
		case 'U':
			n = 8
		}

		hexdigit := make([]rune, n)

		for i := 0; i < n; i++ {
			c, size := next()
			if !isHexdigit(rune(c)) {
				return utf8.RuneError, w, errInvalidEscape
			}
			w += size
			hexdigit[i] = c
		}

		v, err := strconv.ParseUint(string(hexdigit), 16, 64)
		if err != nil {
			return utf8.RuneError, w, errInvalidEscape
		}

		r = rune(v)
	}

	return r, w, nil
}

func isHexdigit(r rune) bool {
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'a', 'b', 'c', 'd', 'e', 'f',
		'A', 'B', 'C', 'D', 'E', 'F':
		return true
	}

	return false
}
