%{
package filter
import ("fmt";"io";"errors";"sync");
var g_grammer *Grammer;
var g_mutex  sync.Mutex
%}

%union {
	grammer *Grammer
	expr *Expr
	term *Term
	factor *Factor
	list *List
	fun *Func
	str string
	dval float64
	fn int
}

%token COMMA LPAREN RPAREN LAND LOR GET DEFAULT 
%type <grammer> grammer
%type <expr> expr
%type <term> term
%type <factor> factor
%type <fun> fun
%type <list> list
%token <str> VAR STR
%token <dval> NUM
%token <fn> CMP CONTAIN FUNC

%%
start: grammer { g_mutex.Lock(); g_grammer = $1; g_mutex.Unlock(); };
grammer: expr GET NUM grammer {var err error; if $$, err = NewGrammer(EGET, $1, $3, $4); err != nil {panic(err); }}
| DEFAULT GET NUM grammer {var err error; if $$, err = NewGrammer(DGET, nil, $3, $4); err != nil { panic(err); }} 
| expr grammer {var err error; if $$, err = NewGrammer(EEXPR, $1, 0, $2); err != nil { panic(err); }}
|              {$$ = nil; }

expr: expr LAND term {var err error; if $$, err = NewExpr(AND, $1, $3); err != nil { panic(err); }}
| expr LOR term {var err error; if $$, err = NewExpr(OR, $1, $3); err != nil { panic(err); }}
| term {var err error; if $$, err = NewExpr(TERM, nil, $1); err != nil { panic(err); }}

term: factor CONTAIN LPAREN list RPAREN {var err error; if $$, err = NewTerm(TKind_t($2), $1, $4, nil, nil);  err != nil {panic(err);} }
| factor CMP factor  {var err error; if $$, err = NewTerm(TKind_t($2), $1, nil, $3, nil); err != nil {panic(err); }}
|  LPAREN expr RPAREN {var err error; if $$, err = NewTerm(EXPR, nil, nil, nil, $2); err != nil { panic(err); }}

factor : VAR {var err error; if $$, err = NewFactor(VARIABLE, 0, "", $1, nil); err != nil { panic(err); }; }
| STR {var err error; if $$, err = NewFactor(STRING, 0, $1, "", nil); err != nil { panic(err); };}
| NUM {var err error; if $$, err = NewFactor(DOUBLE, $1, "", "", nil); err != nil { panic(err); };}
| fun {var err error; if $$, err = NewFactor(FUNCTION, 0, "", "", $1); err !=nil { panic(err); };}

list : factor {var err error; if $$, err = NewList($1, nil); err != nil { panic(err); };}
| factor COMMA list {var err error; if $$, err = NewList($1, $3); err != nil { panic(err);};}

fun: FUNC LPAREN list RPAREN  {var err error; if $$, err = NewFunc(FnKind_t($1), $3); err != nil { panic(err); }}
| FUNC LPAREN RPAREN {var err error; if $$, err = NewFunc(FnKind_t($1),nil); err != nil { panic(err); }} 

%%

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
		return h, errors.New("invalid rule");
	}	
    return h, err; 
}

/*
   get parse result
   symlist is created by calling QueryToSymlist() or JsonToSymlist() API
 */
func (h *Parser)Parse(symlist *SymList) (ret int, err error) {
	defer func() {
		if e := recover(); e != nil {
			ret, err = -1, errors.New(fmt.Sprint(e)) 
		}
	}()
	ret, err = EvalGrammer(h.grammer, symlist);
	return
}

