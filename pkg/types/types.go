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
	SYMBOL
	MODULE
	PRIMITIVE
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

// type envScope map[BrutSymbol]BrutType
// type BrutEnv struct{
// 	stack []envScope
// }

// func NewBrutEnv() BrutEnv{
// 	env := BrutEnv{stack: make([]envScope, 0)}
// 	env.AddScope()
// 	return env
// }

// func (e *BrutEnv) AddScope(){
// 	e.stack = append(e.stack, make(envScope))
// }

// func (e *BrutEnv) SetScope(sym BrutSymbol, val BrutType){
// 	last := len(e.stack) - 1
// 	e.stack[last][sym] = val
// }

// func (e *BrutEnv) PopScope(){
// 	e.stack = e.stack[:len(e.stack) - 1]
// }

// func (e *BrutEnv) Set(sym BrutSymbol, val BrutType){
// 	e.stack[0][sym] = val
// }

// func (e BrutEnv) LookUp(sym BrutSymbol) BrutType{
// 	for cursor := len(e.stack) - 1; cursor >= 0; cursor-- {
// 		stack := e.stack[cursor]
// 		if val, ok := stack[sym]; ok {
// 			return val
// 		}
// 	}
// 	panic("lookup failed" + sym.String())
// }

// func (BrutEnv) GetType()ObjectType{
// 	return ENV
// }

// func (e BrutEnv)String()string{
// 	result := ""
// 	ind := ""
// 	for _, scope := range e.stack{
// 		result += ind + "-> new scope\n"
// 		for sym, val := range scope{
// 			result += ind + string(sym) + " : " + val.String() + "\n"
// 		}
// 		ind += "\t"
// 	}
// 	return result
// }


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

func (bs *BrutEnv) Set(sym BrutSymbol, val BrutType){
	bs.first.vals[sym] = val
	if val.GetType() == ENV{
		fmt.Println("SETTING ENV")
	}
}

func (bs *BrutEnv) SetParam(sym BrutSymbol, val BrutType){
	bs.vals[sym] = val
	if val.GetType() == ENV{
		fmt.Println("SETTING ENV")
	}
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

func (bs *BrutEnv) Test(){
	fmt.Println(bs.parent == bs.parent.first)
	fmt.Println(bs.parent.vals)
	// fmt.Println(bs.parent.parent.parent.isRoot)
}

// func (bs *BrutEnv) String() string {
// 	result := ""
// 	fmt.Println(bs.parent.vals)
// 	// for sym, val := range bs.vals {
// 	// 	result += fmt.Sprintf("%v : %v", sym, val)
// 	// }
// 	return result
// }

func (bs *BrutEnv) isFirst(){
	if bs == bs.first {
		fmt.Println("I am first")
	}
	fmt.Println("I am not first")
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
