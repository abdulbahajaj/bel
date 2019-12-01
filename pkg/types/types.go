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
	STACK
	ENV
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

type BrutPrimitive func(BrutList)BrutType

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

func (BrutEnv) GetType()ObjectType{
	return ENV
}

func (e BrutEnv)String()string{
	result := ""
	for key, val := range e.value{
		result += key + ": " + val.String() + "\n"
	}
	return result
}

func (e BrutEnv) LookUp(sym BrutSymbol)BrutType{
	if sym == "scope" {
		return e
	}
	if val, ok := e.value[sym.String()]; ok {
		return val
	} else {
		panic("Lookup failed: " + sym.String())
	}
}


/*
* Stacks
*/

type BrutStack []BrutType

func (st BrutStack) String() string {
	result := ""
	for _, obj := range st{
		result += obj.String() + "\n"
	}
	return result
}

func (BrutStack) GetType()ObjectType{
	return STACK
}


/*
* NIL
*/

type BrutNil bool

func (BrutNil) String() string {
	return "NIL"
}

func (BrutNil) GetType()ObjectType{
	return NIL
}
