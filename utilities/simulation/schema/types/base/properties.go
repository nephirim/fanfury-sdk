package base

import (
	"math"
	"math/rand"

	"github.com/persistenceOne/persistenceSDK/schema/types"
	"github.com/persistenceOne/persistenceSDK/schema/types/base"
)

func GenerateRandomProperties(r *rand.Rand) types.Properties {
	randomPositiveInt := int(math.Abs(float64(r.Int()))) % 11

	propertyList := make([]types.Property, randomPositiveInt)

	for i := 0; i <= randomPositiveInt; i++ {
		propertyList[i] = GenerateRandomProperty(r)
	}

	return base.NewProperties(propertyList...)
}
