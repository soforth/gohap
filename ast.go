/************************************************************************************
 author: soforth
 date: 2014-9-13
 email: soforth@qq.com
 description: a HTTP package filter, supported BNF as bellow:

 grammer -> expr => RET grammer     // if ( expr ) return RET else return grammer
		| default => RET grammer    // if grammer != 0 return grammer default: return RET
		| expr grammer              // if ( expr) return grammer else return 0
		;
 expr -> expr || term               // logic OR operation
		| expr && term              // logic AND operation
		| term
		;
 term -> factor  @ ( list )         // variable/string/double is in list
		| factor !@ ( list )        // variable/string/double is not in list
		| factor > factor           // left great than right
		| factor < factor           // left less than right
		| factor == factor          // left equal right
		| factor != factor          // left not equal right
		| factor >= factor          // left great than or equal right
		| factor <= factor          // left less than or equal right
		| factor # REGEX            // left regex match right pattern
		| factor !# REGEX           // left regex not match right pattern
		| ( expr )                  // term can be a expr in paren
		;
 list -> factor list;               // list is recursive defined
 factor -> DOUBLE_const             // factor can be double immediate constant
		| STRING_const              // also can be string immediate constant
		| VAR_str                   // also can be variable as symbol input by user
		| func                      // also can be a internal function
		;
 func -> VAR_str ( list );          // function has zero or more arguments
*************************************************************************************/
package filter

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"math"
	"regexp"
	"strconv"
)

type GKind_t int  // for Grammer
type EKind_t int  // for Expr
type TKind_t int  // for Term
type FKind_t int  // for Factor & SymList
type FnKind_t int // for Func

const (
	EGET  = GKind_t(0)
	DGET  = GKind_t(1)
	EEXPR = GKind_t(2)

	AND  = EKind_t(0)
	OR   = EKind_t(1)
	TERM = EKind_t(2)

	IN   = TKind_t(0)
	NI   = TKind_t(1)
	GT   = TKind_t(2)
	LT   = TKind_t(3)
	EQ   = TKind_t(4)
	NE   = TKind_t(5)
	GE   = TKind_t(6)
	LE   = TKind_t(7)
	MA   = TKind_t(8)
	NM   = TKind_t(9)
	EXPR = TKind_t(10)

	DOUBLE   = FKind_t(0)
	STRING   = FKind_t(1)
	VARIABLE = FKind_t(2)
	FUNCTION = FKind_t(3)

	LEN   = FnKind_t(0)
	MD5   = FnKind_t(1)
	COUNT = FnKind_t(2)
	ATOI  = FnKind_t(3)
	ITOA  = FnKind_t(4)
)

type Grammer struct {
	Kind    GKind_t // EGET, DGET, EEXPR
	Expr    *Expr
	Ret     float64
	Grammer *Grammer
}

type Expr struct {
	Kind  EKind_t // AND, OR, TERM
	Left  *Expr
	Right *Term
}

type Term struct {
	Kind  TKind_t // IN, NI, GT, LT, EQ, NE, GE, LE, MA, NM, EXPR
	Left  *Factor
	Right interface{}
}

type Factor struct {
	Kind  FKind_t // DOUBLE, STRING, VARIABLE, FUNCTION
	Value interface{}
}

type Func struct {
	Kind FnKind_t // LEN, MD5, COUNT, ATOI, ITOA
	List *List
}

type List struct {
	Factor *Factor
	Next   *List
}

type SymList struct {
	Kind  FKind_t // DOUBLE, STRING
	Name  string
	Value interface{}
	Next  *SymList
}

func gkind2str(kind GKind_t) string {
	switch kind {
	case EGET:
		return "<expr> =>"
	case DGET:
		return "default =>"
	case EEXPR:
		return "<expr>"
	}
	return fmt.Sprintf("%d", int(kind))
}

func ekind2str(kind EKind_t) string {
	switch kind {
	case AND:
		return "&&"
	case OR:
		return "||"
	case TERM:
		return "<term alone>"
	}
	return fmt.Sprintf("%d", int(kind))
}

func tkind2str(kind TKind_t) string {
	switch kind {
	case IN:
		return "@"
	case NI:
		return "!@"
	case GT:
		return ">"
	case LT:
		return "<"
	case GE:
		return ">="
	case LE:
		return "<="
	case EQ:
		return "=="
	case NE:
		return "!="
	case MA:
		return "#"
	case NM:
		return "!#"
	}
	return fmt.Sprintf("%d", int(kind))
}

func fkind2str(kind FKind_t) string {
	switch kind {
	case DOUBLE:
		return "float64"
	case STRING:
		return "string"
	case VARIABLE:
		return "var"
	case FUNCTION:
		return "func"
	}
	return fmt.Sprintf("%d", int(kind))
}

func fnkind2str(kind FnKind_t) string {
	switch kind {
	case LEN:
		return "len"
	case COUNT:
		return "count"
	case MD5:
		return "md5"
	case ITOA:
		return "itoa"
	case ATOI:
		return "atoi"
	}
	return fmt.Sprintf("%d", int(kind))
}

func bool2int(v bool) int {
	if v {
		return 1
	}
	return 0
}

func cast2string(i interface{}) (string, error) {
	if v, ok := i.(string); ok == true {
		return v, nil
	}
	return "", errors.New("not a 'string'")
}

func cast2float64(i interface{}) (float64, error) {
	if v, ok := i.(float64); ok == true {
		return v, nil
	}
	return 0, errors.New("not a 'float64'")
}

func cast2func(i interface{}) (*Func, error) {
	if v, ok := i.(*Func); ok == true {
		return v, nil
	}
	return nil, errors.New("not a '*Func'")
}

func NewGrammer(kind GKind_t, expr *Expr, ret float64, grammer *Grammer) (*Grammer, error) {
	g := new(Grammer)
	g.Kind = kind
	g.Expr = expr
	g.Ret = ret
	g.Grammer = grammer
	return g, nil
}

func NewExpr(kind EKind_t, expr *Expr, term *Term) (*Expr, error) {
	e := new(Expr)
	e.Kind = kind
	e.Left = expr
	e.Right = term
	return e, nil
}

func NewTerm(kind TKind_t, lfactor *Factor, list *List, rfactor *Factor, expr *Expr) (*Term, error) {
	t := new(Term)
	t.Kind = kind
	switch kind {
	case IN, NI:
		t.Left = lfactor
		t.Right = list
	case GT, LT, EQ, NE, GE, LE:
		t.Left = lfactor
		t.Right = rfactor
	case MA, NM:
		t.Left = lfactor
		if v, err := cast2string(rfactor.Value); err != nil {
			return t, err
		} else {
			if t.Right, err = regexp.CompilePOSIX(v); err != nil {
				return t, err
			}
		}
	case EXPR:
		t.Left = nil
		t.Right = expr
	}
	return t, nil
}

func NewFactor(kind FKind_t, dbl float64, str string, vari string, fn *Func) (*Factor, error) {
	f := new(Factor)
	f.Kind = kind
	switch kind {
	case DOUBLE:
		f.Value = dbl
	case STRING:
		f.Value = str
	case VARIABLE:
		f.Value = vari
	case FUNCTION:
		f.Value = fn
	}
	return f, nil
}

func NewFunc(kind FnKind_t, list *List) (*Func, error) {
	fn := new(Func)
	fn.Kind = kind
	fn.List = list
	return fn, nil
}

func NewList(factor *Factor, next *List) (*List, error) {
	l := new(List)
	l.Factor = factor
	l.Next = next
	return l, nil
}

func EvalGrammer(grammer *Grammer, symlist *SymList) (int, error) {
	var err error
	ret := -1
	if grammer == nil {
		return 0, nil
	}

	switch grammer.Kind {
	case EGET:
		ret, err = EvalExpr(grammer.Expr, symlist)
		if err != nil {
			return -1, err
		}
		if ret == 1 {
			return int(grammer.Ret), nil
		}
		return EvalGrammer(grammer.Grammer, symlist)
	case DGET:
		ret, err = EvalGrammer(grammer.Grammer, symlist)
		if err != nil {
			return -1, err
		}
		if ret == 0 {
			return int(grammer.Ret), nil
		}
		return ret, nil
	case EEXPR:
		ret, err = EvalExpr(grammer.Expr, symlist)
		if err != nil {
			return -1, err
		}
		if ret != 0 {
			return ret, nil
		}
		return EvalGrammer(grammer.Grammer, symlist)
	}

	return -1, errors.New(fmt.Sprintf("grammer operator '%s' not supported", gkind2str(grammer.Kind)))
}

func EvalExpr(expr *Expr, symlist *SymList) (int, error) {
	var err error
	var left int
	if expr == nil {
		return -1, errors.New("expr with invalid parameter")
	}

	switch expr.Kind {
	case AND:
		if left, err = EvalExpr(expr.Left, symlist); err != nil {
			return -1, err
		}
		if left <= 0 {
			return left, nil
		}
		return EvalTerm((*Term)(expr.Right), symlist)
	case OR:
		if left, err = EvalExpr(expr.Left, symlist); err != nil {
			return -1, err
		}
		if left != 0 {
			return left, nil
		}
		return EvalTerm((*Term)(expr.Right), symlist)
	case TERM:
		return EvalTerm((*Term)(expr.Right), symlist)
	}
	return -1, errors.New(fmt.Sprintf("expr operator '%s' not supported", ekind2str(expr.Kind)))
}

func EvalTerm(term *Term, symlist *SymList) (int, error) {
	if term == nil {
		return -1, errors.New("term with invalid parameter")
	}

	switch term.Kind {
	case IN, NI:
		switch v := term.Right.(type) {
		case *List:
			return EvalList(term.Kind, term.Left, v, symlist)
		}
	case GT, LT, EQ, NE, GE, LE:
		switch v := term.Right.(type) {
		case *Factor:
			return EvalCmp(term.Kind, term.Left, v, symlist)
		}
	case MA, NM:
		switch v := term.Right.(type) {
		case *regexp.Regexp:
			return EvalRegex(term.Kind, term.Left, v, symlist)
		}
	case EXPR:
		switch v := term.Right.(type) {
		case *Expr:
			return EvalExpr(v, symlist)
		}
	}
	return -1, errors.New(fmt.Sprintf("term with invalid kind '%s'", tkind2str(term.Kind)))
}

func EvalList(kind TKind_t, factor *Factor, list *List, symlist *SymList) (int, error) {
	found := false
	for p := list; p != nil; p = p.Next {
		if rc, err := EvalCmp(EQ, factor, p.Factor, symlist); err != nil {
			return -1, err
		} else if rc > 0 {
			found = true
			break
		}
	}
	if kind == IN {
		return bool2int(found), nil
	}
	return bool2int(!found), nil
}

func EvalRegex(kind TKind_t, lfactor *Factor, regex *regexp.Regexp, symlist *SymList) (int, error) {
	if lfactor == nil || regex == nil {
		return -1, errors.New("regexp with invalid parameter")
	}

	err := errors.New("regex match: parameter should be 'string'")
	lv := lfactor
	if lfactor.Kind == VARIABLE {
		if v, err := cast2string(lfactor.Value); err != nil {
			return -1, err
		} else if lv, err = SymbolLookup(symlist, v); err != nil {
			return -1, err
		}
	} else if lfactor.Kind == FUNCTION {
		if v, ok := lfactor.Value.(*Func); ok != true {
			return -1, err
		} else {
			var value *Factor
			value, err := EvalFunc(v, symlist)
			if err != nil {
				return -1, err
			}
			if value.Kind != STRING {
				return -1, errors.New("func ret should be 'string'")
			}
			lv = value
		}
	}

	if lv.Kind != STRING {
		return -1, err
	}

	if v, err := cast2string(lv.Value); err == nil {
		rc := regex.MatchString(v)
		return bool2int((kind == MA && rc == true) ||
			(kind == NM && rc == false)), nil
	}
	return -1, errors.New("left value has invalid type")
}

func EvalLen(list *List, symlist *SymList) (*Factor, error) {
	if list == nil || list.Factor == nil {
		return nil, errors.New("len() with invalid parameter")
	}

	switch list.Factor.Kind {
	case STRING:
		if v, err := cast2string(list.Factor.Value); err != nil {
			return nil, err
		} else {
			return NewFactor(DOUBLE, float64(len(v)), "", "", nil)
		}
	case VARIABLE:
		if v, err := cast2string(list.Factor.Value); err != nil {
			return nil, err
		} else {
			var value *Factor
			if value, err = SymbolLookup(symlist, v); err != nil {
				return nil, err
			}
			if value.Kind != STRING {
				return nil, err
			}
			if v2, err := cast2string(value.Value); err != nil {
				return nil, err
			} else {
				return NewFactor(DOUBLE, float64(len(v2)), "", "", nil)
			}
		}
	case FUNCTION:
		if v, ok := list.Factor.Value.(*Func); ok != true {
			return nil, errors.New("len(): parameter should be 'Func'")
		} else {
			var value *Factor
			value, err := EvalFunc(v, symlist)
			if err != nil {
				return nil, err
			}
			if value.Kind != STRING {
				return nil, errors.New("func ret should be 'string'")
			}
			if v2, err := cast2string(value.Value); err != nil {
				return nil, errors.New("func ret should be 'string'")
			} else {
				return NewFactor(DOUBLE, float64(len(v2)), "", "", nil)
			}
		}
	}

	return nil, errors.New(fmt.Sprintf("len with invalid kind '%s'", fkind2str(list.Factor.Kind)))
}

func EvalMD5(list *List, symlist *SymList) (*Factor, error) {
	if list == nil || list.Factor == nil {
		return nil, errors.New("md5() with invalid parameter")
	}

	deferr := errors.New("md5() parameter should be 'string'")
	h := md5.New()
	for p := list; p != nil; p = p.Next {
		switch p.Factor.Kind {
		case STRING:
			if v, err := cast2string(p.Factor.Value); err != nil {
				return nil, err
			} else {
				io.WriteString(h, v)
			}
		case VARIABLE:
			if v, err := cast2string(p.Factor.Value); err != nil {
				return nil, err
			} else {
				var value *Factor
				value, err := SymbolLookup(symlist, v)
				if err != nil {
					return nil, err
				}
				if value.Kind != STRING {
					return nil, deferr
				}
				if v2, err := cast2string(value.Value); err != nil {
					return nil, err
				} else {
					io.WriteString(h, v2)
				}
			}
		case FUNCTION:
			if v, err := cast2func(p.Factor.Value); err != nil {
				return nil, err
			} else {
				var value *Factor
				if value, err = EvalFunc(v, symlist); err != nil {
					return nil, err
				}
				if value.Kind != STRING {
					return nil, deferr
				}
				if v2, err := cast2string(value.Value); err != nil {
					return nil, err
				} else {
					io.WriteString(h, v2)
				}
			}
		}
	}
	return NewFactor(STRING, 0, fmt.Sprintf("%x", h.Sum(nil)), "", nil)
}

func EvalCount(symlist *SymList) (*Factor, error) {
	count := 0
	for p := symlist; p != nil; p = p.Next {
		count += 1
	}
	return NewFactor(DOUBLE, float64(count), "", "", nil)
}

func EvalAtoi(list *List, symlist *SymList) (*Factor, error) {
	if list == nil || list.Factor == nil {
		return nil, errors.New("atoi() with invalid parameter")
	}

	deferr := errors.New("atoi() parameter should be 'string'")
	switch list.Factor.Kind {
	case STRING:
		if v, err := cast2string(list.Factor.Value); err != nil {
			return nil, err
		} else {
			if dbl, err := strconv.ParseFloat(v, 64); err != nil {
				return nil, err
			} else {
				return NewFactor(DOUBLE, dbl, "", "", nil)
			}
		}
	case VARIABLE:
		if v, err := cast2string(list.Factor.Value); err != nil {
			return nil, err
		} else {
			var value *Factor
			if value, err = SymbolLookup(symlist, v); err != nil {
				return nil, err
			}
			if value.Kind != STRING {
				return nil, deferr
			}
			if v2, err := cast2string(value.Value); err != nil {
				return nil, err
			} else {
				if dbl, err := strconv.ParseFloat(v2, 64); err != nil {
					return nil, err
				} else {
					return NewFactor(DOUBLE, dbl, "", "", nil)
				}
			}
		}
	case FUNCTION:
		if v, err := cast2func(list.Factor.Value); err != nil {
			return nil, err
		} else {
			var value *Factor
			if value, err = EvalFunc(v, symlist); err != nil {
				return nil, err
			}
			if value.Kind != STRING {
				return nil, deferr
			}
			if v2, err := cast2string(value.Value); err != nil {
				return nil, err
			} else {
				if dbl, err := strconv.ParseFloat(v2, 64); err != nil {
					return nil, err
				} else {
					return NewFactor(DOUBLE, dbl, "", "", nil)
				}
			}
		}
	}

	return nil, errors.New(fmt.Sprintf("atoi with invalid kind '%s'", fkind2str(list.Factor.Kind)))
}

func EvalItoa(list *List, symlist *SymList) (*Factor, error) {
	if list == nil || symlist == nil {
		return nil, errors.New("itoa() with invalid parameter")
	}

	deferr := errors.New("itoa() parameter should be 'double'")
	switch list.Factor.Kind {
	case DOUBLE:
		if v, err := cast2float64(list.Factor.Value); err != nil {
			return nil, err
		} else {
			return NewFactor(STRING, 0, fmt.Sprintf("%f", v), "", nil)
		}
	case VARIABLE:
		if v, err := cast2string(list.Factor.Value); err != nil {
			return nil, err
		} else {
			var value *Factor
			if value, err = SymbolLookup(symlist, v); err != nil {
				return nil, err
			}
			if value.Kind != DOUBLE {
				return nil, deferr
			}
			if v2, err := cast2float64(value.Value); err != nil {
				return nil, err
			} else {
				return NewFactor(STRING, 0, fmt.Sprintf("%.2f", v2), "", nil)
			}
		}
	case FUNCTION:
		var value *Factor
		if v, err := cast2func(list.Factor.Value); err != nil {
			return nil, err
		} else {
			if value, err = EvalFunc(v, symlist); err != nil {
				return nil, err
			}
			if value.Kind != DOUBLE {
				return nil, deferr
			}
			if v2, err := cast2float64(value.Value); err != nil {
				return nil, err
			} else {
				return NewFactor(STRING, 0, fmt.Sprintf("%.2f", v2), "", nil)
			}
		}
	}

	return nil, errors.New(fmt.Sprintf("itoa with invalid kind '%s'", fkind2str(list.Factor.Kind)))
}

func EvalFunc(fn *Func, symlist *SymList) (*Factor, error) {
	if fn == nil || symlist == nil {
		return nil, errors.New(fmt.Sprintf("func '%s' with invalid parameter", fnkind2str(fn.Kind)))
	}
	switch fn.Kind {
	case LEN:
		return EvalLen(fn.List, symlist)
	case MD5:
		return EvalMD5(fn.List, symlist)
	case COUNT:
		return EvalCount(symlist)
	case ATOI:
		return EvalAtoi(fn.List, symlist)
	case ITOA:
		return EvalItoa(fn.List, symlist)
	}

	return nil, errors.New(fmt.Sprintf("function '%s' not supported", fnkind2str(fn.Kind)))
}

func EvalCmp(kind TKind_t, lfactor, rfactor *Factor, symlist *SymList) (int, error) {
	lv := lfactor
	rv := rfactor
	if lfactor.Kind == VARIABLE {
		if v, err := cast2string(lfactor.Value); err != nil {
			return -1, err
		} else {
			if lv, err = SymbolLookup(symlist, v); err != nil {
				return -1, err
			}
		}
	} else if lfactor.Kind == FUNCTION {
		if v, err := cast2func(lfactor.Value); err != nil {
			return -1, err
		} else {
			if lv, err = EvalFunc(v, symlist); err != nil {
				return -1, err
			}
		}
	}

	if rfactor.Kind == VARIABLE {
		if v, err := cast2string(rfactor.Value); err != nil {
			return -1, err
		} else {
			if rv, err = SymbolLookup(symlist, v); err != nil {
				return -1, err
			}
		}
	} else if rfactor.Kind == FUNCTION {
		if v, err := cast2func(rfactor.Value); err != nil {
			return -1, err
		} else {
			if rv, err = EvalFunc(v, symlist); err != nil {
				return -1, err
			}
		}
	}

	if lv.Kind != rv.Kind {
		return 0, nil // just ignore
	}

	if lv.Kind == DOUBLE {
		if v1, err := cast2float64(lv.Value); err != nil {
			return -1, err
		} else {
			if v2, err := cast2float64(rv.Value); err != nil {
				return -1, err
			} else {
				return CmpDbl(kind, v1, v2)
			}
		}
	} else if lv.Kind == STRING {
		if v1, err := cast2string(lv.Value); err != nil {
			return -1, err
		} else {
			if v2, err := cast2string(rv.Value); err != nil {
				return -1, err
			} else {
				return CmpStr(kind, v1, v2)
			}
		}
	}

	return -1, errors.New(fmt.Sprintf("operator '%s' not supported", tkind2str(kind)))
}

func CmpDbl(kind TKind_t, d1, d2 float64) (int, error) {
	switch kind {
	case GT:
		return bool2int(d1 > d2), nil
	case LT:
		return bool2int(d1 < d2), nil
	case EQ:
		return bool2int(math.Abs(d1-d2) < 0.001), nil
	case NE:
		return bool2int(math.Abs(d1-d2) > 0.001), nil
	case GE:
		return bool2int(d1 > d2 || math.Abs(d1-d2) < 0.001), nil
	case LE:
		return bool2int(d1 < d2 || math.Abs(d1-d2) < 0.001), nil
	}

	return -1, errors.New(fmt.Sprintf("double operator '%s' not supported", tkind2str(kind)))
}

func CmpStr(kind TKind_t, s1, s2 string) (int, error) {
	switch kind {
	case GT:
		return bool2int(s1 > s2), nil
	case LT:
		return bool2int(s1 < s2), nil
	case GE:
		return bool2int(s1 >= s2), nil
	case LE:
		return bool2int(s1 <= s2), nil
	case EQ:
		return bool2int(s1 == s2), nil
	case NE:
		return bool2int(s1 != s2), nil
	}

	return -1, errors.New(fmt.Sprintf("string operator '%s' not supported", tkind2str(kind)))
}
