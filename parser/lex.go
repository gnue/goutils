package parser

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	itemError = -(1 + iota)
	itemEOF
	itemIdentifier
	itemText
	itemString
	itemComment
	itemKeyword
	itemExport
)

var tokenString = map[rune]string{
	itemError:      "ERROR",
	itemEOF:        "EOF",
	itemIdentifier: "identifier",
	itemText:       "text",
	itemString:     "string",
	itemComment:    "comment",
	itemExport:     "export",
}

var key = map[string]rune{
	"export": itemExport,
}

const eof = -1

type item struct {
	token rune
	val   string
}

func (i item) String() string {
	s := i.tokenString()

	switch i.token {
	case itemIdentifier, itemString, itemText, itemComment:
		return fmt.Sprintf("%s(%s)", s, i.value())
	}

	return s
}

func (i item) tokenString() string {
	if s, ok := tokenString[i.token]; ok {
		return s
	}
	return fmt.Sprintf("%q", string(i.token))
}

func (i item) value() string {
	switch i.token {
	case itemIdentifier, itemText:
		return i.val
	case itemString:
		s, _ := unquoteString(i.val)
		return s
	case itemComment:
		return strings.Trim(i.val, " \t")
	}

	return ""
}

type stateFn func(*lexer) stateFn

type lexer struct {
	input string
	start int
	pos   int
	width int
	pre   stateFn
	items chan item
}

func lex(input string) *lexer {
	l := &lexer{
		input: input,
		items: make(chan item),
	}
	go l.run()
	return l
}

func (l *lexer) run() {
	var curr stateFn

	for state := lexStart; state != nil; {
		l.pre = curr
		curr = state

		state = state(l)
	}
	close(l.items)
}

func (l *lexer) emit(t rune) {
	l.items <- item{t, l.text()}
	l.start = l.pos
}

func (l *lexer) text() string {
	return l.input[l.start:l.pos]
}

func (l *lexer) nextItem() item {
	item := <-l.items
	return item
}

func (l *lexer) next() (r rune) {
	if len(l.input) <= l.pos {
		l.width = 0
		return eof
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) accept(valid string) bool {
	if 0 <= strings.IndexRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(valid string) {
	for 0 <= strings.IndexRune(valid, l.next()) {
	}
	l.backup()
}

func (l *lexer) acceptN(valid string, n int) bool {
	pos := l.pos

	for 0 <= strings.IndexRune(valid, l.next()) {
		n -= 1
		if n == 0 {
			return true
		}
	}

	l.pos = pos
	return false
}

func (l *lexer) acceptRune(r rune) bool {
	if l.next() == r {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{
		itemError,
		fmt.Sprintf(format, args...),
	}
	return nil
}

// stateFn

func lexStart(l *lexer) stateFn {
	for {
		r := l.next()

		switch r {
		case eof:
			l.emit(itemEOF)
			return nil
		case '\r':
			l.acceptRune('\n')
			l.emit('\n')
		case '\n':
			l.emit('\n')
		case '#':
			l.ignore()
			return lexLineComment
		case '=', ':':
			l.emit(r)
			return lexRHS
		case '"', '\'':
			l.backup()
			return lexString
		default:
			if unicode.IsSpace(r) {
				l.ignore()
			} else {
				return lexIdentifier
			}
		}
	}
}

func lexRHS(l *lexer) stateFn {
	for {
		r := l.next()

		switch r {
		case eof, '\n', '\r':
			l.backup()
			return lexStart
		case '#':
			l.ignore()
			return lexLineComment
		case '"', '\'':
			l.backup()
			return lexString
		default:
			if unicode.IsSpace(r) {
				l.ignore()
			} else {
				return lexText
			}
		}
	}
}

func lexString(l *lexer) stateFn {
	qoute := l.next()

	for {
		switch l.next() {
		case eof:
			break
		case qoute:
			l.emit(itemString)
			return l.pre
		case '\\':
			_, _, err := espcaeRune(func() (rune, int) {
				r := l.next()
				return r, l.width
			})

			if err != nil {
				l.backup()
				return l.errorf("invalid escape")
			}
		}
	}

	return l.errorf("unclosed action")
}

func lexIdentifier(l *lexer) stateFn {
	defer func() {
		word := l.text()
		if token, ok := key[word]; ok {
			l.emit(token)
		} else {
			l.emit(itemIdentifier)
		}
	}()

	for {
		r := l.peek()

		switch r {
		case eof, '=', ':':
			return lexStart
		default:
			if unicode.IsSpace(r) {
				return lexStart
			}
		}

		l.next()
	}
}

func lexText(l *lexer) stateFn {
	defer l.emit(itemText)

	for {
		r := l.peek()

		switch {
		case r == eof, unicode.IsSpace(r):
			return lexStart
		}

		l.next()
	}
}

func lexLineComment(l *lexer) stateFn {
	defer l.emit(itemComment)

	for {
		switch l.peek() {
		case eof, '\n', '\r':
			return lexStart
		}

		l.next()
	}
}
