package svg

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

const endSymbol rune = 1114112

/* The rule types inferred from the grammar are below. */
type pegRule uint8

const (
	ruleUnknown pegRule = iota
	rulepath
	rulemovetoDrawtos
	rulemovetoDrawto
	ruledrawtos
	ruledrawto
	rulemoveto
	rulemovetoArgs
	rulelineto
	ruleclosepath
	rulelinetoArgs
	rulecubicto
	rulecubictoArgs
	rulecubictoArg
	rulequadto
	rulequadtoArgs
	rulequadtoArg
	ruledigits
	rulepair
	rulecoord
	rulenumber
	rulenonneg
	rulecomma
	ruleinteger
	rulefloat
	rulefract
	ruleexponent
	rulewsp
	rulesign
	ruleAction0
	ruleAction1
	ruleAction2
	ruleAction3
	ruleAction4
	ruleAction5
	ruleAction6
	ruleAction7
	ruleAction8
	rulePegText
	ruleAction9

	rulePre
	ruleIn
	ruleSuf
)

var rul3s = [...]string{
	"Unknown",
	"path",
	"movetoDrawtos",
	"movetoDrawto",
	"drawtos",
	"drawto",
	"moveto",
	"movetoArgs",
	"lineto",
	"closepath",
	"linetoArgs",
	"cubicto",
	"cubictoArgs",
	"cubictoArg",
	"quadto",
	"quadtoArgs",
	"quadtoArg",
	"digits",
	"pair",
	"coord",
	"number",
	"nonneg",
	"comma",
	"integer",
	"float",
	"fract",
	"exponent",
	"wsp",
	"sign",
	"Action0",
	"Action1",
	"Action2",
	"Action3",
	"Action4",
	"Action5",
	"Action6",
	"Action7",
	"Action8",
	"PegText",
	"Action9",

	"Pre_",
	"_In_",
	"_Suf",
}

type node32 struct {
	token32
	up, next *node32
}

func (node *node32) print(depth int, buffer string) {
	for node != nil {
		for c := 0; c < depth; c++ {
			fmt.Printf(" ")
		}
		fmt.Printf("\x1B[34m%v\x1B[m %v\n", rul3s[node.pegRule], strconv.Quote(string(([]rune(buffer)[node.begin:node.end]))))
		if node.up != nil {
			node.up.print(depth+1, buffer)
		}
		node = node.next
	}
}

func (node *node32) Print(buffer string) {
	node.print(0, buffer)
}

type element struct {
	node *node32
	down *element
}

/* ${@} bit structure for abstract syntax tree */
type token32 struct {
	pegRule
	begin, end, next uint32
}

func (t *token32) isZero() bool {
	return t.pegRule == ruleUnknown && t.begin == 0 && t.end == 0 && t.next == 0
}

func (t *token32) isParentOf(u token32) bool {
	return t.begin <= u.begin && t.end >= u.end && t.next > u.next
}

func (t *token32) getToken32() token32 {
	return token32{pegRule: t.pegRule, begin: uint32(t.begin), end: uint32(t.end), next: uint32(t.next)}
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v %v", rul3s[t.pegRule], t.begin, t.end, t.next)
}

type tokens32 struct {
	tree    []token32
	ordered [][]token32
}

func (t *tokens32) trim(length int) {
	t.tree = t.tree[0:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) Order() [][]token32 {
	if t.ordered != nil {
		return t.ordered
	}

	depths := make([]int32, 1, math.MaxInt16)
	for i, token := range t.tree {
		if token.pegRule == ruleUnknown {
			t.tree = t.tree[:i]
			break
		}
		depth := int(token.next)
		if length := len(depths); depth >= length {
			depths = depths[:depth+1]
		}
		depths[depth]++
	}
	depths = append(depths, 0)

	ordered, pool := make([][]token32, len(depths)), make([]token32, len(t.tree)+len(depths))
	for i, depth := range depths {
		depth++
		ordered[i], pool, depths[i] = pool[:depth], pool[depth:], 0
	}

	for i, token := range t.tree {
		depth := token.next
		token.next = uint32(i)
		ordered[depth][depths[depth]] = token
		depths[depth]++
	}
	t.ordered = ordered
	return ordered
}

type state32 struct {
	token32
	depths []int32
	leaf   bool
}

func (t *tokens32) AST() *node32 {
	tokens := t.Tokens()
	stack := &element{node: &node32{token32: <-tokens}}
	for token := range tokens {
		if token.begin == token.end {
			continue
		}
		node := &node32{token32: token}
		for stack != nil && stack.node.begin >= token.begin && stack.node.end <= token.end {
			stack.node.next = node.up
			node.up = stack.node
			stack = stack.down
		}
		stack = &element{node: node, down: stack}
	}
	return stack.node
}

func (t *tokens32) PreOrder() (<-chan state32, [][]token32) {
	s, ordered := make(chan state32, 6), t.Order()
	go func() {
		var states [8]state32
		for i := range states {
			states[i].depths = make([]int32, len(ordered))
		}
		depths, state, depth := make([]int32, len(ordered)), 0, 1
		write := func(t token32, leaf bool) {
			S := states[state]
			state, S.pegRule, S.begin, S.end, S.next, S.leaf = (state+1)%8, t.pegRule, t.begin, t.end, uint32(depth), leaf
			copy(S.depths, depths)
			s <- S
		}

		states[state].token32 = ordered[0][0]
		depths[0]++
		state++
		a, b := ordered[depth-1][depths[depth-1]-1], ordered[depth][depths[depth]]
	depthFirstSearch:
		for {
			for {
				if i := depths[depth]; i > 0 {
					if c, j := ordered[depth][i-1], depths[depth-1]; a.isParentOf(c) &&
						(j < 2 || !ordered[depth-1][j-2].isParentOf(c)) {
						if c.end != b.begin {
							write(token32{pegRule: ruleIn, begin: c.end, end: b.begin}, true)
						}
						break
					}
				}

				if a.begin < b.begin {
					write(token32{pegRule: rulePre, begin: a.begin, end: b.begin}, true)
				}
				break
			}

			next := depth + 1
			if c := ordered[next][depths[next]]; c.pegRule != ruleUnknown && b.isParentOf(c) {
				write(b, false)
				depths[depth]++
				depth, a, b = next, b, c
				continue
			}

			write(b, true)
			depths[depth]++
			c, parent := ordered[depth][depths[depth]], true
			for {
				if c.pegRule != ruleUnknown && a.isParentOf(c) {
					b = c
					continue depthFirstSearch
				} else if parent && b.end != a.end {
					write(token32{pegRule: ruleSuf, begin: b.end, end: a.end}, true)
				}

				depth--
				if depth > 0 {
					a, b, c = ordered[depth-1][depths[depth-1]-1], a, ordered[depth][depths[depth]]
					parent = a.isParentOf(b)
					continue
				}

				break depthFirstSearch
			}
		}

		close(s)
	}()
	return s, ordered
}

func (t *tokens32) PrintSyntax() {
	tokens, ordered := t.PreOrder()
	max := -1
	for token := range tokens {
		if !token.leaf {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[36m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
			}
			fmt.Printf(" \x1B[36m%v\x1B[m\n", rul3s[token.pegRule])
		} else if token.begin == token.end {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[31m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
			}
			fmt.Printf(" \x1B[31m%v\x1B[m\n", rul3s[token.pegRule])
		} else {
			for c, end := token.begin, token.end; c < end; c++ {
				if i := int(c); max+1 < i {
					for j := max; j < i; j++ {
						fmt.Printf("skip %v %v\n", j, token.String())
					}
					max = i
				} else if i := int(c); i <= max {
					for j := i; j <= max; j++ {
						fmt.Printf("dupe %v %v\n", j, token.String())
					}
				} else {
					max = int(c)
				}
				fmt.Printf("%v", c)
				for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
					fmt.Printf(" \x1B[34m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
				}
				fmt.Printf(" \x1B[34m%v\x1B[m\n", rul3s[token.pegRule])
			}
			fmt.Printf("\n")
		}
	}
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	tokens, _ := t.PreOrder()
	for token := range tokens {
		for c := 0; c < int(token.next); c++ {
			fmt.Printf(" ")
		}
		fmt.Printf("\x1B[34m%v\x1B[m %v\n", rul3s[token.pegRule], strconv.Quote(string(([]rune(buffer)[token.begin:token.end]))))
	}
}

func (t *tokens32) Add(rule pegRule, begin, end, depth uint32, index int) {
	t.tree[index] = token32{pegRule: rule, begin: uint32(begin), end: uint32(end), next: uint32(depth)}
}

func (t *tokens32) Tokens() <-chan token32 {
	s := make(chan token32, 16)
	go func() {
		for _, v := range t.tree {
			s <- v.getToken32()
		}
		close(s)
	}()
	return s
}

func (t *tokens32) Error() []token32 {
	ordered := t.Order()
	length := len(ordered)
	tokens, length := make([]token32, length), length-1
	for i := range tokens {
		o := ordered[length-i]
		if len(o) > 1 {
			tokens[i] = o[len(o)-2].getToken32()
		}
	}
	return tokens
}

func (t *tokens32) Expand(index int) {
	tree := t.tree
	if index >= len(tree) {
		expanded := make([]token32, 2*len(tree))
		copy(expanded, tree)
		t.tree = expanded
	}
}

type pathParser struct {
	builder

	Buffer string
	buffer []rune
	rules  [40]func() bool
	Parse  func(rule ...int) error
	Reset  func()
	Pretty bool
	tokens32
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer []rune, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range buffer {
		if c == '\n' {
			line, symbol = line+1, 0
		} else {
			symbol++
		}
		if i == positions[j] {
			translations[positions[j]] = textPosition{line, symbol}
			for j++; j < length; j++ {
				if i != positions[j] {
					continue search
				}
			}
			break search
		}
	}

	return translations
}

type parseError struct {
	p   *pathParser
	max token32
}

func (e *parseError) Error() string {
	tokens, error := []token32{e.max}, "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.buffer, positions)
	format := "parse error near %v (line %v symbol %v - line %v symbol %v):\n%v\n"
	if e.p.Pretty {
		format = "parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n"
	}
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		error += fmt.Sprintf(format,
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			strconv.Quote(string(e.p.buffer[begin:end])))
	}

	return error
}

func (p *pathParser) PrintSyntaxTree() {
	p.tokens32.PrintSyntaxTree(p.Buffer)
}

func (p *pathParser) Highlighter() {
	p.PrintSyntax()
}

func (p *pathParser) Execute() {
	buffer, _buffer, text, begin, end := p.Buffer, p.buffer, "", 0, 0
	for token := range p.Tokens() {
		switch token.pegRule {

		case rulePegText:
			begin, end = int(token.begin), int(token.end)
			text = string(_buffer[begin:end])

		case ruleAction0:
			p.MoveTo(true)
		case ruleAction1:
			p.MoveTo(false)
		case ruleAction2:
			p.LineTo(true)
		case ruleAction3:
			p.LineTo(false)
		case ruleAction4:
			p.ClosePath()
		case ruleAction5:
			p.CubicTo(true)
		case ruleAction6:
			p.CubicTo(false)
		case ruleAction7:
			p.QuadraticTo(true)
		case ruleAction8:
			p.QuadraticTo(false)
		case ruleAction9:
			p.Coord(buffer[begin:end])

		}
	}
	_, _, _, _, _ = buffer, _buffer, text, begin, end
}

func (p *pathParser) Init() {
	p.buffer = []rune(p.Buffer)
	if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != endSymbol {
		p.buffer = append(p.buffer, endSymbol)
	}

	tree := tokens32{tree: make([]token32, math.MaxInt16)}
	var max token32
	position, depth, tokenIndex, buffer, _rules := uint32(0), uint32(0), 0, p.buffer, p.rules

	p.Parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokens32 = tree
		if matches {
			p.trim(tokenIndex)
			return nil
		}
		return &parseError{p, max}
	}

	p.Reset = func() {
		position, tokenIndex, depth = 0, 0, 0
	}

	add := func(rule pegRule, begin uint32) {
		tree.Expand(tokenIndex)
		tree.Add(rule, begin, position, depth, tokenIndex)
		tokenIndex++
		if begin != position && position > max.end {
			max = token32{rule, begin, position, depth}
		}
	}

	matchDot := func() bool {
		if buffer[position] != endSymbol {
			position++
			return true
		}
		return false
	}

	/*matchChar := func(c byte) bool {
		if buffer[position] == c {
			position++
			return true
		}
		return false
	}*/

	/*matchRange := func(lower byte, upper byte) bool {
		if c := buffer[position]; c >= lower && c <= upper {
			position++
			return true
		}
		return false
	}*/

	_rules = [...]func() bool{
		nil,
		/* 0 path <- <(wsp* movetoDrawtos? wsp* !.)> */
		func() bool {
			position0, tokenIndex0, depth0 := position, tokenIndex, depth
			{
				position1 := position
				depth++
			l2:
				{
					position3, tokenIndex3, depth3 := position, tokenIndex, depth
					if !_rules[rulewsp]() {
						goto l3
					}
					goto l2
				l3:
					position, tokenIndex, depth = position3, tokenIndex3, depth3
				}
				{
					position4, tokenIndex4, depth4 := position, tokenIndex, depth
					if !_rules[rulemovetoDrawtos]() {
						goto l4
					}
					goto l5
				l4:
					position, tokenIndex, depth = position4, tokenIndex4, depth4
				}
			l5:
			l6:
				{
					position7, tokenIndex7, depth7 := position, tokenIndex, depth
					if !_rules[rulewsp]() {
						goto l7
					}
					goto l6
				l7:
					position, tokenIndex, depth = position7, tokenIndex7, depth7
				}
				{
					position8, tokenIndex8, depth8 := position, tokenIndex, depth
					if !matchDot() {
						goto l8
					}
					goto l0
				l8:
					position, tokenIndex, depth = position8, tokenIndex8, depth8
				}
				depth--
				add(rulepath, position1)
			}
			return true
		l0:
			position, tokenIndex, depth = position0, tokenIndex0, depth0
			return false
		},
		/* 1 movetoDrawtos <- <((movetoDrawto wsp* movetoDrawtos) / movetoDrawto)> */
		func() bool {
			position9, tokenIndex9, depth9 := position, tokenIndex, depth
			{
				position10 := position
				depth++
				{
					position11, tokenIndex11, depth11 := position, tokenIndex, depth
					if !_rules[rulemovetoDrawto]() {
						goto l12
					}
				l13:
					{
						position14, tokenIndex14, depth14 := position, tokenIndex, depth
						if !_rules[rulewsp]() {
							goto l14
						}
						goto l13
					l14:
						position, tokenIndex, depth = position14, tokenIndex14, depth14
					}
					if !_rules[rulemovetoDrawtos]() {
						goto l12
					}
					goto l11
				l12:
					position, tokenIndex, depth = position11, tokenIndex11, depth11
					if !_rules[rulemovetoDrawto]() {
						goto l9
					}
				}
			l11:
				depth--
				add(rulemovetoDrawtos, position10)
			}
			return true
		l9:
			position, tokenIndex, depth = position9, tokenIndex9, depth9
			return false
		},
		/* 2 movetoDrawto <- <(moveto wsp* drawtos?)> */
		func() bool {
			position15, tokenIndex15, depth15 := position, tokenIndex, depth
			{
				position16 := position
				depth++
				if !_rules[rulemoveto]() {
					goto l15
				}
			l17:
				{
					position18, tokenIndex18, depth18 := position, tokenIndex, depth
					if !_rules[rulewsp]() {
						goto l18
					}
					goto l17
				l18:
					position, tokenIndex, depth = position18, tokenIndex18, depth18
				}
				{
					position19, tokenIndex19, depth19 := position, tokenIndex, depth
					if !_rules[ruledrawtos]() {
						goto l19
					}
					goto l20
				l19:
					position, tokenIndex, depth = position19, tokenIndex19, depth19
				}
			l20:
				depth--
				add(rulemovetoDrawto, position16)
			}
			return true
		l15:
			position, tokenIndex, depth = position15, tokenIndex15, depth15
			return false
		},
		/* 3 drawtos <- <((drawto wsp* drawtos) / drawto)> */
		func() bool {
			position21, tokenIndex21, depth21 := position, tokenIndex, depth
			{
				position22 := position
				depth++
				{
					position23, tokenIndex23, depth23 := position, tokenIndex, depth
					if !_rules[ruledrawto]() {
						goto l24
					}
				l25:
					{
						position26, tokenIndex26, depth26 := position, tokenIndex, depth
						if !_rules[rulewsp]() {
							goto l26
						}
						goto l25
					l26:
						position, tokenIndex, depth = position26, tokenIndex26, depth26
					}
					if !_rules[ruledrawtos]() {
						goto l24
					}
					goto l23
				l24:
					position, tokenIndex, depth = position23, tokenIndex23, depth23
					if !_rules[ruledrawto]() {
						goto l21
					}
				}
			l23:
				depth--
				add(ruledrawtos, position22)
			}
			return true
		l21:
			position, tokenIndex, depth = position21, tokenIndex21, depth21
			return false
		},
		/* 4 drawto <- <(closepath / cubicto / quadto / lineto)> */
		func() bool {
			position27, tokenIndex27, depth27 := position, tokenIndex, depth
			{
				position28 := position
				depth++
				{
					position29, tokenIndex29, depth29 := position, tokenIndex, depth
					if !_rules[ruleclosepath]() {
						goto l30
					}
					goto l29
				l30:
					position, tokenIndex, depth = position29, tokenIndex29, depth29
					if !_rules[rulecubicto]() {
						goto l31
					}
					goto l29
				l31:
					position, tokenIndex, depth = position29, tokenIndex29, depth29
					if !_rules[rulequadto]() {
						goto l32
					}
					goto l29
				l32:
					position, tokenIndex, depth = position29, tokenIndex29, depth29
					if !_rules[rulelineto]() {
						goto l27
					}
				}
			l29:
				depth--
				add(ruledrawto, position28)
			}
			return true
		l27:
			position, tokenIndex, depth = position27, tokenIndex27, depth27
			return false
		},
		/* 5 moveto <- <(('M' wsp* movetoArgs Action0) / ('m' wsp* movetoArgs Action1))> */
		func() bool {
			position33, tokenIndex33, depth33 := position, tokenIndex, depth
			{
				position34 := position
				depth++
				{
					position35, tokenIndex35, depth35 := position, tokenIndex, depth
					if buffer[position] != rune('M') {
						goto l36
					}
					position++
				l37:
					{
						position38, tokenIndex38, depth38 := position, tokenIndex, depth
						if !_rules[rulewsp]() {
							goto l38
						}
						goto l37
					l38:
						position, tokenIndex, depth = position38, tokenIndex38, depth38
					}
					if !_rules[rulemovetoArgs]() {
						goto l36
					}
					if !_rules[ruleAction0]() {
						goto l36
					}
					goto l35
				l36:
					position, tokenIndex, depth = position35, tokenIndex35, depth35
					if buffer[position] != rune('m') {
						goto l33
					}
					position++
				l39:
					{
						position40, tokenIndex40, depth40 := position, tokenIndex, depth
						if !_rules[rulewsp]() {
							goto l40
						}
						goto l39
					l40:
						position, tokenIndex, depth = position40, tokenIndex40, depth40
					}
					if !_rules[rulemovetoArgs]() {
						goto l33
					}
					if !_rules[ruleAction1]() {
						goto l33
					}
				}
			l35:
				depth--
				add(rulemoveto, position34)
			}
			return true
		l33:
			position, tokenIndex, depth = position33, tokenIndex33, depth33
			return false
		},
		/* 6 movetoArgs <- <((pair comma? linetoArgs) / pair)> */
		func() bool {
			position41, tokenIndex41, depth41 := position, tokenIndex, depth
			{
				position42 := position
				depth++
				{
					position43, tokenIndex43, depth43 := position, tokenIndex, depth
					if !_rules[rulepair]() {
						goto l44
					}
					{
						position45, tokenIndex45, depth45 := position, tokenIndex, depth
						if !_rules[rulecomma]() {
							goto l45
						}
						goto l46
					l45:
						position, tokenIndex, depth = position45, tokenIndex45, depth45
					}
				l46:
					if !_rules[rulelinetoArgs]() {
						goto l44
					}
					goto l43
				l44:
					position, tokenIndex, depth = position43, tokenIndex43, depth43
					if !_rules[rulepair]() {
						goto l41
					}
				}
			l43:
				depth--
				add(rulemovetoArgs, position42)
			}
			return true
		l41:
			position, tokenIndex, depth = position41, tokenIndex41, depth41
			return false
		},
		/* 7 lineto <- <(('L' wsp* linetoArgs Action2) / ('l' wsp* linetoArgs Action3))> */
		func() bool {
			position47, tokenIndex47, depth47 := position, tokenIndex, depth
			{
				position48 := position
				depth++
				{
					position49, tokenIndex49, depth49 := position, tokenIndex, depth
					if buffer[position] != rune('L') {
						goto l50
					}
					position++
				l51:
					{
						position52, tokenIndex52, depth52 := position, tokenIndex, depth
						if !_rules[rulewsp]() {
							goto l52
						}
						goto l51
					l52:
						position, tokenIndex, depth = position52, tokenIndex52, depth52
					}
					if !_rules[rulelinetoArgs]() {
						goto l50
					}
					if !_rules[ruleAction2]() {
						goto l50
					}
					goto l49
				l50:
					position, tokenIndex, depth = position49, tokenIndex49, depth49
					if buffer[position] != rune('l') {
						goto l47
					}
					position++
				l53:
					{
						position54, tokenIndex54, depth54 := position, tokenIndex, depth
						if !_rules[rulewsp]() {
							goto l54
						}
						goto l53
					l54:
						position, tokenIndex, depth = position54, tokenIndex54, depth54
					}
					if !_rules[rulelinetoArgs]() {
						goto l47
					}
					if !_rules[ruleAction3]() {
						goto l47
					}
				}
			l49:
				depth--
				add(rulelineto, position48)
			}
			return true
		l47:
			position, tokenIndex, depth = position47, tokenIndex47, depth47
			return false
		},
		/* 8 closepath <- <(('z' / 'Z') Action4)> */
		func() bool {
			position55, tokenIndex55, depth55 := position, tokenIndex, depth
			{
				position56 := position
				depth++
				{
					position57, tokenIndex57, depth57 := position, tokenIndex, depth
					if buffer[position] != rune('z') {
						goto l58
					}
					position++
					goto l57
				l58:
					position, tokenIndex, depth = position57, tokenIndex57, depth57
					if buffer[position] != rune('Z') {
						goto l55
					}
					position++
				}
			l57:
				if !_rules[ruleAction4]() {
					goto l55
				}
				depth--
				add(ruleclosepath, position56)
			}
			return true
		l55:
			position, tokenIndex, depth = position55, tokenIndex55, depth55
			return false
		},
		/* 9 linetoArgs <- <((pair comma? linetoArgs) / pair)> */
		func() bool {
			position59, tokenIndex59, depth59 := position, tokenIndex, depth
			{
				position60 := position
				depth++
				{
					position61, tokenIndex61, depth61 := position, tokenIndex, depth
					if !_rules[rulepair]() {
						goto l62
					}
					{
						position63, tokenIndex63, depth63 := position, tokenIndex, depth
						if !_rules[rulecomma]() {
							goto l63
						}
						goto l64
					l63:
						position, tokenIndex, depth = position63, tokenIndex63, depth63
					}
				l64:
					if !_rules[rulelinetoArgs]() {
						goto l62
					}
					goto l61
				l62:
					position, tokenIndex, depth = position61, tokenIndex61, depth61
					if !_rules[rulepair]() {
						goto l59
					}
				}
			l61:
				depth--
				add(rulelinetoArgs, position60)
			}
			return true
		l59:
			position, tokenIndex, depth = position59, tokenIndex59, depth59
			return false
		},
		/* 10 cubicto <- <(('C' wsp* cubictoArgs Action5) / ('c' wsp* cubictoArgs Action6))> */
		func() bool {
			position65, tokenIndex65, depth65 := position, tokenIndex, depth
			{
				position66 := position
				depth++
				{
					position67, tokenIndex67, depth67 := position, tokenIndex, depth
					if buffer[position] != rune('C') {
						goto l68
					}
					position++
				l69:
					{
						position70, tokenIndex70, depth70 := position, tokenIndex, depth
						if !_rules[rulewsp]() {
							goto l70
						}
						goto l69
					l70:
						position, tokenIndex, depth = position70, tokenIndex70, depth70
					}
					if !_rules[rulecubictoArgs]() {
						goto l68
					}
					if !_rules[ruleAction5]() {
						goto l68
					}
					goto l67
				l68:
					position, tokenIndex, depth = position67, tokenIndex67, depth67
					if buffer[position] != rune('c') {
						goto l65
					}
					position++
				l71:
					{
						position72, tokenIndex72, depth72 := position, tokenIndex, depth
						if !_rules[rulewsp]() {
							goto l72
						}
						goto l71
					l72:
						position, tokenIndex, depth = position72, tokenIndex72, depth72
					}
					if !_rules[rulecubictoArgs]() {
						goto l65
					}
					if !_rules[ruleAction6]() {
						goto l65
					}
				}
			l67:
				depth--
				add(rulecubicto, position66)
			}
			return true
		l65:
			position, tokenIndex, depth = position65, tokenIndex65, depth65
			return false
		},
		/* 11 cubictoArgs <- <((cubictoArg comma? cubictoArgs) / cubictoArg)> */
		func() bool {
			position73, tokenIndex73, depth73 := position, tokenIndex, depth
			{
				position74 := position
				depth++
				{
					position75, tokenIndex75, depth75 := position, tokenIndex, depth
					if !_rules[rulecubictoArg]() {
						goto l76
					}
					{
						position77, tokenIndex77, depth77 := position, tokenIndex, depth
						if !_rules[rulecomma]() {
							goto l77
						}
						goto l78
					l77:
						position, tokenIndex, depth = position77, tokenIndex77, depth77
					}
				l78:
					if !_rules[rulecubictoArgs]() {
						goto l76
					}
					goto l75
				l76:
					position, tokenIndex, depth = position75, tokenIndex75, depth75
					if !_rules[rulecubictoArg]() {
						goto l73
					}
				}
			l75:
				depth--
				add(rulecubictoArgs, position74)
			}
			return true
		l73:
			position, tokenIndex, depth = position73, tokenIndex73, depth73
			return false
		},
		/* 12 cubictoArg <- <(pair comma? pair comma? pair)> */
		func() bool {
			position79, tokenIndex79, depth79 := position, tokenIndex, depth
			{
				position80 := position
				depth++
				if !_rules[rulepair]() {
					goto l79
				}
				{
					position81, tokenIndex81, depth81 := position, tokenIndex, depth
					if !_rules[rulecomma]() {
						goto l81
					}
					goto l82
				l81:
					position, tokenIndex, depth = position81, tokenIndex81, depth81
				}
			l82:
				if !_rules[rulepair]() {
					goto l79
				}
				{
					position83, tokenIndex83, depth83 := position, tokenIndex, depth
					if !_rules[rulecomma]() {
						goto l83
					}
					goto l84
				l83:
					position, tokenIndex, depth = position83, tokenIndex83, depth83
				}
			l84:
				if !_rules[rulepair]() {
					goto l79
				}
				depth--
				add(rulecubictoArg, position80)
			}
			return true
		l79:
			position, tokenIndex, depth = position79, tokenIndex79, depth79
			return false
		},
		/* 13 quadto <- <(('Q' wsp* quadtoArgs Action7) / ('q' wsp* quadtoArgs Action8))> */
		func() bool {
			position85, tokenIndex85, depth85 := position, tokenIndex, depth
			{
				position86 := position
				depth++
				{
					position87, tokenIndex87, depth87 := position, tokenIndex, depth
					if buffer[position] != rune('Q') {
						goto l88
					}
					position++
				l89:
					{
						position90, tokenIndex90, depth90 := position, tokenIndex, depth
						if !_rules[rulewsp]() {
							goto l90
						}
						goto l89
					l90:
						position, tokenIndex, depth = position90, tokenIndex90, depth90
					}
					if !_rules[rulequadtoArgs]() {
						goto l88
					}
					if !_rules[ruleAction7]() {
						goto l88
					}
					goto l87
				l88:
					position, tokenIndex, depth = position87, tokenIndex87, depth87
					if buffer[position] != rune('q') {
						goto l85
					}
					position++
				l91:
					{
						position92, tokenIndex92, depth92 := position, tokenIndex, depth
						if !_rules[rulewsp]() {
							goto l92
						}
						goto l91
					l92:
						position, tokenIndex, depth = position92, tokenIndex92, depth92
					}
					if !_rules[rulequadtoArgs]() {
						goto l85
					}
					if !_rules[ruleAction8]() {
						goto l85
					}
				}
			l87:
				depth--
				add(rulequadto, position86)
			}
			return true
		l85:
			position, tokenIndex, depth = position85, tokenIndex85, depth85
			return false
		},
		/* 14 quadtoArgs <- <((quadtoArg comma? quadtoArgs) / quadtoArg)> */
		func() bool {
			position93, tokenIndex93, depth93 := position, tokenIndex, depth
			{
				position94 := position
				depth++
				{
					position95, tokenIndex95, depth95 := position, tokenIndex, depth
					if !_rules[rulequadtoArg]() {
						goto l96
					}
					{
						position97, tokenIndex97, depth97 := position, tokenIndex, depth
						if !_rules[rulecomma]() {
							goto l97
						}
						goto l98
					l97:
						position, tokenIndex, depth = position97, tokenIndex97, depth97
					}
				l98:
					if !_rules[rulequadtoArgs]() {
						goto l96
					}
					goto l95
				l96:
					position, tokenIndex, depth = position95, tokenIndex95, depth95
					if !_rules[rulequadtoArg]() {
						goto l93
					}
				}
			l95:
				depth--
				add(rulequadtoArgs, position94)
			}
			return true
		l93:
			position, tokenIndex, depth = position93, tokenIndex93, depth93
			return false
		},
		/* 15 quadtoArg <- <(pair comma? pair)> */
		func() bool {
			position99, tokenIndex99, depth99 := position, tokenIndex, depth
			{
				position100 := position
				depth++
				if !_rules[rulepair]() {
					goto l99
				}
				{
					position101, tokenIndex101, depth101 := position, tokenIndex, depth
					if !_rules[rulecomma]() {
						goto l101
					}
					goto l102
				l101:
					position, tokenIndex, depth = position101, tokenIndex101, depth101
				}
			l102:
				if !_rules[rulepair]() {
					goto l99
				}
				depth--
				add(rulequadtoArg, position100)
			}
			return true
		l99:
			position, tokenIndex, depth = position99, tokenIndex99, depth99
			return false
		},
		/* 16 digits <- <[0-9]+> */
		func() bool {
			position103, tokenIndex103, depth103 := position, tokenIndex, depth
			{
				position104 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l103
				}
				position++
			l105:
				{
					position106, tokenIndex106, depth106 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l106
					}
					position++
					goto l105
				l106:
					position, tokenIndex, depth = position106, tokenIndex106, depth106
				}
				depth--
				add(ruledigits, position104)
			}
			return true
		l103:
			position, tokenIndex, depth = position103, tokenIndex103, depth103
			return false
		},
		/* 17 pair <- <(coord comma? coord)> */
		func() bool {
			position107, tokenIndex107, depth107 := position, tokenIndex, depth
			{
				position108 := position
				depth++
				if !_rules[rulecoord]() {
					goto l107
				}
				{
					position109, tokenIndex109, depth109 := position, tokenIndex, depth
					if !_rules[rulecomma]() {
						goto l109
					}
					goto l110
				l109:
					position, tokenIndex, depth = position109, tokenIndex109, depth109
				}
			l110:
				if !_rules[rulecoord]() {
					goto l107
				}
				depth--
				add(rulepair, position108)
			}
			return true
		l107:
			position, tokenIndex, depth = position107, tokenIndex107, depth107
			return false
		},
		/* 18 coord <- <(<number> Action9)> */
		func() bool {
			position111, tokenIndex111, depth111 := position, tokenIndex, depth
			{
				position112 := position
				depth++
				{
					position113 := position
					depth++
					if !_rules[rulenumber]() {
						goto l111
					}
					depth--
					add(rulePegText, position113)
				}
				if !_rules[ruleAction9]() {
					goto l111
				}
				depth--
				add(rulecoord, position112)
			}
			return true
		l111:
			position, tokenIndex, depth = position111, tokenIndex111, depth111
			return false
		},
		/* 19 number <- <((sign? float) / (sign? integer))> */
		func() bool {
			position114, tokenIndex114, depth114 := position, tokenIndex, depth
			{
				position115 := position
				depth++
				{
					position116, tokenIndex116, depth116 := position, tokenIndex, depth
					{
						position118, tokenIndex118, depth118 := position, tokenIndex, depth
						if !_rules[rulesign]() {
							goto l118
						}
						goto l119
					l118:
						position, tokenIndex, depth = position118, tokenIndex118, depth118
					}
				l119:
					if !_rules[rulefloat]() {
						goto l117
					}
					goto l116
				l117:
					position, tokenIndex, depth = position116, tokenIndex116, depth116
					{
						position120, tokenIndex120, depth120 := position, tokenIndex, depth
						if !_rules[rulesign]() {
							goto l120
						}
						goto l121
					l120:
						position, tokenIndex, depth = position120, tokenIndex120, depth120
					}
				l121:
					if !_rules[ruleinteger]() {
						goto l114
					}
				}
			l116:
				depth--
				add(rulenumber, position115)
			}
			return true
		l114:
			position, tokenIndex, depth = position114, tokenIndex114, depth114
			return false
		},
		/* 20 nonneg <- <(float / integer)> */
		nil,
		/* 21 comma <- <((wsp+ ','? wsp*) / (',' wsp*))> */
		func() bool {
			position123, tokenIndex123, depth123 := position, tokenIndex, depth
			{
				position124 := position
				depth++
				{
					position125, tokenIndex125, depth125 := position, tokenIndex, depth
					if !_rules[rulewsp]() {
						goto l126
					}
				l127:
					{
						position128, tokenIndex128, depth128 := position, tokenIndex, depth
						if !_rules[rulewsp]() {
							goto l128
						}
						goto l127
					l128:
						position, tokenIndex, depth = position128, tokenIndex128, depth128
					}
					{
						position129, tokenIndex129, depth129 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l129
						}
						position++
						goto l130
					l129:
						position, tokenIndex, depth = position129, tokenIndex129, depth129
					}
				l130:
				l131:
					{
						position132, tokenIndex132, depth132 := position, tokenIndex, depth
						if !_rules[rulewsp]() {
							goto l132
						}
						goto l131
					l132:
						position, tokenIndex, depth = position132, tokenIndex132, depth132
					}
					goto l125
				l126:
					position, tokenIndex, depth = position125, tokenIndex125, depth125
					if buffer[position] != rune(',') {
						goto l123
					}
					position++
				l133:
					{
						position134, tokenIndex134, depth134 := position, tokenIndex, depth
						if !_rules[rulewsp]() {
							goto l134
						}
						goto l133
					l134:
						position, tokenIndex, depth = position134, tokenIndex134, depth134
					}
				}
			l125:
				depth--
				add(rulecomma, position124)
			}
			return true
		l123:
			position, tokenIndex, depth = position123, tokenIndex123, depth123
			return false
		},
		/* 22 integer <- <digits> */
		func() bool {
			position135, tokenIndex135, depth135 := position, tokenIndex, depth
			{
				position136 := position
				depth++
				if !_rules[ruledigits]() {
					goto l135
				}
				depth--
				add(ruleinteger, position136)
			}
			return true
		l135:
			position, tokenIndex, depth = position135, tokenIndex135, depth135
			return false
		},
		/* 23 float <- <((fract exponent?) / (digits exponent))> */
		func() bool {
			position137, tokenIndex137, depth137 := position, tokenIndex, depth
			{
				position138 := position
				depth++
				{
					position139, tokenIndex139, depth139 := position, tokenIndex, depth
					if !_rules[rulefract]() {
						goto l140
					}
					{
						position141, tokenIndex141, depth141 := position, tokenIndex, depth
						if !_rules[ruleexponent]() {
							goto l141
						}
						goto l142
					l141:
						position, tokenIndex, depth = position141, tokenIndex141, depth141
					}
				l142:
					goto l139
				l140:
					position, tokenIndex, depth = position139, tokenIndex139, depth139
					if !_rules[ruledigits]() {
						goto l137
					}
					if !_rules[ruleexponent]() {
						goto l137
					}
				}
			l139:
				depth--
				add(rulefloat, position138)
			}
			return true
		l137:
			position, tokenIndex, depth = position137, tokenIndex137, depth137
			return false
		},
		/* 24 fract <- <((digits? '.' digits) / (digits '.'))> */
		func() bool {
			position143, tokenIndex143, depth143 := position, tokenIndex, depth
			{
				position144 := position
				depth++
				{
					position145, tokenIndex145, depth145 := position, tokenIndex, depth
					{
						position147, tokenIndex147, depth147 := position, tokenIndex, depth
						if !_rules[ruledigits]() {
							goto l147
						}
						goto l148
					l147:
						position, tokenIndex, depth = position147, tokenIndex147, depth147
					}
				l148:
					if buffer[position] != rune('.') {
						goto l146
					}
					position++
					if !_rules[ruledigits]() {
						goto l146
					}
					goto l145
				l146:
					position, tokenIndex, depth = position145, tokenIndex145, depth145
					if !_rules[ruledigits]() {
						goto l143
					}
					if buffer[position] != rune('.') {
						goto l143
					}
					position++
				}
			l145:
				depth--
				add(rulefract, position144)
			}
			return true
		l143:
			position, tokenIndex, depth = position143, tokenIndex143, depth143
			return false
		},
		/* 25 exponent <- <(('e' / 'E') sign? digits)> */
		func() bool {
			position149, tokenIndex149, depth149 := position, tokenIndex, depth
			{
				position150 := position
				depth++
				{
					position151, tokenIndex151, depth151 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l152
					}
					position++
					goto l151
				l152:
					position, tokenIndex, depth = position151, tokenIndex151, depth151
					if buffer[position] != rune('E') {
						goto l149
					}
					position++
				}
			l151:
				{
					position153, tokenIndex153, depth153 := position, tokenIndex, depth
					if !_rules[rulesign]() {
						goto l153
					}
					goto l154
				l153:
					position, tokenIndex, depth = position153, tokenIndex153, depth153
				}
			l154:
				if !_rules[ruledigits]() {
					goto l149
				}
				depth--
				add(ruleexponent, position150)
			}
			return true
		l149:
			position, tokenIndex, depth = position149, tokenIndex149, depth149
			return false
		},
		/* 26 wsp <- <(' ' / '\t' / '\r' / '\n')> */
		func() bool {
			position155, tokenIndex155, depth155 := position, tokenIndex, depth
			{
				position156 := position
				depth++
				{
					position157, tokenIndex157, depth157 := position, tokenIndex, depth
					if buffer[position] != rune(' ') {
						goto l158
					}
					position++
					goto l157
				l158:
					position, tokenIndex, depth = position157, tokenIndex157, depth157
					if buffer[position] != rune('\t') {
						goto l159
					}
					position++
					goto l157
				l159:
					position, tokenIndex, depth = position157, tokenIndex157, depth157
					if buffer[position] != rune('\r') {
						goto l160
					}
					position++
					goto l157
				l160:
					position, tokenIndex, depth = position157, tokenIndex157, depth157
					if buffer[position] != rune('\n') {
						goto l155
					}
					position++
				}
			l157:
				depth--
				add(rulewsp, position156)
			}
			return true
		l155:
			position, tokenIndex, depth = position155, tokenIndex155, depth155
			return false
		},
		/* 27 sign <- <('-' / '+')> */
		func() bool {
			position161, tokenIndex161, depth161 := position, tokenIndex, depth
			{
				position162 := position
				depth++
				{
					position163, tokenIndex163, depth163 := position, tokenIndex, depth
					if buffer[position] != rune('-') {
						goto l164
					}
					position++
					goto l163
				l164:
					position, tokenIndex, depth = position163, tokenIndex163, depth163
					if buffer[position] != rune('+') {
						goto l161
					}
					position++
				}
			l163:
				depth--
				add(rulesign, position162)
			}
			return true
		l161:
			position, tokenIndex, depth = position161, tokenIndex161, depth161
			return false
		},
		/* 29 Action0 <- <{ p.MoveTo(true) }> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 30 Action1 <- <{ p.MoveTo(false) }> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 31 Action2 <- <{ p.LineTo(true) }> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 32 Action3 <- <{ p.LineTo(false) }> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 33 Action4 <- <{ p.ClosePath() }> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 34 Action5 <- <{ p.CubicTo(true) }> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 35 Action6 <- <{ p.CubicTo(false) }> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 36 Action7 <- <{ p.QuadraticTo(true) }> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		/* 37 Action8 <- <{ p.QuadraticTo(false) }> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		nil,
		/* 39 Action9 <- <{ p.Coord(buffer[begin:end]) }> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
	}
	p.rules = _rules
}
