package filter

import (
	"encoding/json"
	"errors"
	"fmt"
	js "github.com/bitly/go-simplejson"
	"strconv"
	"strings"
)

func NewSymlist(name, value string, kind FKind_t) (*SymList, error) {
	s := new(SymList)
	s.Kind = kind
	s.Name = name
	if kind == DOUBLE {
		if dbl, err := strconv.ParseFloat(value, 64); err != nil {
			return nil, err
		} else {
			s.Value = dbl
		}
	} else {
		s.Value = value
	}
	s.Next = nil
	return s, nil
}

func NewSymlistString(name, value string) (*SymList, error) {
	return NewSymlist(name, value, STRING)
}

func NewSymlistDouble(name string, value float64) (*SymList, error) {
	s := new(SymList)
	s.Kind = DOUBLE
	s.Name = name
	s.Value = value
	s.Next = nil
	return s, nil
}

func AppendSymlist(symlist *SymList, name, value string, kind FKind_t) (*SymList, error) {
	pre := symlist
	found := false
	for p := symlist; p != nil; p = p.Next {
		if name == p.Name {
			found = true
			break
		}
		pre = p
	}

	if !found {
		pre.Next, _ = NewSymlist(name, value, kind)
	}

	return symlist, nil
}

func AppendSymlistString(symlist *SymList, name, value string) (*SymList, error) {
	return AppendSymlist(symlist, name, value, STRING)
}

func AppendSymlistDouble(symlist *SymList, name string, value float64) (*SymList, error) {
	pre := symlist
	found := false
	for p := symlist; p != nil; p = p.Next {
		if name == p.Name {
			found = true
			break
		}
		pre = p
	}

	if !found {
		pre.Next, _ = NewSymlistDouble(name, value)
	}

	return symlist, nil
}

func DeleteSymlist(symlist *SymList) {
	pre := symlist
	for p := symlist; p != nil; {
		pre = p
		p = p.Next
		pre.Next = nil
	}
	symlist = nil
}

func SymbolLookup(symlist *SymList, name string) (*Factor, error) {
	for p := symlist; p != nil; p = p.Next {
		if name == p.Name {
			if p.Kind == DOUBLE {
				if v, err := cast2float64(p.Value); err != nil {
					return nil, err
				} else {
					return NewFactor(DOUBLE, v, "", "", nil)
				}
			} else {
				if v, err := cast2string(p.Value); err != nil {
					return nil, err
				} else {
					return NewFactor(STRING, 0, v, "", nil)
				}
			}
		}
	}
	return nil, errors.New(fmt.Sprintf("symbol '%s' not found", name))
}

func DumpSymlist(symlist *SymList) {
	dump := "NAME\tVALUE\tTYPE\n======================\n"
	for p := symlist; p != nil; p = p.Next {
		if p.Kind == DOUBLE {
			dump += fmt.Sprintf("%s\t%.2f\tDOUBLE\n", p.Name, p.Value)
		} else {
			dump += fmt.Sprintf("%s\t'%s'\tSTRING\n", p.Name, p.Value)
		}
	}
	fmt.Print(dump)
}

/*
 * parse a query string to symlist_t struct , string format should be:
 * nmq=testmq&mac=xxxx&bootid=xxxx...
 */
func QueryToSymlist(query string) (symlist *SymList, err error) {
	var middle int = 0
	var end int = 0
	data := query
	for {
		middle = strings.IndexByte(data, '=')
		if middle == -1 {
			break
		}

		end = strings.IndexByte(data[middle:], '&')
		if end == -1 {
			if symlist == nil {
				symlist, _ = NewSymlist(data[:middle], data[middle+1:], STRING)
			} else {
				symlist, _ = AppendSymlistString(symlist, data[:middle], data[middle+1:])
			}
			break
		} else {
			if symlist == nil {
				symlist, _ = NewSymlist(data[:middle], data[middle+1:middle+end], STRING)
			} else {
				symlist, _ = AppendSymlistString(symlist, data[:middle], data[middle+1:middle+end])
			}
			data = data[middle+end+1:]
		}
	}
	return symlist, err
}

/*
 * parse a JSON string to symlist_t struct, string format should be:
 * {"double_name":10.0, "interger_name": 99, "string_name":"FIFA WC 2014", ...}
 */
func JsonToSymlist(jstr string) (symlist *SymList, err error) {
	jsroot, err := js.NewJson([]byte(jstr))
	if err != nil {
		return nil, err
	}

	jsMap, err := jsroot.Map()
	if err != nil {
		return nil, err
	}

	for k, v := range jsMap {
		switch u := v.(type) {
		case json.Number:
			dbl, err := v.(json.Number).Float64()
			if err != nil {
				continue
			}
			if symlist == nil {
				symlist, _ = NewSymlistDouble(k, dbl)
			} else {
				symlist, _ = AppendSymlistDouble(symlist, k, dbl)
			}
		case string:
			if symlist == nil {
				symlist, _ = NewSymlistString(k, u)
			} else {
				symlist, _ = AppendSymlistString(symlist, k, u)
			}
		default:
			continue
		}

	}

	return symlist, nil
}
