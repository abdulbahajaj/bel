package eval

import (
	"github.com/abdulbahajaj/brutus/pkg/types"
	"github.com/abdulbahajaj/brutus/pkg/common"
	// "fmt"
)


/*
* Evaluate callables
*/

func notCallablePanic(name string){
	panic("Not callable: " + name)
}

func callClo(lit types.BrutList, args types.BrutList, env types.BrutEnv) (types.BrutType, types.BrutEnv){
	paramNames := lit[3].(types.BrutList)
	body := lit[4]

	if len(paramNames) != len(args){
		panic("Wrong arity")
	}

	var cloEnv types.BrutEnv
	cloEnvBType := lit[2]
	switch cloEnvBType.GetType(){
		case types.NIL:
		cloEnv = types.NewBrutEnv()

		case types.SYMBOL:
		evaluated, newEnv := RecEval(cloEnvBType, env)
		env = newEnv
		cloEnv = evaluated.(types.BrutEnv)
	}

	cloEnv = cloEnv.SetParams(paramNames, args)

	returnVal, cloEnv := RecEval(body, cloEnv)

	cloEnv = cloEnv.ClearParams()
	//TODO override the variables from env with the ones from the cloEnv
	return returnVal, env
}

func invokeCallable(call types.BrutList, env types.BrutEnv) (types.BrutType, types.BrutEnv){
	if len(call) == 0 {
		notCallablePanic("()")
	}

	first := call[0]
	evFirst, newEnv := RecEval(first, env)
	env = newEnv
	if evFirst.GetType() != types.LIST{
		notCallablePanic(first.String())
	}

	lit := evFirst.(types.BrutList)
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
		return callClo(lit, args, env)

		case "prim":
		tempFunc := lit[3]
		if tempFunc.GetType() != types.PRIMITIVE {
			notCallablePanic(first.String())
		}
		function := tempFunc.(types.BrutPrimitive)
		args, newEnv := seqEval(call[1:], env)
		env = newEnv
		return function(args, env)
	}
	panic("error")

}


/*
* Special form evaluation
*/

func evalIf(bList types.BrutList, env types.BrutEnv) (types.BrutType, types.BrutEnv){
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

func evalSet(exp types.BrutList, env types.BrutEnv) (types.BrutType, types.BrutEnv){
	for cursor := 1; cursor < len(exp); cursor += 2{
		val, newEnv := RecEval(exp[cursor + 1], env)
		env = newEnv

		env = env.Set(exp[cursor].(types.BrutSymbol), val)
	}
	return types.BrutSymbol("t"), env
}


/*
* Evaluation
*/

//Sequentially evaluate elements in a list
func seqEval(list types.BrutList, env types.BrutEnv) (types.BrutList, types.BrutEnv){
	evaluatedList := types.NewBrutList()
	for _, val := range list{
		val, newEnv := RecEval(val, env)
		env = newEnv
		evaluatedList = evaluatedList.Append(val)
	}
	return evaluatedList, env
}

//An eval function that is recursively called
func RecEval(bType types.BrutType, env types.BrutEnv) (types.BrutType, types.BrutEnv){
	if bType.GetType() == types.SYMBOL{
		return env.LookUp(bType.(types.BrutSymbol)), env
	} else if common.IsAtom(bType) {
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
	return invokeCallable(bType.(types.BrutList), env)
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
