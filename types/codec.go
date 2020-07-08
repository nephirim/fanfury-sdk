package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

func RegisterCodec(codec *codec.Codec) {
	codec.RegisterInterface((*Chain)(nil), nil)
	codec.RegisterInterface((*Chains)(nil), nil)
	codec.RegisterInterface((*Classification)(nil), nil)
	codec.RegisterInterface((*Fact)(nil), nil)
	codec.RegisterInterface((*GenesisState)(nil), nil)
	codec.RegisterInterface((*Height)(nil), nil)
	codec.RegisterInterface((*ID)(nil), nil)
	codec.RegisterInterface((*Immutables)(nil), nil)
	codec.RegisterInterface((*InterNFT)(nil), nil)
	codec.RegisterInterface((*InterNFTs)(nil), nil)
	codec.RegisterInterface((*InterNFTWallet)(nil), nil)
	codec.RegisterInterface((*Maintainer)(nil), nil)
	codec.RegisterInterface((*Maintainers)(nil), nil)
	codec.RegisterInterface((*Mutables)(nil), nil)
	codec.RegisterInterface((*NFT)(nil), nil)
	codec.RegisterInterface((*NFTWallet)(nil), nil)
	codec.RegisterInterface((*Properties)(nil), nil)
	codec.RegisterInterface((*Property)(nil), nil)
	codec.RegisterInterface((*QueryRequest)(nil), nil)
	codec.RegisterInterface((*QueryResponse)(nil), nil)
	codec.RegisterInterface((*Request)(nil), nil)
	codec.RegisterInterface((*Share)(nil), nil)
	codec.RegisterInterface((*Signature)(nil), nil)
	codec.RegisterInterface((*Signatures)(nil), nil)
	codec.RegisterInterface((*Trait)(nil), nil)
	codec.RegisterInterface((*Traits)(nil), nil)
	codec.RegisterInterface((*TransactionRequest)(nil), nil)

	codec.RegisterConcrete(fact{}, "xprt/fact", nil)
	codec.RegisterConcrete(height{}, "xprt/height", nil)
	codec.RegisterConcrete(id{}, "xprt/id", nil)
	codec.RegisterConcrete(immutables{}, "xprt/immutables", nil)
	codec.RegisterConcrete(mutables{}, "xprt/mutables", nil)
	codec.RegisterConcrete(properties{}, "xprt/properties", nil)
	codec.RegisterConcrete(property{}, "xprt/property", nil)
	codec.RegisterConcrete(signature{}, "xprt/signature", nil)
	codec.RegisterConcrete(signatures{}, "xprt/signatures", nil)
}