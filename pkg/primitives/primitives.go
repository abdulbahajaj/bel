package primitives

import (
	"fmt"
	"github.com/abdulbahajaj/brutus/pkg/types"
	"github.com/abdulbahajaj/brutus/pkg/eval"
	"github.com/abdulbahajaj/brutus/pkg/common"
)

func sum(l types.BrutList, env *types.BrutEnv)(types.BrutType, *types.BrutEnv){
	var sum float64 = 0.0
	for _, el := range l{
		if el.GetType() == types.NUMBER{
			sum += float64(el.(types.BrutNumber))
		}
	}
	return types.BrutNumber(sum), env
}

func id(l types.BrutList, env *types.BrutEnv)(types.BrutType, *types.BrutEnv){
	first := l[0]
	for _, el := range l[1:]{
		if el.GetType() != first.GetType(){
			return types.NewBrutList(), env
		}
		if el.GetType() == types.LIST {
			if len(el.(types.BrutList)) == 0 && len(el.(types.BrutList)) == 0 {
			} else if &el != &first {
				return types.NewBrutList(), env
			}
		} else if el != first {
			return types.NewBrutList(), env
		}
	}
	return types.BrutSymbol("t"), env
}

func prn(l types.BrutList, env *types.BrutEnv)(types.BrutType, *types.BrutEnv){
	for _, el := range l{
		fmt.Print(el.String() + " ")
	}
	fmt.Print("\n")
	return types.BrutNumber(1), env
}

func evaluate(exp types.BrutList, env *types.BrutEnv)(types.BrutType, *types.BrutEnv){
	return eval.RecEval(exp, env)
}

func list(exp types.BrutList, env *types.BrutEnv)(types.BrutType, *types.BrutEnv){
	return exp, env
}

func append(exp types.BrutList, env *types.BrutEnv)(types.BrutType, *types.BrutEnv){
	result := types.NewBrutList()
	for _, el := range exp {
		if common.IsAtom(el){
			result = result.Append(el)
		} else {
			for _, el2 := range el.(types.BrutList){
				result = result.Append(el2)
			}
		}
	}
	return result, env
}

func cons(exp types.BrutList, env *types.BrutEnv)(types.BrutType, *types.BrutEnv){
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

func bmap(l types.BrutList, env *types.BrutEnv)(types.BrutType, *types.BrutEnv){
	f := l[0]
	toIterateOver := l[1].(types.BrutList)
	results := types.NewBrutList()
	for _, el := range toIterateOver{
		call := types.NewBrutList()
		call = call.Append(f)
		call = call.Append(el)
		callRes, newEnv := eval.RecEval(call, env)
		env = newEnv
		results = results.Append(callRes)
	}
	return results, env
}

func length(l types.BrutList, env *types.BrutEnv)(types.BrutType, *types.BrutEnv){
	el := l[0]
	if el.GetType() != types.LIST {
		return types.BrutNumber(1), env
	}
	return types.BrutNumber(len(el.(types.BrutList))), env
}

func biggerThan(l types.BrutList, env *types.BrutEnv)(types.BrutType, *types.BrutEnv){
	if l[0].(types.BrutNumber) > l[1].(types.BrutNumber) {
		return types.BrutSymbol("t"), env
	}
	return types.NewBrutList(), env
}

func smallerThan(l types.BrutList, env *types.BrutEnv)(types.BrutType, *types.BrutEnv){
	if l[0].(types.BrutNumber) < l[1].(types.BrutNumber) {
		return types.BrutSymbol("t"), env
	}
	return types.NewBrutList(), env
}

func first(l types.BrutList, env *types.BrutEnv)(types.BrutType, *types.BrutEnv){
	li := l[0].(types.BrutList)
	return li[0], env
}

func second(l types.BrutList, env *types.BrutEnv)(types.BrutType, *types.BrutEnv){
	li := l[0].(types.BrutList)
	return li[0], env
}

func nth(l types.BrutList, env *types.BrutEnv)(types.BrutType, *types.BrutEnv){
	li := l[0].(types.BrutList)
	index := l[1].(types.BrutNumber)
	return  li[int(index)], env
}

func mod(l types.BrutList, env *types.BrutEnv)(types.BrutType, *types.BrutEnv){
	return types.BrutNumber(int(l[0].(types.BrutNumber)) % int(l[1].(types.BrutNumber))), env
}

func filter(l types.BrutList, env *types.BrutEnv)(types.BrutType, *types.BrutEnv){
	f := l[0]
	toIterateOver := l[1].(types.BrutList)
	results := types.NewBrutList()
	for _, el := range toIterateOver{
		call := types.NewBrutList()
		call = call.Append(f)
		call = call.Append(el)
		callRes, newEnv := eval.RecEval(call, env)
		env = newEnv
		if callRes.GetType() == types.LIST {
			if len(callRes.(types.BrutList)) == 0 {
				continue
			}
		}
		results = results.Append(callRes)
	}
	return results, env
}

func bRange(l types.BrutList, env *types.BrutEnv)(types.BrutType, *types.BrutEnv){
	results := types.NewBrutList()
	start := int(l[0].(types.BrutNumber))
	end := int(l[1].(types.BrutNumber))

	for cur := start; cur < end; cur++{
		results = results.Append(types.BrutNumber(cur))
	}

	return results, env

}

func index(l types.BrutList, env *types.BrutEnv)(types.BrutType, *types.BrutEnv){
	results := types.NewBrutList()
	for key, val := range results{
		pair := types.NewBrutList()
		pair = pair.Append(types.BrutNumber(key))
		pair = pair.Append(val)
		results.Append(pair)
	}
	return results, env
}

func table(l types.BrutList, env *types.BrutEnv)(types.BrutType, *types.BrutEnv){
	tb := types.NewBrutTable()
	if len(l) % 2 != 0{
		panic("Key, value pairs in tables must be even")
	}
	for cur := 0; cur < len(l); cur += 2{
		tb[l[cur]] = l[cur + 1]
	}

	lit := types.NewBrutList()
	lit = lit.Append(types.BrutSymbol("lit"))
	lit = lit.Append(types.BrutSymbol("tab"))
	lit = lit.Append(tb)

	return lit, env
}

func setPrimitive(name string, env *types.BrutEnv, fn func(exp types.BrutList, env *types.BrutEnv)(types.BrutType, *types.BrutEnv)){
	lit := types.NewBrutList()
	lit = lit.Append(types.BrutSymbol("lit"))
	lit = lit.Append(types.BrutSymbol("prim"))
	lit = lit.Append(types.BrutSymbol(name))
	lit = lit.Append(types.BrutPrimitive(fn))

	env.Def(types.BrutSymbol(name), lit)
}

func GetPrimitiveEnv() *types.BrutEnv{
	env := types.NewBrutEnv()
	env.MakeGlobal()


	setPrimitive("nth", env, nth)
	setPrimitive("table", env, table)
	setPrimitive("range", env, bRange)
	setPrimitive("filter", env, filter)
	setPrimitive(">", env, biggerThan)
	setPrimitive("<", env, smallerThan)
	setPrimitive("mod", env, mod)
	setPrimitive("len", env, length)
	setPrimitive("+", env, sum)
	setPrimitive("prn", env, prn)
	setPrimitive("id", env, id)
	setPrimitive("eval", env, evaluate)
	setPrimitive("list", env, list)
	setPrimitive("cons", env, cons)
	setPrimitive("append", env, append)
	setPrimitive("map", env, bmap)
	return env
}
