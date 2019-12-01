package primitives

import (
	"fmt"
	"github.com/abdulbahajaj/brutus/pkg/types"
	"github.com/abdulbahajaj/brutus/pkg/eval"
)

func Sum(l types.BrutList, env types.BrutEnv)(types.BrutType, types.BrutEnv){
	var sum float64 = 0.0
	for _, el := range l{
		if el.GetType() == types.NUMBER{
			sum += float64(el.(types.BrutNumber))
		}
	}
	return types.BrutNumber(sum), env
}

func id(l types.BrutList, env types.BrutEnv)(types.BrutType, types.BrutEnv){
	first := l[0]
	for _, el := range l[1:]{
		if el.GetType() != first.GetType(){
			return types.BrutNil(false), env
		}
		if el.GetType() == types.LIST {
			if &el != &first {
				return types.BrutNil(false), env
			}
		} else if el != first {
			return types.BrutNil(false), env
		}
	}
	return types.BrutSymbol("t"), env
}

func prn(l types.BrutList, env types.BrutEnv)(types.BrutType, types.BrutEnv){
	for _, el := range l{
		fmt.Print(el.String() + " ")
	}
	fmt.Print("\n")
	return l, env
}

func evaluate(exp types.BrutList, env types.BrutEnv)(types.BrutType, types.BrutEnv){
	return eval.RecEval(exp, env)
}

func GetPrimitiveEnv() types.BrutEnv{
	env := types.NewBrutEnv()
	env = env.Set(types.BrutSymbol("+"), types.BrutPrimitive(Sum))
	env = env.Set(types.BrutSymbol("prn"), types.BrutPrimitive(prn))
	env = env.Set(types.BrutSymbol("id"), types.BrutPrimitive(id))
	env = env.Set(types.BrutSymbol("eval"), types.BrutPrimitive(evaluate))
	return env
}
