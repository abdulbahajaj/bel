package eval

import (
	"fmt"
	"github.com/abdulbahajaj/brutus/pkg/types"
	"github.com/abdulbahajaj/brutus/pkg/common"
)

func notCallablePanic(name string){
	panic("Not callable: " + name)
}

func unwrapParams(params types.BrutList, args types.BrutList, env *types.BrutEnv) *types.BrutEnv{
	paramsLen := len(params)
	argsLen := len(args)

	for cursor := 0; cursor < paramsLen; cursor++ {
		currentParam := params[cursor]
		var currentArg types.BrutType
		if cursor == paramsLen - 1 && argsLen > paramsLen {
			currentArg = args[cursor:]
		} else {
			currentArg = args[cursor]
		}
		if common.IsAtom(currentParam) {
			env.SetParam(currentParam.(types.BrutSymbol), currentArg)
		} else if !common.IsAtom(currentParam) && !common.IsAtom(currentArg) {
			env = unwrapParams(currentParam.(types.BrutList), currentArg.(types.BrutList), env)
		} else {
			panic("args don't match params")
		}
	}
	return env
}

func callClo(lit types.BrutList, args types.BrutList, env *types.BrutEnv, evalArgs bool) (types.BrutType, *types.BrutEnv) {
	envDesc, newEnv := RecEval(lit[2], env)
	env = newEnv
	params := lit[3].(types.BrutList)
	body := lit[4]

	// if len(params) != len(args){
	// 	panic("Wrong arity")
	// }

	var cloEnv *types.BrutEnv
	switch envDesc.GetType(){
		case types.LIST:
		cloEnv = types.NewBrutEnv()
		envDescL := envDesc.(types.BrutList)
		for cursor := 0; cursor < len(envDescL); cursor += 2{
			cloEnv.SetParam(envDescL[cursor].(types.BrutSymbol), envDescL[cursor + 1] )
		}
		case types.ENV:
		// envTemp :=
		cloEnv = envDesc.(*types.BrutEnv).AddScope()

	}

	cloEnv = unwrapParams(params, args, cloEnv)
	returnVal, cloEnv := RecEval(body, cloEnv)

	return returnVal, env
}

func invokeCallable(call types.BrutList, env *types.BrutEnv) (types.BrutType, *types.BrutEnv){
	if len(call) == 0 {
		notCallablePanic("()")
	}

	first := call[0]
	envFirst, newEnv := RecEval(first, env)
	env = newEnv
	if envFirst.GetType() != types.LIST{
		notCallablePanic(first.String())
	}

	lit := envFirst.(types.BrutList)
	if len(lit) < 3 {
		fmt.Println("LIT IS " + lit.String())
		notCallablePanic(first.String())
	}

	litFirst := lit[0]
	litType := lit[1]

	if litFirst.GetType() != types.SYMBOL || litType.GetType() != types.SYMBOL {
		notCallablePanic(first.String())
	}

	switch litType.(types.BrutSymbol){
		case "clo":
		args, newEnv := seqEval(call[1:], env)
		env = newEnv
		return callClo(lit, args, env, true)

		case "prim":
		tempFunc := lit[3]
		if tempFunc.GetType() != types.PRIMITIVE {
			notCallablePanic(first.String())
		}
		function := tempFunc.(types.BrutPrimitive)
		args, newEnv := seqEval(call[1:], env)
		env = newEnv
		return function(args, env)

		case "mac":
		clo := lit[2]
		if clo.GetType() != types.LIST{
			notCallablePanic(first.String())
		}
		args := call[1:]
		expansion, env := callClo(clo.(types.BrutList), args, env, false)
		return RecEval(expansion, env)

	}
	panic("Unrecognized callable")
}


/*
* Special form evaluation
*/

func evalIf(bList types.BrutList, env *types.BrutEnv) (types.BrutType, *types.BrutEnv){
	for cursor := 1; cursor < len(bList); cursor += 2 {
		if cursor == len(bList) - 1 {
			return RecEval(bList[cursor], env)
		}
	
		test := bList[cursor]
		conseq := bList[cursor + 1]

		testRes, newEnv := RecEval(test, env)
		env = newEnv

		exec := true
		if testRes.GetType() == types.LIST{
			if len(testRes.(types.BrutList)) == 0 {
				exec = false
			}
		}

		if exec {
			return RecEval(conseq, env)
		}
	}
	return types.BrutType(types.NewBrutList()), env
}

func evalSet(exp types.BrutList, env *types.BrutEnv) (types.BrutType, *types.BrutEnv){
	for cursor := 1; cursor < len(exp); cursor += 2{
		val, newEnv := RecEval(exp[cursor + 1], env)
		env = newEnv

		env.Set(exp[cursor].(types.BrutSymbol), val)
	}
	return types.BrutSymbol("t"), env
}


/*
* Evaluation
*/

//Sequentially evaluate elements in a list
func seqEval(list types.BrutList, env *types.BrutEnv) (types.BrutList, *types.BrutEnv){
	evaluatedList := types.NewBrutList()
	for _, val := range list{
		val, newEnv := RecEval(val, env)
		env = newEnv
		evaluatedList = evaluatedList.Append(val)
	}
	return evaluatedList, env
}

func setMac(def types.BrutList, env *types.BrutEnv) (*types.BrutEnv){
	name := def[1].(types.BrutSymbol)
	args := def[2]
	body := def[3]

	clo := types.NewBrutList()
	clo = clo.Append(types.BrutSymbol("lit"))
	clo = clo.Append(types.BrutSymbol("clo"))
	clo = clo.Append(types.BrutSymbol("scope"))
	clo = clo.Append(args)
	clo = clo.Append(body)

	mac := types.NewBrutList()
	mac = mac.Append(types.BrutSymbol("lit"))
	mac = mac.Append(types.BrutSymbol("mac"))
	mac = mac.Append(clo)

	env.Set(name, mac)

// (mac nilwith (x)
//   (list 'cons nil x))

	return env
}

//An eval function that is recursively called
func RecEval(bType types.BrutType, env *types.BrutEnv) (types.BrutType, *types.BrutEnv){
	if bType.GetType() == types.SYMBOL{
		return env.LookUp(bType.(types.BrutSymbol)), env
	} else if common.IsAtom(bType) {
		return bType, env
	}

	bList := bType.(types.BrutList)

	if len(bList) == 0{
		return types.NewBrutList(), env
	}

	// Evaluate special forms
	first := bList[0]
	if first.GetType() == types.SYMBOL{
		switch first := bList[0].(types.BrutSymbol); first {
		case "if":
			return evalIf(bList, env)
		case "lit":
			return bList, env
		case "quote":
			return bList[1], env
		case "set":
			return evalSet(bList, env)
		case "do":
			res, newEnv := seqEval(bList[1:], env)
			env = newEnv
			return res[len(res) - 1], env
		case "thread":
			// TODO there is an issue where threads don't have a reliable access to the global env
			go RecEval(bList, env)
			return types.BrutSymbol("t"), env
		}
	}

	//Default case - evaluate as functions
	return invokeCallable(bType.(types.BrutList), env)
}

//Sets up the environment and calls RecEval
func Eval(bType types.BrutType, env *types.BrutEnv) (types.BrutType, *types.BrutEnv){
	return_stack := types.NewBrutList()

	for _, exp := range bType.(types.BrutList) {
		return_val, newEnv := RecEval(exp, env)
		env = newEnv
		return_stack = return_stack.Append(return_val)
	}

	return return_stack, env
}


