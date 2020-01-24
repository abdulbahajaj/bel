package eval

import (
	"fmt"
	"github.com/abdulbahajaj/brutus/pkg/types"
	"github.com/abdulbahajaj/brutus/pkg/common"
)


/*
* Arg destructuring
*/

func destructure(params types.BrutList, args types.BrutList, env *types.BrutEnv) *types.BrutEnv {
	argsLen := len(args) - 1
	paramsLen := len(params) - 1
	if paramsLen == -1 {
		return env
	}

	last := params[paramsLen]
	if last.GetType() == types.SYMBOL {
		if lastSym := last.(types.BrutSymbol); lastSym[0] == '&' {
			var curArg types.BrutType
			if argsLen < paramsLen {
				curArg = types.NewBrutList()
			} else {
				curArg = args[paramsLen:]
			}
			env.SetParam(lastSym[1:], curArg)
			paramsLen = paramsLen - 1
		}
	}

	if argsLen < paramsLen {
		fmt.Println(params)
		fmt.Println(args)
		panic("Wrong Arity")
	}

	for cur := 0; cur <= paramsLen; cur++ {
		curParam := params[cur]
		curArg := args[cur]
		if curParam.GetType() == types.SYMBOL {
			env.SetParam(curParam.(types.BrutSymbol), curArg)
		} else if curParam.GetType() == types.LIST {
			env = destructure(curParam.(types.BrutList), curArg.(types.BrutList), env)
		} else {
			panic("Param can be either a symbol or a list of symbols")
		}
	}
	return env
}

/*
* Handle Callables
*/

func notCallablePanic(name string){
	panic("Not callable: " + name)
}

func callClo(lit types.BrutList, args types.BrutList, env *types.BrutEnv, evalArgs bool) (types.BrutType, *types.BrutEnv) {
	envDesc, newEnv := RecEval(lit[2], env)
	env = newEnv
	params := lit[3].(types.BrutList)
	body := lit[4]

	var cloEnv *types.BrutEnv
	switch envDesc.GetType(){
		case types.LIST:
		cloEnv = types.NewBrutEnv()
		envDescL := envDesc.(types.BrutList)
		for cursor := 0; cursor < len(envDescL); cursor += 2{
			cloEnv.SetParam(envDescL[cursor].(types.BrutSymbol), envDescL[cursor + 1] )
		}
		case types.ENV:
		cloEnv = envDesc.(*types.BrutEnv).AddScope()
	}

	cloEnv = destructure(params, args, cloEnv)
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

		case "tab":
		table := lit[2].(types.BrutTable)
		key := call[1]
		return table[key], env
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

func evalDef(exp types.BrutList, env *types.BrutEnv) (types.BrutType, *types.BrutEnv){
	for cursor := 1; cursor < len(exp); cursor += 2{
		val, newEnv := RecEval(exp[cursor + 1], env)
		env = newEnv

		env.Def(exp[cursor].(types.BrutSymbol), val)
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

	env.Def(name, mac)

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
		case "def":
			return evalDef(bList, env)
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


