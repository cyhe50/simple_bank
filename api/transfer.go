package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/cyhe50/simple_bank/db/sqlc"
	"github.com/cyhe50/simple_bank/token"
	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (s *Server) createTransfer(ctx *gin.Context) {
	var req *transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// client send unexpected request
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fromAccount, isValid := s.validAccount(ctx, req.FromAccountID, req.Currency)
	if !isValid {
		return
	}

	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if fromAccount.Owner != payload.Username {
		err := errors.New("fromAccount doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	}

	if _, isValid := s.validAccount(ctx, req.ToAccountID, req.Currency); !isValid {
		return
	}

	arg := db.TransferResponseParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	transfer, err := s.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transfer)
}

func (s *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := s.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency %s mismatch: expect to transfer %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	return account, true
}
