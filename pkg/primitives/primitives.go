package primitives

import (
	// "fmt"
	"github.com/abdulbahajaj/brutus/pkg/types"
)

func Sum(l types.BrutList)types.BrutType{
	sum := 0.0
	for _, el := range l.Elements{
		if el.GetType() == types.NUMBER{
			sum += el.(types.BrutNumber).Value
		}
	}
	return types.NewBrutNumber(sum)
}

func GetPrimitiveEnv() types.BrutEnv{
	env := types.NewBrutEnv()
	env = env.Set(types.NewBrutSymbol("+"), types.NewBrutPrimitive(Sum))
	return env
}
