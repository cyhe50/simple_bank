package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateTransferWithTx(t *testing.T) {
	s := NewStore(testDB)

	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	amount := int64(10)

	// run n concurrent transfer transactions
	n := 5
	errs := make(chan error)
	results := make(chan transferResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := s.TransferTx(context.Background(), transferResponseParams{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	//check results
	for i := 0; i < n; i++ {
		err := <-errs
		result := <-results

		require.NoError(t, err)
		require.NotEmpty(t, result)

		// TODO: update account later
		//require.Equal(t, result.FromAccount, fromAccount)
		//require.Equal(t, result.ToAccount, toAccount)

		require.Equal(t, result.FromEntry.AccountID, fromAccount.ID)
		require.Equal(t, result.FromEntry.Amount, -amount)

		require.Equal(t, result.ToEntry.AccountID, toAccount.ID)
		require.Equal(t, result.ToEntry.Amount, amount)

		require.Equal(t, result.Transfer.FromAccountID, fromAccount.ID)
		require.Equal(t, result.Transfer.ToAccountID, toAccount.ID)
		require.Equal(t, result.Transfer.Amount, amount)
	}
}
