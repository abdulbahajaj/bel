package primitives

import (
	// "fmt"
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

func GetPrimitiveEnv() types.BrutEnv{
	env := types.NewBrutEnv()
	env = env.Set(types.BrutSymbol("+"), types.BrutPrimitive(Sum))
	return env
}
