package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/khafizullokh02/simplebank/db/sqlc"
)

type CreateAccountRequest struct {
	Owner    string `json:"owner" biding:"required"`
	Currency string `json:"currency" biding:"required, oneof:USD EUR"`
}

func (server *Server) CreateAccount(ctx *gin.Context) {
	var req CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}
