package common

import (
	"github.com/abdulbahajaj/brutus/pkg/types"
)

func IsAtom(bType types.BrutType) bool{
	for _, objType := range []types.ObjectType{
		types.LIST,
		// types.STACK,
	} {
		if bType.GetType() == objType{
			return false
		}
	}
	return true
}


