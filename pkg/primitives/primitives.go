package primitives

import (
	"fmt"
	"github.com/abdulbahajaj/brutus/pkg/types"
	"github.com/abdulbahajaj/brutus/pkg/eval"
	"github.com/abdulbahajaj/brutus/pkg/common"
)

func sum(l types.BrutList, env types.BrutEnv)(types.BrutType, types.BrutEnv){
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

func list(exp types.BrutList, env types.BrutEnv)(types.BrutType, types.BrutEnv){
	return exp, env
}

func append(exp types.BrutList, env types.BrutEnv)(types.BrutType, types.BrutEnv){
	result := types.NewBrutList()
	for _, el := range exp {
		if common.IsAtom(el){
			result = result.Append(el)
		}else {
			for _, el2 := range el.(types.BrutList){
				result = result.Append(el2)
			}
		}
	}
	return result, env
}

func cons(exp types.BrutList, env types.BrutEnv)(types.BrutType, types.BrutEnv){
	lastEl := exp[len(exp) -1].(types.BrutList)
	result := make(types.BrutList, 0, len(lastEl) + len(exp) - 1)
	for _, el := range exp[:len(exp)-1]{
		result = result.Append(el)
	}
	for _, el := range lastEl {
		result = result.Append(el)
	}

	return result, env
}

func setPrimitive(name string, env types.BrutEnv, fn func(exp types.BrutList, env types.BrutEnv)(types.BrutType, types.BrutEnv))types.BrutEnv{
	lit := types.NewBrutList()
	lit = lit.Append(types.BrutSymbol("lit"))
	lit = lit.Append(types.BrutSymbol("prim"))
	lit = lit.Append(types.BrutSymbol(name))
	lit = lit.Append(types.BrutPrimitive(fn))

	env = env.Set(types.BrutSymbol(name), lit)
	return env
}

func GetPrimitiveEnv() types.BrutEnv{
	env := types.NewBrutEnv()

	env = setPrimitive("+", env, sum)
	env = setPrimitive("prn", env, prn)
	env = setPrimitive("id", env, id)
	env = setPrimitive("eval", env, evaluate)
	env = setPrimitive("list", env, list)
	env = setPrimitive("cons", env, cons)
	env = setPrimitive("append", env, append)

	// env = env.Set(types.BrutSymbol("+"), types.BrutPrimitive(sum))
	// env = env.Set(types.BrutSymbol("prn"), types.BrutPrimitive(prn))
	// env = env.Set(types.BrutSymbol("id"), types.BrutPrimitive(id))
	// env = env.Set(types.BrutSymbol("eval"), types.BrutPrimitive(evaluate))
	// env = env.Set(types.BrutSymbol("list"), types.BrutPrimitive(list))
	// env = env.Set(types.BrutSymbol("cons"), types.BrutPrimitive(cons))
	// env = env.Set(types.BrutSymbol("append"), types.BrutPrimitive(append))
	return env
}
