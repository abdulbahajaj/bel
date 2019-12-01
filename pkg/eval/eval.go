package eval

import (
	"github.com/abdulbahajaj/brutus/pkg/types"
	// "fmt"
)


/*
* Function invocation
*/

func invoke_func(fn_call types.BrutList, env types.BrutEnv) (types.BrutType, types.BrutEnv){
	evaluatedList := types.NewBrutList()
	for _, val := range fn_call {
		val, newEnv := RecEval(val, env)
		env = newEnv
		evaluatedList = evaluatedList.Append(val)
	}

	first := evaluatedList[0]

	if first.GetType() == types.PRIMITIVE{
		function := first.(types.BrutPrimitive)
		evaluatedList = evaluatedList[1:]
		function.String()
		return function(evaluatedList, env)
	} else {
		clo := first.(types.BrutList)
		if clo[1].(types.BrutSymbol) != "clo"{
			panic("Not callable")
		}

		// parsing closure
		cloEnvBType := clo[2]
		paramNames := clo[3].(types.BrutList)
		body := clo[4]
		paramValues := evaluatedList[1:]

		if len(paramNames) != len(paramValues){
			panic("Wrong arity")
		}

		if body.GetType() != types.LIST{
			return body, env
		} else {
			body = body.(types.BrutList)
		}

		// Getting closure local environment
		// isGlobalEnv := false
		var cloEnv types.BrutEnv

		switch cloEnvBType.GetType(){
			case types.NIL:
			cloEnv = types.NewBrutEnv()
			case types.SYMBOL:
			evaluated, newEnv := RecEval(cloEnvBType, env)
			env = newEnv
			cloEnv = evaluated.(types.BrutEnv)
		}

		// Setting closure parameters in env
		cloEnv = cloEnv.SetParams(paramNames, paramValues)

		//Running the expression
		returnVal, cloEnv := RecEval(body, cloEnv)
		cloEnv = cloEnv.ClearParams()

		return returnVal, env
	}
}



func isAtom(bType types.BrutType) bool{
	for _, objType := range []types.ObjectType{
		types.LIST,
		types.STACK,
	} {
		if bType.GetType() == objType{
			return false
		}
	}
	return true
}

func evalIf(bList types.BrutList, env types.BrutEnv)(types.BrutType, types.BrutEnv){
	for cursor := 1; cursor < len(bList); cursor += 2 {
		// fmt.Println(bList[cursor].GetType())
		if cursor == len(bList) - 1 {
			return RecEval(bList[cursor], env)
		}

		bType, newEnv := RecEval(bList[cursor], env)
		env = newEnv

		if bType.GetType() != types.NIL{
			return RecEval(bList[cursor + 1], env)
		}
	}
	return types.BrutNil(false), env
}

func evalSet(exp types.BrutList, env types.BrutEnv)(types.BrutType, types.BrutEnv){
	for cursor := 1; cursor < len(exp); cursor += 2{
		val, newEnv := RecEval(exp[cursor + 1], env)
		env = newEnv

		env = env.Set(exp[cursor].(types.BrutSymbol), val)
	}
	return types.BrutSymbol("t"), env
}

//An eval function that is recursively called
func RecEval(bType types.BrutType, env types.BrutEnv) (types.BrutType, types.BrutEnv){
	if bType.GetType() == types.SYMBOL{
		return env.LookUp(bType.(types.BrutSymbol)), env
	} else if isAtom(bType) {
		return bType, env
	}

	bList := bType.(types.BrutList)

	if len(bList) == 0{
		return types.BrutNil(false), env
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
		case "thread":
			// TODO there is an issue where threads don't have a reliable access to the global env
			go RecEval(bList, env)
			return types.BrutSymbol("t"), env
		}
	}

	//Default case - evaluate as functions
	return invoke_func(bType.(types.BrutList), env)
}

//Sets up the environment and calls RecEval
func Eval(bType types.BrutType, env types.BrutEnv) (types.BrutType, types.BrutEnv){
	// env.Set(types.BrutSymbol("eval"), )
	return_stack := types.NewBrutList()
	if bType.GetType() == types.STACK {
		for _, exp := range bType.(types.BrutStack) {
			return_val, newEnv := RecEval(exp, env)
			env = newEnv
			return_stack = return_stack.Append(return_val)
		}
		return return_stack, env
	}
	return RecEval(bType, env)
}
