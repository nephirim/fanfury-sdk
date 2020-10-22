/*
 Copyright [2019] - [2020], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceSDK contributors
 SPDX-License-Identifier: Apache-2.0
*/

package mutate

import (
	"github.com/asaskevich/govalidator"
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	xprtErrors "github.com/persistenceOne/persistenceSDK/constants/errors"
	"github.com/persistenceOne/persistenceSDK/modules/assets/internal/common"
	"github.com/persistenceOne/persistenceSDK/modules/assets/internal/module"
	"github.com/persistenceOne/persistenceSDK/schema/helpers"
	"github.com/persistenceOne/persistenceSDK/schema/types"
	codecUtilities "github.com/persistenceOne/persistenceSDK/utilities/codec"
)

type message struct {
	From                  sdkTypes.AccAddress  `json:"from" valid:"required~required field from missing"`
	FromID                types.ID             `json:"fromID" valid:"required~required field fromID missing"`
	AssetID               types.ID             `json:"assetID" valid:"required~required field assetID missing"`
	MutableMetaProperties types.MetaProperties `json:"mutableMetaProperties" valid:"required~required field mutableMetaProperties missing"`
	MutableProperties     types.Properties     `json:"mutableProperties" valid:"required~required field mutableProperties missing"`
}

var _ helpers.Message = message{}

func (message message) Route() string { return module.Name }
func (message message) Type() string  { return Transaction.GetName() }
func (message message) ValidateBasic() error {
	var _, Error = govalidator.ValidateStruct(message)
	if Error != nil {
		return errors.Wrap(xprtErrors.IncorrectMessage, Error.Error())
	}
	return nil
}
func (message message) GetSignBytes() []byte {
	return sdkTypes.MustSortJSON(common.Codec.MustMarshalJSON(message))
}
func (message message) GetSigners() []sdkTypes.AccAddress {
	return []sdkTypes.AccAddress{message.From}
}
func (message) RegisterCodec(codec *codec.Codec) {
	codecUtilities.RegisterXPRTConcrete(codec, module.Name, message{})
}
func messageFromInterface(msg sdkTypes.Msg) message {
	switch value := msg.(type) {
	case message:
		return value
	default:
		return message{}
	}
}
func messagePrototype() helpers.Message {
	return message{}
}
func newMessage(from sdkTypes.AccAddress, fromID types.ID, assetID types.ID, mutableMetaProperties types.MetaProperties, mutableProperties types.Properties) sdkTypes.Msg {
	return message{
		From:                  from,
		FromID:                fromID,
		AssetID:               assetID,
		MutableMetaProperties: mutableMetaProperties,
		MutableProperties:     mutableProperties,
	}
}