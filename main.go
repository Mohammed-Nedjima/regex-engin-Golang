package main

import "fmt"
import "strings"
import "strconv"

type tokenType uint8

const (
	group tokenType = iota
	bracket
	or
	repeat
	literal
	groupUncaptured
)
const (
	startOfText uint8 = iota
	endOfText
)
const repeatInfinity = -1

type token struct {
	tokenType tokenType
	value     interface{}
}

type parseContext struct {
	pos    int
	tokens []token
}
type repeatPayload struct {
	min   int
	max   int
	token token
}

type state struct {
	start       bool
	terminal    bool
	transitions map[uint8][]*state
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
	fmt.Println("parseGroup is working ...")
	groupCtx.pos++
	for regex[groupCtx.pos] != ')' {
		process(regex, groupCtx)
		groupCtx.pos++
	}
}

func parseBracket(regex string, ctx *parseContext) {
	fmt.Println("parseBracket is working ...")
	ctx.pos++
	var literals []string
	ch := regex[ctx.pos]
	for ch != ']' {
		if ch == '-' {
			next := regex[ctx.pos+1]
			prev := literals[len(literals)-1][0]
			literals[len(literals)-1] = fmt.Sprintf("%c%c", prev, next)
			ctx.pos++
		} else {
			literals = append(literals, fmt.Sprintf("%c", ch))
		}
		ctx.pos++
	}

	literalsSet := map[uint8]bool{}
	for _, l := range literals {
		for i := l[0]; i <= l[len(l)-1]; i++ {
			literalsSet[i] = true
		}
	}

	ctx.tokens = append(ctx.tokens, token{
		tokenType: bracket,
		value:     literalsSet,
	})
}

func parseRepeat(regex string, ctx *parseContext) {
	fmt.Println("parseRepeat is working ...")
	ch := regex[ctx.pos]
	var min, max int
	if ch == '*' {
		min = 0
		max = repeatInfinity
	} else if ch == '?' {
		min = 0
		max = 1
	} else {
		min = 1
		max = repeatInfinity
	}
	lastToken := ctx.tokens[len(ctx.tokens)-1]
	ctx.tokens[len(ctx.tokens)-1] = token{
		tokenType: repeat,
		value: repeatPayload{
			min:   min,
			max:   max,
			token: lastToken,
		},
	}
}

func parseRepeatSpecified(regex string, ctx *parseContext) {
	fmt.Println("parseRepeatSpecified is working ...")
	start := ctx.pos + 1
	for regex[ctx.pos] != '}' {
		ctx.pos++
	}

	boundariesStr := regex[start:ctx.pos]
	pieces := strings.Split(boundariesStr, ",")
	var min, max int
	if len(pieces) == 1 {
		if value, err := strconv.Atoi(pieces[0]); err != nil {
			panic(err.Error())
		} else {
			min = value
			max = value
		}
	} else if len(pieces) == 2 {
		if value, err := strconv.Atoi(pieces[0]); err != nil {
			panic(err.Error())
		} else {
			min = value
		}

		if pieces[1] == "" {
			max = repeatInfinity
		} else if value, err := strconv.Atoi(pieces[1]); err != nil {
			panic(err.Error())
		} else {
			max = value
		}
	} else {
		panic(fmt.Sprintf("There must be either 1 or 2 values specified for the quantifier: provided '%s'", boundariesStr))
	}

	lastToken := ctx.tokens[len(ctx.tokens)-1]
	ctx.tokens[len(ctx.tokens)-1] = token{
		tokenType: repeat,
		value: repeatPayload{
			min:   min,
			max:   max,
			token: lastToken,
		},
	}
}

func parseOr(regex string, ctx *parseContext) {
	fmt.Println("parseOr is working ...")
	rhsContext := &parseContext{
		pos:    ctx.pos,
		tokens: []token{},
	}
	rhsContext.pos += 1
	for rhsContext.pos < len(regex) && regex[rhsContext.pos] != ')' {
		process(regex, rhsContext)
		rhsContext.pos += 1
	}

	left := token{
		tokenType: groupUncaptured,
		value:     ctx.tokens,
	}

	right := token{
		tokenType: groupUncaptured,
		value:     rhsContext.tokens,
	}
	ctx.pos = rhsContext.pos

	ctx.tokens = []token{{
		tokenType: or,
		value:     []token{left, right},
	}}
}

/* building the NFA */

const epcilonChar uint8 = 0

func toNfa(ctx *parseContext) *state {
	startState, endState := tokenToNfa(&ctx.tokens[0])
	for i := 1; i < len(ctx.tokens); i++ {
		startNext, endNext := tokenToNfa(&ctx.tokens[i])
		endState.transitions[epcilonChar] = append(
			endState.transitions[epcilonChar],
			startNext,
		)
		endState = endNext
	}

	start := &state{
		start: true,
		transitions: map[uint8][]*state{
			epcilonChar: {startState},
		},
	}

	end := &state{
		terminal:    true,
		transitions: map[uint8][]*state{},
	}

	endState.transitions[epcilonChar] = append(
		endState.transitions[epcilonChar],
		end,
	)

	return start
}

func tokenToNfa(t *token) (*state, *state) {
	start := &state{
		start:       true,
		transitions: map[uint8][]*state{},
	}
	end := &state{
		terminal:    true,
		transitions: map[uint8][]*state{},
	}

	switch t.tokenType {
	case literal:
		ch := t.value.(uint8)
		start.transitions[ch] = []*state{end}
	case or:
		values := t.value.([]token)
		left := values[0]
		right := values[1]

		s1, e1 := tokenToNfa(&left)
		s2, e2 := tokenToNfa(&right)

		start.transitions = map[uint8][]*state{
			epcilonChar: []*state{s1, s2},
		}
		e1.transitions[epcilonChar] = []*state{end}
		e2.transitions[epcilonChar] = []*state{end}

	case bracket:
		literals := t.value.(map[uint8]bool)
		for l := range literals {
			start.transitions[l] = []*state{end}
		}

	case repeat:
		p := t.value.(repeatPayload)
		if p.min == 0 {
			start.transitions[epcilonChar] = []*state{end}
		}

		var copyCount int
		if p.max == repeatInfinity {
			if p.min == 0 {
				copyCount = 0
			} else {
				copyCount = p.min
			}
		} else {
			copyCount = p.max
		}

		from, to := tokenToNfa(&p.token)
		start.transitions[epcilonChar] = append(
			start.transitions[epcilonChar],
			from,
		)
		for i := 2; i <= copyCount; i++ {
			s, e := tokenToNfa(&p.token)
			to.transitions[epcilonChar] = append(
				to.transitions[epcilonChar],
				s,
			)

			from = s
			to = e

			if i > p.min {
				s.transitions[epcilonChar] = append(
					s.transitions[epcilonChar],
					end,
				)
			}
		}
		to.transitions[epcilonChar] = append(
			to.transitions[epcilonChar],
			end,
		)

		if p.max == repeatInfinity {
			end.transitions[epcilonChar] = append(
				end.transitions[epcilonChar],
				from,
			)
		}

	case group, groupUncaptured:
		tokens := t.value.([]token)
		start, end = tokenToNfa(&tokens[0])
		for i := 1; i < len(tokens); i++ {
			ts, te := tokenToNfa(&tokens[i])
			end.transitions[epcilonChar] = append(
				end.transitions[epcilonChar],
				ts,
			)
			end = te
		}
	default:
		panic("unknown type of token")
	}

	return start, end
}

// the matching logic

func getChar(str string, pos int) uint8 {
	if pos >= len(str) {
		return endOfText
	}
	if pos < 0 {
		return startOfText
	}
	return str[pos]
}

func (s *state) check(str string, pos int) bool {
	ch := getChar(str, pos)

	if ch == endOfText && s.terminal {
		return true
	}
	if states := s.transitions[ch]; len(states) > 0 {
		nextState := states[0]
		if nextState.check(str, pos+1) {
			return true
		}

		for _, state := range s.transitions[epcilonChar] {
			if state.check(str, pos) {
				return true
			}
		}

		if ch == startOfText && s.check(str, pos) {
			return true
		}
	}
	return false
}

/*
func main() {
	fmt.Println("hello world")
	ctx := parse("[a-zA-Z][a-zA-Z0-9_.]+@[a-zA-Z0-9]+.[a-zA-Z]{2,}")
	fmt.Println("finished parsing")
	nfa := toNfa(ctx)
	fmt.Println("converted to nfa successfully")
	if nfa.check("ayoub@gmail.com", -1) {
		fmt.Println("email is valid")
	} else {
		fmt.Println("email not valid")
	}
}
*/
