package db

import (
	"testing"
	"github.com/stretchr/testify/require"
	"context"
	// "fmt"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	// fmt.Printf("account1 ID: %v\n", account1)
	// fmt.Printf("account2 ID: %v\n", account2)
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID: account2.ID,
				Amount: amount,
			})

			errs <- err
			results <- result
			// fmt.Println("----------------------------------------------------")
			// fmt.Printf("Transfer: %v\n", result.Transfer) 
			// fmt.Printf("FromAccount: %v\n", result.FromAccount) 
			// fmt.Printf("ToAccount: %v\n", result.ToAccount) 
			// fmt.Printf("FromEntry: %v\n", result.FromEntry) 
			// fmt.Printf("ToEntry: %v\n", result.ToEntry) 
		}()
	}
	for i := 0; i < n; i++ {
		err := <- errs
		require.NoError(t, err)

		result := <- results
		require.NotEmpty(t, result)

		// check transfers
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)


		// check entries out
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)


		// check entries in
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
		require.Positive(t, toEntry.Amount)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)


		// check from accounts  
		fromAccount := result.FromAccount

		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		// check into accounts 
		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// check balances
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1/amount)
		require.True(t, k >= 1)
	}
}