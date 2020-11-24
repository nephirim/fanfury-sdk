package provision

import (
	"github.com/persistenceOne/persistenceSDK/constants/errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Provision_Response(t *testing.T) {

	testTransactionResponse := newTransactionResponse(errors.IncorrectFormat)
	testTransactionResponse2 := newTransactionResponse(nil)

	require.Equal(t, transactionResponse{Success: false, Error: errors.IncorrectFormat}, testTransactionResponse)
	require.Equal(t, false, testTransactionResponse.IsSuccessful())
	require.Equal(t, true, testTransactionResponse2.IsSuccessful())

	require.Equal(t, errors.IncorrectFormat, testTransactionResponse.GetError())
	require.Equal(t, nil, testTransactionResponse2.GetError())
}
