/*
 Copyright [2019] - [2020], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceSDK contributors
 SPDX-License-Identifier: Apache-2.0
*/

package make

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceSDK/constants"
	assetsMapper "github.com/persistenceOne/persistenceSDK/modules/assets/mapper"
	"github.com/persistenceOne/persistenceSDK/modules/exchanges/auxiliaries/custody"
	"github.com/persistenceOne/persistenceSDK/modules/identities/auxiliaries/verify"
	"github.com/persistenceOne/persistenceSDK/modules/orders/mapper"
	"github.com/persistenceOne/persistenceSDK/schema/helpers"
	"github.com/persistenceOne/persistenceSDK/schema/types"
	"github.com/persistenceOne/persistenceSDK/schema/types/base"
)

type transactionKeeper struct {
	mapper                    helpers.Mapper
	identitiesVerifyAuxiliary helpers.Auxiliary
	exchangesCustodyAuxiliary helpers.Auxiliary
}

var _ helpers.TransactionKeeper = (*transactionKeeper)(nil)

func (transactionKeeper transactionKeeper) Transact(context sdkTypes.Context, msg sdkTypes.Msg) helpers.TransactionResponse {
	message := messageFromInterface(msg)

	if auxiliaryResponse := transactionKeeper.identitiesVerifyAuxiliary.GetKeeper().Help(context, verify.NewAuxiliaryRequest(message.From, message.MakerID)); !auxiliaryResponse.IsSuccessful() {
		return newTransactionResponse(auxiliaryResponse.GetError())
	}

	makerIsAsset := assetsMapper.NewAssets(assetsMapper.Mapper, context).Fetch(message.MakerSplitID).Get(message.MakerSplitID) != nil
	takerIsAsset := assetsMapper.NewAssets(assetsMapper.Mapper, context).Fetch(message.TakerSplitID).Get(message.TakerSplitID) != nil

	if makerIsAsset && takerIsAsset {
		if !sdkTypes.OneDec().Equal(message.ExchangeRate) {
			return newTransactionResponse(constants.IncorrectMessage)
		}
	} else if !makerIsAsset && takerIsAsset {
		if !message.MakerSplit.Mul(message.ExchangeRate).Equal(sdkTypes.OneDec()) {
			return newTransactionResponse(constants.IncorrectMessage)
		}
	}

	var immutablePropertyList []types.Property
	immutablePropertyList = append(immutablePropertyList,
		base.NewProperty(constants.MakerIDProperty, base.NewFact(message.MakerID.String())),
		base.NewProperty(constants.TakerIDProperty, base.NewFact(message.TakerID.String())),
		base.NewProperty(constants.MakerSplitIDProperty, base.NewFact(message.MakerSplitID.String())),
		base.NewProperty(constants.ExchangeRateProperty, base.NewFact(message.ExchangeRate.String())),
		base.NewProperty(constants.TakerSplitIDProperty, base.NewFact(message.TakerSplitID.String())))
	immutables := base.NewImmutables(base.NewProperties(immutablePropertyList))

	orderID := mapper.NewOrderID(base.NewID(context.ChainID()), message.MaintainersID, immutables.GetHashID())
	orders := mapper.NewOrders(transactionKeeper.mapper, context).Fetch(orderID)

	var makerSplit sdkTypes.Dec
	if orders.Get(orderID) != nil {
		oldMakerSplitFact := orders.Get(orderID).GetMutables().Get().Get(base.NewID(constants.MakerSplitProperty)).GetFact()
		oldMakerSplit, Error := sdkTypes.NewDecFromStr(oldMakerSplitFact.GetHash())
		if Error != nil {
			return newTransactionResponse(Error)
		}
		makerSplit = oldMakerSplit.Add(message.MakerSplit)
	} else {
		makerSplit = message.MakerSplit
	}

	var mutablePropertyList []types.Property
	mutablePropertyList = append(mutablePropertyList,
		base.NewProperty(base.NewID(constants.MakerSplitProperty), base.NewFact(makerSplit.String())))
	mutableProperties := base.NewProperties(mutablePropertyList)
	mutables := base.NewMutables(mutableProperties)

	order := mapper.NewOrder(orderID, mutables, immutables)
	if auxiliaryResponse := transactionKeeper.exchangesCustodyAuxiliary.GetKeeper().Help(context, custody.NewAuxiliaryRequest(message.MakerID, message.MakerSplit, message.MakerSplitID)); !auxiliaryResponse.IsSuccessful() {
		return newTransactionResponse(auxiliaryResponse.GetError())
	}

	if orders.Get(orderID) != nil {
		orders.Mutate(order)
	} else {
		orders.Add(order)
	}
	return newTransactionResponse(nil)
}

func initializeTransactionKeeper(mapper helpers.Mapper, externalKeepers []interface{}) helpers.TransactionKeeper {
	transactionKeeper := transactionKeeper{mapper: mapper}
	for _, externalKeeper := range externalKeepers {
		switch value := externalKeeper.(type) {
		case helpers.Auxiliary:
			switch value.GetName() {
			case verify.Auxiliary.GetName():
				transactionKeeper.identitiesVerifyAuxiliary = value
			case custody.Auxiliary.GetName():
				transactionKeeper.exchangesCustodyAuxiliary = value
			}
		}
	}
	return transactionKeeper
}
