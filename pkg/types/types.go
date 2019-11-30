package types

import (
	"fmt"
	// "strconv"
)

type ObjectType int
const (
	STRING ObjectType=iota
	LIST
	NUMBER
	CHARACHTER
	LITERAL
	NIL
	SYMBOL
	MODULE
	LAMBDA
	PRIMITIVE
)

type BrutAny interface {}

type BrutType interface {
	GetType() ObjectType
	String() string
}


/*
* Lists
*/

type BrutList struct {
	Elements []BrutType
}

func NewBrutList() BrutList {
	return BrutList{Elements: make([]BrutType,0 )}
}

func (bList BrutList) Append(el BrutType) BrutList{
	exist := ""
	for _, el := range bList.Elements {
		exist += el.String() + " "
	}
	bList.Elements = append(bList.Elements, el)

	exist = ""
	for _, el := range bList.Elements {
		exist += el.String() + " "
	}

	return bList
}

func (BrutList) GetType() ObjectType {
	return LIST
}

func (bList BrutList) String() string {
	result := ""

	if len(bList.Elements) == 0 {
		return "()"
	}

	for _, el := range bList.Elements {
		result += el.String()
		result += " "
	}


	result = "( " + result[:len(result)-1] + " )"
	return result
}


/*
* Tables
*/

type BrutTable struct {
	Value map[BrutType]BrutType
}

func NewBrutTable() BrutTable{
	return BrutTable{}
}


/*
* Module
*/

type BrutModule struct {
	Expressions []BrutList
	Env map[string]BrutList
}

func NewBrutModule() BrutModule{
	return BrutModule{Expressions: make([]BrutList,0)}
}

func (BrutModule) GetType() ObjectType{
	return MODULE
}

func (bModule BrutModule) AppendExp(bList BrutList) BrutModule{
	bModule.Expressions = append(bModule.Expressions, bList)
	return bModule
}

func (bModule BrutModule) String() string{
	return ""
}

/*
* Numbers
*/

type BrutNumber struct {
	Value float64
}

func NewBrutNumber(num float64) BrutNumber{
	return BrutNumber{Value:num}
}

func (BrutNumber) GetType() ObjectType{
	return NUMBER
}

func (bNumber BrutNumber) String() string{
	return fmt.Sprintf("%v", bNumber.Value)
}

/*
* Symbols
*/

type BrutSymbol struct {
	sym string
}

func NewBrutSymbol(sym string) BrutSymbol{
	return BrutSymbol{sym: sym}
}

func (BrutSymbol) GetType() ObjectType{
	return SYMBOL
}

func (bSym BrutSymbol) String() string{
	return bSym.sym
}


/*
* lambda
*/

// type BrutLambda struct{
// 	params []BrutSymbol
// 	body BrutList
// }

// func (BrutLambda) GetType() ObjectType{
// 	return LAMBDA
// }

// func (bLambda BrutLambda) String() string{
// 	result := ""
// 	for _,current_param :=range bLambda.params {
// 		result += current_param + " "
// 	}
// 	result += bLambda.String()

// 	return result
// }

/*
* BrutPrimitive
*/

type BrutPrimitive struct{
	FN func(BrutList)BrutType
}

func NewBrutPrimitive(fn func(BrutList)BrutType)BrutPrimitive{
	return BrutPrimitive{FN:fn}
}

func (BrutPrimitive) GetType() ObjectType{
	return PRIMITIVE
}

func (BrutPrimitive) String() string{
	return "Primitive"
}


/*
* Environment
*/

type BrutEnv struct{
	value map[string]BrutType
}

func NewBrutEnv() BrutEnv{
	return BrutEnv{value: make(map[string]BrutType)}
}

func (e BrutEnv) Set(sym BrutSymbol, val BrutType) BrutEnv{
	e.value[sym.String()] = val
	return e
}

func (e BrutEnv) LookUp(sym BrutSymbol)BrutType{
	if val, ok := e.value[sym.String()]; ok {
		return val
	} else {
		panic("Lookup failed: " + sym.String())
	}
}

// func (e Env) Set(sym types.BrutSymbol, val types.BrutType) Env{
// 	e.value[sym.String()] = val
// 	return e
// }

// func (e Env) setPrimitives() Env{
// 	return e
// }
