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
	results := make(chan TransferResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := s.TransferTx(context.Background(), TransferResponseParams{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	prevFromAccountBalance := fromAccount.Balance
	prevToAccountBalance := toAccount.Balance

	//check results
	for i := 0; i < n; i++ {
		err := <-errs
		result := <-results

		require.NoError(t, err)
		require.NotEmpty(t, result)

		require.Equal(t, result.FromAccount.ID, fromAccount.ID)
		require.Equal(t, result.FromAccount.Balance, prevFromAccountBalance-amount)

		require.Equal(t, result.ToAccount.ID, toAccount.ID)
		require.Equal(t, result.ToAccount.Balance, prevToAccountBalance+amount)

		require.Equal(t, result.FromEntry.AccountID, fromAccount.ID)
		require.Equal(t, result.FromEntry.Amount, -amount)

		require.Equal(t, result.ToEntry.AccountID, toAccount.ID)
		require.Equal(t, result.ToEntry.Amount, amount)

		require.Equal(t, result.Transfer.FromAccountID, fromAccount.ID)
		require.Equal(t, result.Transfer.ToAccountID, toAccount.ID)
		require.Equal(t, result.Transfer.Amount, amount)

		prevFromAccountBalance = result.FromAccount.Balance
		prevToAccountBalance = result.ToAccount.Balance
	}
}

func TestCreateTransferWithTxAndDeadlock(t *testing.T) {
	s := NewStore(testDB)

	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	amount := int64(10)

	n := 10
	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := fromAccount.ID
		toAccountID := toAccount.ID

		if i%2 == 1 {
			fromAccountID = toAccount.ID
			toAccountID = fromAccount.ID
		}

		go func() {
			_, err := s.TransferTx(context.Background(), TransferResponseParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	updatedFromAccount, _ := s.GetAccount(context.Background(), fromAccount.ID)
	updatedToAccount, _ := s.GetAccount(context.Background(), toAccount.ID)
	require.Equal(t, fromAccount.Balance, updatedFromAccount.Balance)
	require.Equal(t, toAccount.Balance, updatedToAccount.Balance)
}
