package eval

import (
	"github.com/abdulbahajaj/brutus/pkg/types"
	"github.com/abdulbahajaj/brutus/pkg/primitives"
)

/*
* Function invocation
*/

func func_invoke(fn_call types.BrutList, env types.BrutEnv) types.BrutType{
	evaluatedList := types.NewBrutList()
	for _, val := range fn_call.Elements {
		evaluatedList = evaluatedList.Append(recEval(val, env))
	}
	function := evaluatedList.Elements[0].(types.BrutPrimitive)
	evaluatedList.Elements = evaluatedList.Elements[1:]
	function.String()
	// return types.NewBrutNumber(10)
	return function.FN(evaluatedList)
}


/*
* Eval func
*/

func isAtom(bType types.BrutType) bool{
	for _, objType := range []types.ObjectType{
		types.LIST,
		types.MODULE,
	} {
		if bType.GetType() == objType{
			return false
		}
	}
	return true
}

//An eval function that is recursively called
func recEval(bType types.BrutType, env types.BrutEnv) types.BrutType{
	if bType.GetType() == types.SYMBOL{
		return env.LookUp(bType.(types.BrutSymbol))
	} else if isAtom(bType) {
		return bType
	} else {
		// return types.NewBrutNumber(10)
		return func_invoke(bType.(types.BrutList), env)
	}
}

//Sets up the environment and calls recEval
func Eval(bType types.BrutType) types.BrutType{
	env := primitives.GetPrimitiveEnv()
	if bType.GetType() == types.MODULE {
		for _, exp := range bType.(types.BrutModule).Expressions {
			return recEval(exp, env)
		}
	}
	return recEval(bType, env)
}
