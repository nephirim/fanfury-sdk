package mint

import (
	"github.com/persistenceOne/persistenceSDK/modules/splits/mapper"
	"github.com/persistenceOne/persistenceSDK/schema/utilities/base"
)

var Auxiliary = base.NewAuxiliary(
	mapper.ModuleName,
	AuxiliaryName,
	AuxiliaryRoute,
	initializeAuxiliaryKeeper,
)