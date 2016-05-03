package parser

import ()

type callbackFn func(key, val string)

type parser struct {
	lexer    *lexer
	item     item
	err      error
	errs     []error
	callback callbackFn
}

func NewParser(input string) *parser {
	return &parser{lexer: lex(input)}
}

func (p *parser) next() {
	if p.item.token == itemEOF {
		return
	}

	for {
		p.item = p.lexer.nextItem()
		if p.item.token != itemComment {
			break
		}
	}
}

func (p *parser) accept(token rune) bool {
	if p.item.token == token {
		p.next()
		return true
	}

	return false
}

func (p *parser) expect(token rune) bool {
	if p.accept(token) {
		return true
	}

	p.err = &expectError{token}
	return false
}

func (p *parser) eof() bool {
	if p.item.token == itemEOF {
		return true
	}

	return false
}

func (p *parser) Parse(callback callbackFn) error {
	p.callback = callback

	p.next()
	for p.statement() || p.ignore() || p.recover() {
	}

	if p.errs != nil && 0 < len(p.errs) {
		return &multiError{p.errs}
	}

	return nil
}

func (p *parser) ignore() bool {
	ok := false

	for p.accept('\n') {
		ok = true
	}

	return ok
}

func (p *parser) error() bool {
	err := p.err

	if err == nil && !p.eof() {
		err = &parseError{p.item.val}
	}

	if err != nil {
		p.errs = append(p.errs, err)
		p.err = nil
	}

	return err != nil
}

func (p *parser) recover() bool {
	if p.error() {
		p.next()

		for !p.accept('\n') {
			p.next()
		}
	}

	return !p.eof()
}

func (p *parser) statement() bool {
	switch {
	case p.accept(itemExport):
		p.expression()
		return true
	case p.expression():
		return true
	}

	return false
}

func (p *parser) expression() bool {
	k := p.item
	if p.key() {
		if p.accept('=') || p.accept(':') {
			s := ""

			v := p.item
			if p.value() {
				s = v.value()
			}

			if p.callback != nil {
				p.callback(k.value(), s)
			}
		}

		if p.accept('\n') || p.accept(itemEOF) {
			return true
		}
	}

	return false
}

func (p *parser) key() bool {
	return p.accept(itemIdentifier) || p.accept(itemString)
}

func (p *parser) value() bool {
	return p.accept(itemText) || p.accept(itemString)
}
