//line tacky.y:2
package main

import __yyfmt__ "fmt"

//line tacky.y:2
//import "fmt"

//line tacky.y:6
type yySymType struct {
	yys  int
	s    string
	expr *Expr
}

const DEF_FORM = 57346
const ASSIGN = 57347
const ID = 57348
const MONEY = 57349
const FRAC = 57350
const XREF = 57351
const FUNC = 57352

var yyToknames = []string{
	"DEF_FORM",
	"ASSIGN",
	"ID",
	"MONEY",
	"FRAC",
	"XREF",
	"FUNC",
	" +",
	" -",
	" *",
	" /",
}
var yyStatenames = []string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line tacky.y:45
func cast(y yyLexer) *Tacky { return y.(*Lexer).p }

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

const yyLast = 43

var yyAct = []int{

	6, 19, 9, 8, 13, 11, 12, 26, 27, 7,
	2, 18, 10, 20, 1, 21, 22, 23, 24, 3,
	5, 9, 8, 13, 11, 12, 16, 17, 28, 0,
	4, 10, 14, 15, 16, 17, 0, 0, 25, 14,
	15, 16, 17,
}
var yyPact = []int{

	-1000, 15, -1000, -1000, -1000, -1000, 28, -1000, -1000, -1000,
	-4, -1000, -4, -1000, -4, -4, -4, -4, 21, -10,
	28, 13, 13, -1000, -1000, -1000, -1000, -4, 28,
}
var yyPgo = []int{

	0, 14, 10, 0, 9, 1,
}
var yyR1 = []int{

	0, 1, 1, 2, 2, 2, 2, 3, 3, 3,
	3, 3, 4, 4, 4, 4, 4, 4, 5, 5,
}
var yyR2 = []int{

	0, 0, 2, 1, 1, 1, 1, 1, 3, 3,
	3, 3, 1, 1, 3, 1, 3, 1, 1, 3,
}
var yyChk = []int{

	-1000, -1, -2, 4, 15, 5, -3, -4, 7, 6,
	16, 9, 10, 8, 11, 12, 13, 14, -3, -5,
	-3, -3, -3, -3, -3, 17, 17, 18, -3,
}
var yyDef = []int{

	1, -2, 2, 3, 4, 5, 6, 7, 12, 13,
	0, 15, 0, 17, 0, 0, 0, 0, 0, 0,
	18, 8, 9, 10, 11, 14, 16, 0, 19,
}
var yyTok1 = []int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	16, 17, 13, 11, 18, 12, 3, 14, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 15,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10,
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

	case 3:
		//line tacky.y:24
		{
			cast(yylex).BeginForm(yyS[yypt-0].s)
		}
	case 4:
		//line tacky.y:25
		{
			cast(yylex).EndForm()
		}
	case 5:
		//line tacky.y:26
		{
			cast(yylex).Assign(yyS[yypt-0].s)
		}
	case 6:
		//line tacky.y:27
		{
			cast(yylex).Expr(yyS[yypt-0].expr)
		}
	case 8:
		//line tacky.y:30
		{
			yyVAL.expr = NewOp(yyS[yypt-1].s, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 9:
		//line tacky.y:31
		{
			yyVAL.expr = NewOp(yyS[yypt-1].s, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 10:
		//line tacky.y:32
		{
			yyVAL.expr = NewOp(yyS[yypt-1].s, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 11:
		//line tacky.y:33
		{
			yyVAL.expr = NewOp(yyS[yypt-1].s, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 12:
		//line tacky.y:35
		{
			yyVAL.expr = NewExpr("$", yyS[yypt-0].s)
		}
	case 13:
		//line tacky.y:36
		{
			yyVAL.expr = NewExpr("ID", yyS[yypt-0].s)
		}
	case 14:
		//line tacky.y:37
		{
			yyVAL.expr = yyS[yypt-1].expr
		}
	case 15:
		//line tacky.y:38
		{
			yyVAL.expr = NewExpr("XREF", yyS[yypt-0].s)
		}
	case 16:
		//line tacky.y:39
		{
			yyVAL.expr = NewFun(yyS[yypt-2].s, yyS[yypt-1].expr)
		}
	case 17:
		//line tacky.y:40
		{
			yyVAL.expr = NewExpr("%", yyS[yypt-0].s)
		}
	case 18:
		//line tacky.y:42
		{
			yyVAL.expr = NewExpr("", "")
			yyVAL.expr.AddKid(yyS[yypt-0].expr)
		}
	case 19:
		//line tacky.y:43
		{
			yyS[yypt-2].expr.AddKid(yyS[yypt-0].expr)
		}
	}
	goto yystack /* stack new state and value */
}
