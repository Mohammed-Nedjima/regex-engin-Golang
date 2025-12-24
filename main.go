package main

type tokenType uint8

const (
	group tokenType = iota
	bracket
	or
	repeat
	literal
	groupUncaptured
)

type token struct {
	tokenType tokenType
	value     interface{}
}

type parseContext struct {
	pos    int
	tokens []token
}

func parse(regex string) *parseContext {
	ctx := &parseContext{
		pos:    0,
		tokens: []token{},
	}
	for ctx.pos < len(regex) {
		process(regex, ctx)
		ctx.pos++
	}
	return ctx
}

func process(regex string, ctx *parseContext) {
	ch := regex[ctx.pos]
	switch ch {
	case '(':
		groupCtx := &parseContext{
			pos:    0,
			tokens: []token{},
		}
		parseGroup(regex, groupCtx)
		ctx.tokens = append(ctx.tokens, token{tokenType: group, value: groupCtx.tokens})
	case '[':
		parseBracket(regex, ctx)
	case '|':
		parseOr(regex, ctx)
	case '*', '?', '+':
		parseRepeat(regex, ctx)
	case '{':
		parseRepeatSpecified(regex, ctx)
	default:
		t := token{
			tokenType: literal,
			value:     ch,
		}
		ctx.tokens = append(ctx.tokens, t)
	}
}

func parseGroup(regex string, groupCtx *parseContext) {
	groupCtx.pos++ // Jumping the opening bracket
	for regex[groupCtx.pos] != ')' {
		process(regex, groupCtx)
		groupCtx.pos++
	}
}

func parseBracket(regex string, ctx *parseContext) {
	ctx.pos++ // Jumping the first bracket
	ch := regex[ctx.pos]
	for ch != ']' {
		if ch == '-' {

		} else {

		}
		ctx.pos++
	}
}
