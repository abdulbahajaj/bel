package primitives

import (
	"fmt"
	"github.com/abdulbahajaj/brutus/pkg/types"
)

func Sum(l types.BrutList)types.BrutType{
	sum := 0.0
	for _, el := range l{
		if el.GetType() == types.NUMBER{
			sum += float64(el.(types.BrutNumber))
		}
	}
	return types.BrutNumber(sum)
}

func id(l types.BrutList)types.BrutType{
	first := l[0]
	for _, el := range l[1:]{
		if el.GetType() != first.GetType(){
			return types.BrutNil(false)
		}
		if el.GetType() == types.LIST {
			if &el != &first {
				return types.BrutNil(false)
			}
		} else if el != first {
			return types.BrutNil(false)
		}
	}
	return types.BrutSymbol("t")
}

func prn(l types.BrutList)types.BrutType{
	for _, el := range l{
		fmt.Print(el.String() + " ")
	}
	fmt.Print("\n")
	return l
}
func GetPrimitiveEnv() types.BrutEnv{
	env := types.NewBrutEnv()
	env = env.Set(types.BrutSymbol("+"), types.BrutPrimitive(Sum))
	env = env.Set(types.BrutSymbol("prn"), types.BrutPrimitive(prn))
	env = env.Set(types.BrutSymbol("id"), types.BrutPrimitive(id))
	return env
}
