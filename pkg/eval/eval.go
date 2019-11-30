package eval

import (
	"github.com/abdulbahajaj/brutus/pkg/types"
)


/*
* Function invocation
*/

func func_invoke(fn_call types.BrutList, env types.BrutEnv) (types.BrutType, types.BrutEnv){
	evaluatedList := types.NewBrutList()
	for _, val := range fn_call {
		val, newEnv := recEval(val, env)
		env = newEnv
		evaluatedList = evaluatedList.Append(val)
	}
	function := evaluatedList[0].(types.BrutPrimitive)
	evaluatedList = evaluatedList[1:]
	function.String()
	// return types.NewBrutNumber(10)
	return function(evaluatedList),env
}


/*
* Eval func
*/

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

//An eval function that is recursively called
func recEval(bType types.BrutType, env types.BrutEnv) (types.BrutType, types.BrutEnv){
	if bType.GetType() == types.SYMBOL{
		return env.LookUp(bType.(types.BrutSymbol)), env
	} else if isAtom(bType) {
		return bType, env
	} else {
		// return types.NewBrutNumber(10)
		return func_invoke(bType.(types.BrutList), env)
	}
}

//Sets up the environment and calls recEval
func Eval(bType types.BrutType, env types.BrutEnv) (types.BrutType, types.BrutEnv){
	return_stack := types.NewBrutList()
	if bType.GetType() == types.STACK {
		for _, exp := range bType.(types.BrutStack) {
			return_val, newEnv := recEval(exp, env)
			env = newEnv
			return_stack = return_stack.Append(return_val)
		}
		return return_stack, env
	}
	return recEval(bType, env)
}
