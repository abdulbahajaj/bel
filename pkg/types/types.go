package types

import (
	"fmt"
)

type ObjectType int
const (
	STRING ObjectType=iota
	LIST
	NUMBER
	CHARACHTER
	LITERAL
	SYMBOL
	MODULE
	PRIMITIVE
	ENV
	UNWRAP
	UNQUOTE
)

type BrutAny interface {}

type BrutType interface {
	GetType() ObjectType
	String() string
}


/*
* Lists
*/

type BrutList []BrutType


func NewBrutList() BrutList {
	return make(BrutList, 0)
}

func (bList BrutList) Append(el BrutType) BrutList{
	bList = append(bList, el)
	return bList
}

func (BrutList) GetType() ObjectType {
	return LIST
}

func (bList BrutList) String() string {
	result := ""

	if len(bList) == 0 {
		return "()"
	}

	for _, el := range bList {
		result += el.String()
		result += " "
	}


	result = "( " + result[:len(result)-1] + " )"
	return result
}


/*
* Tables
* TODO Implement
*/

type BrutTable map[BrutType]BrutType



/*
* Numbers
*/

type BrutNumber float64

func (BrutNumber) GetType() ObjectType{
	return NUMBER
}

func (bNumber BrutNumber) String() string{
	return fmt.Sprintf("%v", float64(bNumber))
}


/*
* Strings
*/

type BrutString string

func (BrutString) GetType() ObjectType{
	return STRING
}

func (s BrutString) String() string {
	return string(s)
}


/*
* Symbols
*/

type BrutSymbol string

func (BrutSymbol) GetType() ObjectType{
	return SYMBOL
}

func (bSym BrutSymbol) String() string{
	return string(bSym)
}


/*
* BrutPrimitive
*/

type BrutPrimitive func(BrutList, *BrutEnv)(BrutType, *BrutEnv)

func (BrutPrimitive) GetType() ObjectType{
	return PRIMITIVE
}

func (BrutPrimitive) String() string{
	return ""
}


/*
* Environment
*/

type scopeVals map[BrutSymbol]BrutType

type BrutEnv struct {
	isRoot bool
	first *BrutEnv
	parent *BrutEnv
	vals scopeVals
}

func NewBrutEnv() *BrutEnv{
	newScope := BrutEnv{vals: make(scopeVals), first: nil, isRoot: false }
	newScope.setParent(&newScope)
	return &newScope
}

func (bs *BrutEnv) MakeGlobal(){
	bs.first = bs
	bs.isRoot = true
}

func (bs *BrutEnv) setParent(parent *BrutEnv){
	bs.parent = parent
	bs.first = parent.first
}

func (bs *BrutEnv) AddScope() *BrutEnv{
	newScope := NewBrutEnv()
	newScope.setParent(bs)
	return newScope
}

func (bs *BrutEnv) Def(sym BrutSymbol, val BrutType){
	bs.first.vals[sym] = val
}

func (bs *BrutEnv) SetParam(sym BrutSymbol, val BrutType){
	bs.vals[sym] = val
}

func (bs *BrutEnv) LookUp(sym BrutSymbol) BrutType{
	if sym == "scope" {
		return bs
	}
	if val, ok := bs.vals[sym]; ok {
		return val
	}
	if bs.first == bs {
		panic("Lookup failed: " + sym.String())
	}
	return bs.parent.LookUp(sym)
}

func (BrutEnv) GetType() ObjectType{
	return ENV
}


func (bs BrutEnv) String() string{
	result := ""

	if !bs.isRoot {//&bs != bs.first {
		result += bs.parent.String()
	}

	result += "{\n"

	for sym, val := range bs.vals{
		switch val.GetType(){
			case LIST:
			result += fmt.Sprintf("\t %v : %p \n", sym, &val)
			case ENV:
			result += fmt.Sprintf("\t %v : %p \n", sym, &val)
			default:
			result += fmt.Sprintf("\t %v : %v \n", sym, val)
		}
	}
	result += "}\n"

	return result
}


/*
* Internal types
*/

// Unwrap
type BrutUnwrap struct {
	Val BrutType
}

func NewUnwrap(val BrutType) BrutUnwrap{
	return BrutUnwrap{Val: val}
}

func (b BrutUnwrap) String() string{
	return b.Val.String()
}

func (b BrutUnwrap) GetType() ObjectType{
	return UNWRAP
}

// Unquote
type BrutUnquote struct {
	Val BrutType
}

func NewUnquote(val BrutType) BrutUnquote{
	return BrutUnquote{Val: val}
}

func (b BrutUnquote) String() string{
	return b.Val.String()
}

func (b BrutUnquote) GetType() ObjectType{
	return UNQUOTE
}
