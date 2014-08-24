//line rule.y:2
package filter

import __yyfmt__ "fmt"

//line rule.y:2
import (
	"errors"
	"fmt"
	"io"
	"sync"
)

var g_grammer *Grammer
var g_mutex sync.Mutex

//line rule.y:8
type yySymType struct {
	yys     int
	grammer *Grammer
	expr    *Expr
	term    *Term
	factor  *Factor
	list    *List
	fun     *Func
	str     string
	dval    float64
	fn      int
}

const COMMA = 57346
const LPAREN = 57347
const RPAREN = 57348
const LAND = 57349
const LOR = 57350
const GET = 57351
const DEFAULT = 57352
const VAR = 57353
const STR = 57354
const NUM = 57355
const CMP = 57356
const CONTAIN = 57357
const FUNC = 57358

var yyToknames = []string{
	"COMMA",
	"LPAREN",
	"RPAREN",
	"LAND",
	"LOR",
	"GET",
	"DEFAULT",
	"VAR",
	"STR",
	"NUM",
	"CMP",
	"CONTAIN",
	"FUNC",
}
var yyStatenames = []string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line rule.y:57

/*
   parser handle
*/
type Parser struct {
	grammer *Grammer
}

/*
   analyze input rule script and generate parser handle
*/
func NewParser(in io.Reader) (h *Parser, err error) {
	defer func() {
		if e := recover(); e != nil {
			h, err = nil, errors.New(fmt.Sprint(e))
		}
	}()
	h = new(Parser)
	lex := NewLexer(in)
	yyParse(lex)
	g_mutex.Lock()
	h.grammer = g_grammer
	g_mutex.Unlock()
	if h.grammer == nil {
		return h, errors.New("invalid rule")
	}
	return h, err
}

/*
   get parse result
   symlist is created by calling QueryToSymlist() or JsonToSymlist() API
*/
func (h *Parser) Parse(symlist *SymList) (ret int, err error) {
	defer func() {
		if e := recover(); e != nil {
			ret, err = -1, errors.New(fmt.Sprint(e))
		}
	}()
	ret, err = EvalGrammer(h.grammer, symlist)
	return
}

//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 20
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 61

var yyAct = []int{

	6, 29, 5, 2, 19, 18, 7, 14, 15, 16,
	13, 4, 8, 9, 10, 25, 22, 12, 23, 24,
	27, 17, 31, 37, 35, 7, 32, 31, 34, 33,
	4, 8, 9, 10, 30, 3, 12, 31, 38, 8,
	9, 10, 7, 20, 12, 28, 15, 16, 8, 9,
	10, 26, 21, 12, 8, 9, 10, 36, 1, 12,
	11,
}
var yyPact = []int{

	20, -1000, -1000, 1, 12, -1000, -10, 37, -1000, -1000,
	-1000, -1000, 47, 3, -1000, 37, 37, 2, 46, 43,
	39, 28, 20, -1000, -1000, 20, 43, -1000, -1000, 18,
	-1000, 53, -1000, -1000, 17, -1000, 43, -1000, -1000,
}
var yyPgo = []int{

	0, 3, 35, 2, 0, 60, 1, 58,
}
var yyR1 = []int{

	0, 7, 1, 1, 1, 1, 2, 2, 2, 3,
	3, 3, 4, 4, 4, 4, 6, 6, 5, 5,
}
var yyR2 = []int{

	0, 1, 4, 4, 2, 0, 3, 3, 1, 5,
	3, 3, 1, 1, 1, 1, 1, 3, 4, 3,
}
var yyChk = []int{

	-1000, -7, -1, -2, 10, -3, -4, 5, 11, 12,
	13, -5, 16, 9, -1, 7, 8, 9, 15, 14,
	-2, 5, 13, -3, -3, 13, 5, -4, 6, -6,
	6, -4, -1, -1, -6, 6, 4, 6, -6,
}
var yyDef = []int{

	5, -2, 1, 5, 0, 8, 0, 0, 12, 13,
	14, 15, 0, 0, 4, 0, 0, 0, 0, 0,
	0, 0, 5, 6, 7, 5, 0, 10, 11, 0,
	19, 16, 2, 3, 0, 18, 0, 9, 17,
}
var yyTok1 = []int{

	1,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16,
}
var yyTok3 = []int{
	0,
}

//line yaccpar:1

/*	parser for yacc output	*/

var yyDebug = 0

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

const yyFlag = -1000

func yyTokname(c int) string {
	// 4 is TOKSTART above
	if c >= 4 && c-4 < len(yyToknames) {
		if yyToknames[c-4] != "" {
			return yyToknames[c-4]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yylex1(lex yyLexer, lval *yySymType) int {
	c := 0
	char := lex.Lex(lval)
	if char <= 0 {
		c = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		c = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			c = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		c = yyTok3[i+0]
		if c == char {
			c = yyTok3[i+1]
			goto out
		}
	}

out:
	if c == 0 {
		c = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(c), uint(char))
	}
	return c
}

func yyParse(yylex yyLexer) int {
	var yyn int
	var yylval yySymType
	var yyVAL yySymType
	yyS := make([]yySymType, yyMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yychar := -1
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yychar), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yychar < 0 {
		yychar = yylex1(yylex, &yylval)
	}
	yyn += yychar
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yychar { /* valid shift */
		yychar = -1
		yyVAL = yylval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yychar < 0 {
			yychar = yylex1(yylex, &yylval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yychar {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error("syntax error")
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yychar))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yychar))
			}
			if yychar == yyEofCode {
				goto ret1
			}
			yychar = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		//line rule.y:32
		{
			g_mutex.Lock()
			g_grammer = yyS[yypt-0].grammer
			g_mutex.Unlock()
		}
	case 2:
		//line rule.y:33
		{
			var err error
			if yyVAL.grammer, err = NewGrammer(EGET, yyS[yypt-3].expr, yyS[yypt-1].dval, yyS[yypt-0].grammer); err != nil {
				panic(err)
			}
		}
	case 3:
		//line rule.y:34
		{
			var err error
			if yyVAL.grammer, err = NewGrammer(DGET, nil, yyS[yypt-1].dval, yyS[yypt-0].grammer); err != nil {
				panic(err)
			}
		}
	case 4:
		//line rule.y:35
		{
			var err error
			if yyVAL.grammer, err = NewGrammer(EEXPR, yyS[yypt-1].expr, 0, yyS[yypt-0].grammer); err != nil {
				panic(err)
			}
		}
	case 5:
		//line rule.y:36
		{
			yyVAL.grammer = nil
		}
	case 6:
		//line rule.y:38
		{
			var err error
			if yyVAL.expr, err = NewExpr(AND, yyS[yypt-2].expr, yyS[yypt-0].term); err != nil {
				panic(err)
			}
		}
	case 7:
		//line rule.y:39
		{
			var err error
			if yyVAL.expr, err = NewExpr(OR, yyS[yypt-2].expr, yyS[yypt-0].term); err != nil {
				panic(err)
			}
		}
	case 8:
		//line rule.y:40
		{
			var err error
			if yyVAL.expr, err = NewExpr(TERM, nil, yyS[yypt-0].term); err != nil {
				panic(err)
			}
		}
	case 9:
		//line rule.y:42
		{
			var err error
			if yyVAL.term, err = NewTerm(TKind_t(yyS[yypt-3].fn), yyS[yypt-4].factor, yyS[yypt-1].list, nil, nil); err != nil {
				panic(err)
			}
		}
	case 10:
		//line rule.y:43
		{
			var err error
			if yyVAL.term, err = NewTerm(TKind_t(yyS[yypt-1].fn), yyS[yypt-2].factor, nil, yyS[yypt-0].factor, nil); err != nil {
				panic(err)
			}
		}
	case 11:
		//line rule.y:44
		{
			var err error
			if yyVAL.term, err = NewTerm(EXPR, nil, nil, nil, yyS[yypt-1].expr); err != nil {
				panic(err)
			}
		}
	case 12:
		//line rule.y:46
		{
			var err error
			if yyVAL.factor, err = NewFactor(VARIABLE, 0, "", yyS[yypt-0].str, nil); err != nil {
				panic(err)
			}
		}
	case 13:
		//line rule.y:47
		{
			var err error
			if yyVAL.factor, err = NewFactor(STRING, 0, yyS[yypt-0].str, "", nil); err != nil {
				panic(err)
			}
		}
	case 14:
		//line rule.y:48
		{
			var err error
			if yyVAL.factor, err = NewFactor(DOUBLE, yyS[yypt-0].dval, "", "", nil); err != nil {
				panic(err)
			}
		}
	case 15:
		//line rule.y:49
		{
			var err error
			if yyVAL.factor, err = NewFactor(FUNCTION, 0, "", "", yyS[yypt-0].fun); err != nil {
				panic(err)
			}
		}
	case 16:
		//line rule.y:51
		{
			var err error
			if yyVAL.list, err = NewList(yyS[yypt-0].factor, nil); err != nil {
				panic(err)
			}
		}
	case 17:
		//line rule.y:52
		{
			var err error
			if yyVAL.list, err = NewList(yyS[yypt-2].factor, yyS[yypt-0].list); err != nil {
				panic(err)
			}
		}
	case 18:
		//line rule.y:54
		{
			var err error
			if yyVAL.fun, err = NewFunc(FnKind_t(yyS[yypt-3].fn), yyS[yypt-1].list); err != nil {
				panic(err)
			}
		}
	case 19:
		//line rule.y:55
		{
			var err error
			if yyVAL.fun, err = NewFunc(FnKind_t(yyS[yypt-2].fn), nil); err != nil {
				panic(err)
			}
		}
	}
	goto yystack /* stack new state and value */
}
